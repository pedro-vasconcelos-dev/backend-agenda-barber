package professional_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers/validation"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

type CreateProfessionalRequest struct {
	BarbershopID string `json:"barbershop_id" binding:"required,uuid4"`
	Name         string `json:"name" binding:"required,min=2,max=120"`
	Phone        string `json:"phone" binding:"required,max=20"`
	Email        string `json:"email" binding:"required,email"`
	Role         string `json:"role" binding:"required,oneof=professional reception"`
}

type CreateProfessionalResponse struct {
	ID           uuid.UUID `json:"id"`
	BarbershopID uuid.UUID `json:"barbershop_id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Password     string    `json:"password"`
}

func CreateProfessional(ctx *gin.Context, db *gorm.DB) {
	var req CreateProfessionalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	barbershopID, ok := helpers.RequireBarbershopOwner(ctx, db, req.BarbershopID)
	if !ok {
		return
	}

	password, err := validation.GenRandomPassword(12)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate password"})
		return
	}

	var resp CreateProfessionalResponse
	txErr := db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		user.Name = req.Name
		user.Email = validation.NormalizeEmail(req.Email)
		user.Phone = req.Phone
		user.PasswordHash, err = validation.HashPassword(password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return err
		}

		if err := tx.Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return err
		}

		professional := models.Professional{
			BarbershopID: barbershopID,
			UserID:       user.ID,
			Name:         req.Name,
			Phone:        req.Phone,
		}

		if err := tx.Create(&professional).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create professional"})
			return err
		}

		var barbershopUser models.BarbershopUser
		barbershopUser.BarbershopID = barbershopID
		barbershopUser.UserID = user.ID
		barbershopUser.Role = req.Role

		if err := tx.Create(&barbershopUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create barbershop user"})
			return err
		}

		resp = CreateProfessionalResponse{
			ID:           professional.ID,
			BarbershopID: barbershopID,
			Name:         professional.Name,
			Phone:        professional.Phone,
			Email:        user.Email,
			Role:         barbershopUser.Role,
			Password:     password,
		}
		return nil
	})

	if txErr != nil {
		// O erro já foi retornado e resposta já enviada
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}
