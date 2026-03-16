package billing_handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"github.com/stripe/stripe-go/v78"
	subapi "github.com/stripe/stripe-go/v78/subscription"
	"gorm.io/gorm"
)

func CancelSubscription(ctx *gin.Context, db *gorm.DB) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "stripe not configured"})
		return
	}

	var req CancelSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.BarbershopID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id is required"})
		return
	}

	barbershopID, ok := requireBarbershopOwner(ctx, db, req.BarbershopID)
	if !ok {
		return
	}

	var subscription models.BillingSubscription
	if err := db.Where("barbershop_id = ? AND status IN ?", barbershopID, []string{"active", "trialing", "past_due"}).
		Order("created_at DESC").
		First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no active subscription"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load subscription"})
		return
	}

	if subscription.StripeSubscriptionID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "missing stripe_subscription_id"})
		return
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	updated, err := subapi.Update(subscription.StripeSubscriptionID, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to cancel subscription"})
		return
	}

	tm := time.Unix(updated.CurrentPeriodEnd, 0)

	if err := db.Model(&models.BillingSubscription{}).
		Where("id = ?", subscription.ID).
		Updates(map[string]interface{}{
			"cancel_at_period_end": true,
			"current_period_end":   tm,
			"status":               string(updated.Status),
			"updated_at":           time.Now(),
		}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist cancellation"})
		return
	}

	cpe := &tm
	ctx.JSON(http.StatusOK, CancelSubscriptionResponse{
		OK:                   true,
		StripeSubscriptionID: subscription.StripeSubscriptionID,
		CancelAtPeriodEnd:    true,
		CurrentPeriodEnd:     cpe,
		Status:               string(updated.Status),
	})
}
