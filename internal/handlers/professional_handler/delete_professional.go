package professional_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

func DeleteProfessional(ctx *gin.Context, db *gorm.DB) {
	professionalIDStr := ctx.Param("id")
	if professionalIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "professional_id is required"})
		return
	}

	professionalID, err := uuid.Parse(professionalIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid professional_id"})
		return
	}

	var professional models.Professional
	if err := db.First(&professional, "id = ?", professionalID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "professional not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query professional"})
		return
	}

	barbershopID := professional.BarbershopID
	if _, ok := helpers.RequireBarbershopOwner(ctx, db, barbershopID.String()); !ok {
		return
	}

	var user models.User
	if err := db.Delete(&user, "id = ?", professional.UserID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete associated user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "professional deleted successfully"})
}
