package billing_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

func GetBillingStatus(ctx *gin.Context, db *gorm.DB) {
	barbershopIDStr := ctx.Query("barbershop_id")
	if barbershopIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id is required"})
		return
	}

	barbershopID, ok := requireBarbershopOwner(ctx, db, barbershopIDStr)
	if !ok {
		return
	}

	var subscription models.BillingSubscription
	err := db.Where("barbershop_id = ? AND status IN ?", barbershopID, []string{"active", "trialing"}).
		Order("created_at DESC").
		First(&subscription).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, BillingStatusResponse{
				BarbershopID:      barbershopID.String(),
				PlanLookupKey:     "",
				Status:            "inactive",
				CurrentPeriodEnd:  nil,
				CancelAtPeriodEnd: false,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query billing status"})
		return
	}

	ctx.JSON(http.StatusOK, BillingStatusResponse{
		BarbershopID:      subscription.BarbershopID.String(),
		PlanLookupKey:     subscription.StripePriceLookupKey,
		Status:            subscription.Status,
		CurrentPeriodEnd:  subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd: subscription.CancelAtPeriodEnd,
	})
}
