package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt_helper "github.com/legitimatech-rpa/backend-agenda-barber/internal/helpers/jwt"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := jwt_helper.ValidateAccessToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			return
		}

		// coloca no context para handlers usarem
		ctx.Set("user_id", claims.UserID) // string UUID
		ctx.Set("role", claims.Role)      // tutor | establishment
		ctx.Set("plan", claims.Plan)      // free | premium

		ctx.Next()
	}
}
