package auth_handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logout(ctx *gin.Context) {
	// Define headers para limpar cache do browser
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	// Limpa qualquer cookie que possa existir no domínio
	// (caso o frontend esteja usando cookies para armazenar tokens)
	clearCookies := []string{
		"access_token", "refresh_token", "session_id",
		"csrf_token", "auth_token", "token",
	}

	for _, cookieName := range clearCookies {
		ctx.SetCookie(
			cookieName, // name
			"",         // value (vazio)
			-1,         // maxAge (negativo para deletar)
			"/",        // path
			"",         // domain (vazio para usar o domínio atual)
			false,      // secure (ajuste conforme seu ambiente)
			true,       // httpOnly
		)
	}

	// Resposta indicando sucesso do logout
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "logout successful",
		"timestamp":    time.Now().Unix(),
		"instructions": "Please clear the access token from your client storage",
	})
}
