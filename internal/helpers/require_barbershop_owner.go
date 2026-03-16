package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RequireBarbershopOwner valida que o usuário autenticado é owner da barbearia.
// Recebe o barbershop_id do body JSON (já extraído pelo caller).
// Retorna o barbershopID parseado, ou escreve o erro no ctx e retorna false.
func RequireBarbershopOwner(ctx *gin.Context, db *gorm.DB, barbershopIDStr string) (uuid.UUID, bool) {
	userID := GetAuthenticatedUUID(ctx)
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	barbershopID, err := uuid.Parse(barbershopIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid barbershop_id"})
		return uuid.Nil, false
	}

	ok, err := IsBarbershopOwner(db, userID, barbershopID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate barbershop access"})
		return uuid.Nil, false
	}
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only the barbershop owner can manage billing"})
		return uuid.Nil, false
	}

	return barbershopID, true
}
