package auth_handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers/jwt"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers/validation"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
}

func Register(ctx *gin.Context, db *gorm.DB) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if email and phone are unique
	var existingUser models.User
	if req.Email != "" {
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
			return
		}
	}

	if req.Phone != "" {
		if err := db.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already in use"})
			return
		}
	}

	hashPassword, err := validation.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashPassword,
	}

	if err := db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := jwt.GenerateAccessToken(user, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{UserID: user.ID, Token: token})
}
