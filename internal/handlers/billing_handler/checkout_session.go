package billing_handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"

	"github.com/stripe/stripe-go/v78"
	checkoutsession "github.com/stripe/stripe-go/v78/checkout/session"
	customerapi "github.com/stripe/stripe-go/v78/customer"
	priceapi "github.com/stripe/stripe-go/v78/price"
)

type CheckoutSessionRequest struct {
	BarbershopID string `json:"barbershop_id"`
}

type CheckoutSessionResponse struct {
	CheckoutURL string `json:"checkout_url"`
}

type BillingStatusResponse struct {
	BarbershopID      string     `json:"barbershop_id"`
	PlanLookupKey     string     `json:"plan_lookup_key"`
	Status            string     `json:"status"`
	CurrentPeriodEnd  *time.Time `json:"current_period_end"`
	CancelAtPeriodEnd bool       `json:"cancel_at_period_end"`
}

type CancelSubscriptionRequest struct {
	BarbershopID string `json:"barbershop_id"`
}

type CancelSubscriptionResponse struct {
	OK                   bool       `json:"ok"`
	StripeSubscriptionID string     `json:"stripe_subscription_id,omitempty"`
	CancelAtPeriodEnd    bool       `json:"cancel_at_period_end"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end,omitempty"`
	Status               string     `json:"status,omitempty"`
}

// --- Stripe helpers ---

func getPriceIDByLookupKey(lookupKey string) (string, error) {
	params := &stripe.PriceListParams{
		LookupKeys: []*string{stripe.String(lookupKey)},
		Active:     stripe.Bool(true),
	}

	it := priceapi.List(params)
	for it.Next() {
		return it.Price().ID, nil
	}

	if err := it.Err(); err != nil {
		return "", err
	}

	return "", &stripe.Error{
		Code: stripe.ErrorCodeResourceMissing,
		Msg:  "price not found for lookup_key: " + lookupKey,
	}
}

// requireBarbershopOwner valida que o usuário autenticado é owner da barbearia.
// Recebe o barbershop_id do body JSON (já extraído pelo caller).
// Retorna o barbershopID parseado, ou escreve o erro no ctx e retorna false.
func requireBarbershopOwner(ctx *gin.Context, db *gorm.DB, barbershopIDStr string) (uuid.UUID, bool) {
	userID := helpers.GetAuthenticatedUUID(ctx)
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	barbershopID, err := uuid.Parse(barbershopIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid barbershop_id"})
		return uuid.Nil, false
	}

	ok, err := helpers.IsBarbershopOwner(db, userID, barbershopID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate barbershop access"})
		return uuid.Nil, false
	}
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only the barbershop owner can manage billing"})
		return uuid.Nil, false
	}

	return barbershopID, true
}

func CheckoutSession(ctx *gin.Context, db *gorm.DB) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "stripe not configured"})
		return
	}

	var req CheckoutSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.BarbershopID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id is required"})
		return
	}

	barbershopID, ok := requireBarbershopOwner(ctx, db, req.BarbershopID)
	if !ok {
		return
	}

	const lookupKey = "assinatura_mensal_padrao"

	// Busca ou cria BillingCustomer para a barbearia
	var bc models.BillingCustomer
	if err := db.Where("barbershop_id = ?", barbershopID).First(&bc).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Busca e-mail do usuário autenticado para registrar no Stripe
		userID := helpers.GetAuthenticatedUUID(ctx)
		var user models.User
		if dbErr := db.Where("id = ?", userID).First(&user).Error; dbErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
			return
		}

		cusParams := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Metadata: map[string]string{
				"barbershop_id": barbershopID.String(),
			},
		}

		cus, stripeErr := customerapi.New(cusParams)
		if stripeErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create stripe customer"})
			return
		}

		bc = models.BillingCustomer{
			BarbershopID:     barbershopID,
			StripeCustomerID: cus.ID,
		}

		if err := db.Create(&bc).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save billing customer"})
			return
		}
	}

	successURL := os.Getenv("STRIPE_SUCCESS_URL")
	cancelURL := os.Getenv("STRIPE_CANCEL_URL")
	if successURL == "" || cancelURL == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "stripe redirect URLs not configured"})
		return
	}

	priceID, err := getPriceIDByLookupKey(lookupKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve price: " + err.Error()})
		return
	}

	params := &stripe.CheckoutSessionParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer:          stripe.String(bc.StripeCustomerID),
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
		ClientReferenceID: stripe.String(barbershopID.String()),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"barbershop_id":           barbershopID.String(),
			"stripe_price_lookup_key": lookupKey,
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"barbershop_id":           barbershopID.String(),
				"stripe_price_lookup_key": lookupKey,
			},
		},
	}

	s, err := checkoutsession.New(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create checkout session: " + err.Error()})
		return
	}

	if s.URL == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid checkout session - no URL generated: " + s.ID})
		return
	}

	ctx.JSON(http.StatusOK, CheckoutSessionResponse{CheckoutURL: s.URL})
}
