package professional_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"gorm.io/gorm"
)

func ListProfessionals(ctx *gin.Context, db *gorm.DB) {
	barbershopIDStr := ctx.Query("barbershop_id")
	if barbershopIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id is required"})
		return
	}

	if _, ok := helpers.RequireBarbershopOwner(ctx, db, barbershopIDStr); !ok {
		return
	}

	var rows []professionalScan
	if err := professionalQuery(db).
		Where("p.barbershop_id = ? AND p.is_active = true", barbershopIDStr).
		Scan(&rows).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query professionals"})
		return
	}

	resp := make([]ProfessionalResponse, 0, len(rows))
	for _, r := range rows {
		resp = append(resp, r.toResponse())
	}

	ctx.JSON(http.StatusOK, resp)
}
