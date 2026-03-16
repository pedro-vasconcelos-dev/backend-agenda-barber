package helpers

import "github.com/gin-gonic/gin"

func GetAuthenticatedUUID(ctx *gin.Context) string {

	if val, ok := ctx.Get("user_id"); ok {
		if s, ok := val.(string); ok {
			return s
		}
	}

	return ""
}
