package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	IsAdmin bool      `json:"is_admin"`
}

type BarbershopResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Address  string    `json:"address"`
	IsActive bool      `json:"is_active"`
	Role     string    `json:"role"`
}

type GetMeResponse struct {
	User        UserResponse         `json:"user"`
	Barbershops []BarbershopResponse `json:"barbershops"`
}

func GetMe(ctx *gin.Context, db *gorm.DB) {
	userID := helpers.GetAuthenticatedUUID(ctx)
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	// Buscar barbearias do usuário via join em barbershop_users e barbershops
	type BarbershopWithRole struct {
		BarbershopID uuid.UUID
		Name         string
		Slug         string
		Address      string
		IsActive     bool
		Role         string
	}
	var barbershopsWithRole []BarbershopWithRole
	if err := db.Table("barbershop_users bu").
		Select("bu.barbershop_id, b.name, b.slug, b.address, b.is_active, bu.role").
		Joins("JOIN barbershops b ON b.id = bu.barbershop_id").
		Where("bu.user_id = ?", userID).
		Scan(&barbershopsWithRole).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve barbershops"})
		return
	}

	response := GetMeResponse{
		User: UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Phone:   user.Phone,
			IsAdmin: user.IsAdmin,
		},
	}

	for _, b := range barbershopsWithRole {
		response.Barbershops = append(response.Barbershops, BarbershopResponse{
			ID:       b.BarbershopID,
			Name:     b.Name,
			Slug:     b.Slug,
			Address:  b.Address,
			IsActive: b.IsActive,
			Role:     b.Role,
		})
	}

	ctx.JSON(http.StatusOK, response)
}
