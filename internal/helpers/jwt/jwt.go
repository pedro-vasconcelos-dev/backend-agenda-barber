package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/models"
)

type AccessClaims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	Plan      string `json:"plan"` // free | premium
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func secret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		// Melhor falhar logo do que rodar com secret vazio
		panic("JWT_SECRET not set")
	}
	return []byte(s)
}

func GenerateAccessToken(user models.User, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = time.Hour // padrão
	}

	claims := AccessClaims{
		UserID:    user.ID.String(),
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret())
}

func ValidateAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		// garante HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	// valida UUID
	if _, err := uuid.Parse(claims.UserID); err != nil {
		return nil, errors.New("invalid user_id")
	}

	return claims, nil
}
