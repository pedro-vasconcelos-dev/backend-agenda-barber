package barbershop_handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

type CreateBarbershopRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=120"`
	Slug     string `json:"slug" binding:"required,min=3,max=80"`
	Timezone string `json:"timezone" binding:"omitempty,max=100"`
	Address  string `json:"address" binding:"max=255"`
	Phone    string `json:"phone" binding:"max=20"`
}

type CreateBarbershopResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Timezone string    `json:"timezone"`
	Address  string    `json:"address"`
	Phone    string    `json:"phone"`
	Role     string    `json:"role"`
}

func CreateBarbershop(ctx *gin.Context, db *gorm.DB) {
	var req CreateBarbershopRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID := helpers.GetAuthenticatedUUID(ctx)
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Slug = normalizeSlug(req.Slug)
	req.Timezone = strings.TrimSpace(req.Timezone)
	req.Address = strings.TrimSpace(req.Address)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.Timezone == "" {
		req.Timezone = "America/Sao_Paulo"
	}

	if req.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if req.Slug == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	if req.Phone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone is required"})
		return
	}

	var createdBarbershop models.Barbershop

	if err := db.Transaction(func(tx *gorm.DB) error {
		barbershop := models.Barbershop{
			Name:     req.Name,
			Slug:     req.Slug,
			Timezone: req.Timezone,
			Address:  req.Address,
			Phone:    req.Phone,
			IsActive: true,
		}

		if err := tx.Create(&barbershop).Error; err != nil {
			return err
		}

		barbershopUser := models.BarbershopUser{
			BarbershopID: barbershop.ID,
			UserID:       uuid.MustParse(userID),
			Role:         "owner",
			IsActive:     true,
		}

		if err := tx.Create(&barbershopUser).Error; err != nil {
			return err
		}

		createdBarbershop = barbershop
		return nil
	}); err != nil {
		if isUniqueViolation(err) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Slug already in use"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create barbershop"})
		return
	}

	ctx.JSON(http.StatusCreated, CreateBarbershopResponse{
		ID:       createdBarbershop.ID,
		Name:     createdBarbershop.Name,
		Slug:     createdBarbershop.Slug,
		Timezone: createdBarbershop.Timezone,
		Address:  createdBarbershop.Address,
		Phone:    createdBarbershop.Phone,
		Role:     "owner",
	})
}

func normalizeSlug(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	return value
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
