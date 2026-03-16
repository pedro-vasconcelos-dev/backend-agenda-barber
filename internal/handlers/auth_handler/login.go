package auth_handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers/jwt"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers/validation"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthUserResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	IsAdmin bool      `json:"is_admin"`
}

type AuthResponse struct {
	AccessToken string           `json:"access_token"`
	User        AuthUserResponse `json:"user"`
}

func toUserResponse(u models.User) AuthUserResponse {
	return AuthUserResponse{
		ID:      u.ID,
		Name:    u.Name,
		Email:   u.Email,
		Phone:   u.Phone,
		IsAdmin: u.IsAdmin,
	}
}

func Login(ctx *gin.Context, db *gorm.DB) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	req.Email = validation.NormalizeEmail(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !validation.CheckPassword(user.PasswordHash, req.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := jwt.GenerateAccessToken(user, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, AuthResponse{
		AccessToken: token,
		User:        toUserResponse(user),
	})
}
