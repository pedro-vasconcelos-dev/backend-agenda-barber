package professional_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"gorm.io/gorm"
)

func GetProfessionalByID(ctx *gin.Context, db *gorm.DB) {
	professionalIDStr := ctx.Param("id")
	if professionalIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "professional_id is required"})
		return
	}

	// Busca profissional com join para obter email e role antes de checar autorização
	var row professionalScan
	if err := professionalQuery(db).
		Where("p.id = ? AND p.is_active = true", professionalIDStr).
		Scan(&row).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query professional"})
		return
	}
	if row.ID.String() == "00000000-0000-0000-0000-000000000000" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "professional not found"})
		return
	}

	if _, ok := helpers.RequireBarbershopOwner(ctx, db, row.BarbershopID.String()); !ok {
		return
	}

	ctx.JSON(http.StatusOK, row.toResponse())
}
