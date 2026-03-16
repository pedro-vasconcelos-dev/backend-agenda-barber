package professional_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers/validation"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

type UpdateProfessionalRequest struct {
	Name  *string `json:"name"  binding:"omitempty,min=2,max=120"`
	Phone *string `json:"phone" binding:"omitempty,max=20"`
	Email *string `json:"email" binding:"omitempty,email"`
	Role  *string `json:"role"  binding:"omitempty,oneof=professional reception"`
}

func UpdateProfessional(ctx *gin.Context, db *gorm.DB) {
	professionalIDStr := ctx.Param("id")
	if professionalIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "professional_id is required"})
		return
	}

	var req UpdateProfessionalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Name == nil && req.Phone == nil && req.Email == nil && req.Role == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	var professional models.Professional
	if err := db.First(&professional, "id = ? AND is_active = true", professionalIDStr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "professional not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query professional"})
		return
	}

	if _, ok := helpers.RequireBarbershopOwner(ctx, db, professional.BarbershopID.String()); !ok {
		return
	}

	txErr := db.Transaction(func(tx *gorm.DB) error {
		// Atualiza users (name, phone, email são compartilhados)
		userUpdates := map[string]interface{}{}
		if req.Name != nil {
			userUpdates["name"] = *req.Name
		}
		if req.Phone != nil {
			userUpdates["phone"] = *req.Phone
		}
		if req.Email != nil {
			userUpdates["email"] = validation.NormalizeEmail(*req.Email)
		}
		if len(userUpdates) > 0 {
			if err := tx.Model(&models.User{}).
				Where("id = ?", professional.UserID).
				Updates(userUpdates).Error; err != nil {
				return err
			}
		}

		// Atualiza professionals (name e phone espelham users)
		profUpdates := map[string]interface{}{}
		if req.Name != nil {
			profUpdates["name"] = *req.Name
		}
		if req.Phone != nil {
			profUpdates["phone"] = *req.Phone
		}
		if len(profUpdates) > 0 {
			if err := tx.Model(&models.Professional{}).
				Where("id = ?", professionalIDStr).
				Updates(profUpdates).Error; err != nil {
				return err
			}
		}

		// Atualiza role em barbershop_users
		if req.Role != nil {
			if err := tx.Model(&models.BarbershopUser{}).
				Where("user_id = ? AND barbershop_id = ?", professional.UserID, professional.BarbershopID).
				Update("role", *req.Role).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update professional"})
		return
	}

	var row professionalScan
	if err := professionalQuery(db).
		Where("p.id = ?", professionalIDStr).
		Scan(&row).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load updated professional"})
		return
	}

	ctx.JSON(http.StatusOK, row.toResponse())
}
