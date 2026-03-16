package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers"
	"gorm.io/gorm"
)

// BarbershopAccessMiddleware garante que o usuário autenticado pertence à barbearia informada
// (qualquer role ativo). Lê barbershop_id de query, form ou path param.
func BarbershopAccessMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := helpers.GetAuthenticatedUUID(ctx)
		if userID == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		barbershopIDStr := ctx.Query("barbershop_id")
		if barbershopIDStr == "" {
			barbershopIDStr = ctx.PostForm("barbershop_id")
		}
		if barbershopIDStr == "" {
			barbershopIDStr = ctx.Param("barbershop_id")
		}
		if barbershopIDStr == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "barbershop_id is required"})
			return
		}

		barbershopID, err := uuid.Parse(barbershopIDStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid barbershop_id"})
			return
		}

		ok, err := helpers.HasBarbershopAccess(db, userID, barbershopID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate barbershop access"})
			return
		}
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have access to this barbershop"})
			return
		}

		ctx.Next()
	}
}
