package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePremium() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		planVal, ok := ctx.Get("plan")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		plan, _ := planVal.(string)
		if plan != "premium" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "premium_required",
				"message": "this feature requires premium",
			})
			return
		}
		ctx.Next()
	}
}
