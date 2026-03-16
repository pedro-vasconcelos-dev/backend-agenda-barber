package stripe_handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"github.com/stripe/stripe-go/v78"
	subapi "github.com/stripe/stripe-go/v78/subscription"
	"github.com/stripe/stripe-go/v78/webhook"
	"gorm.io/gorm"
)

func HandleStripeWebhook(ctx *gin.Context, db *gorm.DB) {
	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if secret == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "missing webhook secret"})
		return
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, ctx.Request.Body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	event, err := webhook.ConstructEventWithOptions(
		buf.Bytes(),
		ctx.GetHeader("Stripe-Signature"),
		secret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	switch event.Type {

	// ==============================
	// 1) Checkout finalizado → cria subscription (pending)
	// ==============================
	case "checkout.session.completed":
		var cs stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &cs); err != nil {
			break
		}

		if cs.Subscription == nil || cs.Customer == nil {
			break
		}

		lookupKey := cs.Metadata["stripe_price_lookup_key"]

		barbershopID, _ := uuid.Parse(cs.ClientReferenceID)

		sub := models.BillingSubscription{
			BarbershopID:         barbershopID,
			StripeCustomerID:     cs.Customer.ID,
			StripeSubscriptionID: cs.Subscription.ID,
			StripePriceLookupKey: lookupKey,
			Status:               "pending",
			CancelAtPeriodEnd:    false,
		}

		if err := upsertBillingSubscription(db, &sub); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

	// ==============================
	// 2) PAGAMENTO CONFIRMADO → ativa plano
	// ==============================
	case "invoice_payment.paid":
		var payload struct {
			Invoice struct {
				Subscription string `json:"subscription"`
				Customer     string `json:"customer"`
			} `json:"invoice"`
		}

		if err := json.Unmarshal(event.Data.Raw, &payload); err != nil {
			break
		}

		if payload.Invoice.Subscription == "" {
			break
		}

		if err := db.Model(&models.BillingSubscription{}).
			Where("stripe_subscription_id = ?", payload.Invoice.Subscription).
			Updates(map[string]interface{}{
				"status":             "active",
				"stripe_customer_id": payload.Invoice.Customer,
				"updated_at":         time.Now(),
			}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to activate"})
			return
		}

		// ==============================
		// 3) Atualização da subscription (cancelamento, downgrade, etc)
		// ==============================
	case "customer.subscription.updated", "customer.subscription.deleted":
		var s stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
			break
		}

		status := string(s.Status)
		if event.Type == "customer.subscription.deleted" {
			status = "canceled"
		}

		// ✅ Fonte de verdade: Stripe (garante current_period_end preenchido)
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		if stripe.Key == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "stripe not configured"})
			return
		}

		fresh, err := subapi.Get(s.ID, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subscription from stripe"})
			return
		}

		cpe := unixPtr(fresh.CurrentPeriodEnd)

		if err := db.Model(&models.BillingSubscription{}).
			Where("stripe_subscription_id = ?", s.ID).
			Updates(map[string]interface{}{
				"status":               status,
				"cancel_at_period_end": fresh.CancelAtPeriodEnd,
				"current_period_end":   cpe,
				"updated_at":           time.Now(),
			}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "sync failed"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// ==============================
// Helpers
// ==============================
func unixPtr(v int64) *time.Time {
	if v == 0 {
		return nil
	}
	t := time.Unix(v, 0)
	return &t
}

func upsertBillingSubscription(db *gorm.DB, sub *models.BillingSubscription) error {
	var existing models.BillingSubscription
	err := db.Where("stripe_subscription_id = ?", sub.StripeSubscriptionID).First(&existing).Error
	if err == nil {
		return db.Model(&existing).Updates(sub).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return db.Create(sub).Error
}
