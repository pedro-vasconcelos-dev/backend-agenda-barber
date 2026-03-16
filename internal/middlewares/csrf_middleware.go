package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CSRFOpts struct {
	CookieName string // "csrf_token"
	HeaderName string // "X-CSRF-Token"
}

func CSRFMiddleware(opts CSRFOpts) gin.HandlerFunc {
	// origem de dev permitida por fallback (pode vir de env)
	devOrigin := "http://localhost:5173"

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			csrfCookie, _ := c.Cookie(opts.CookieName)
			csrfHeader := c.GetHeader(opts.HeaderName)
			origin := c.GetHeader("Origin")

			// --- Caminho normal (prod): cookie e header DEVEM bater ---
			if csrfCookie != "" && csrfHeader != "" && csrfCookie == csrfHeader {
				c.Next()
				return
			}

			// --- Fallback DEV: se a origem for o localhost de dev e houver header, aceite ---
			// (ÚNICO alvo: facilitar desenvolvimento quando o front não consegue ler o cookie de outro domínio)
			if csrfCookie == "" && csrfHeader != "" && strings.EqualFold(origin, devOrigin) {
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "invalid csrf token",
			})
			return
		}
		c.Next()
	}
}
