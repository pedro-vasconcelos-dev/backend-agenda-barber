package middlewares

import (
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func isAllowedOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := u.Hostname() // só o host, sem esquema/porta

	// localhost dev
	if (u.Scheme == "http" || u.Scheme == "https") &&
		(host == "localhost" && (u.Port() == "8080" || u.Port() == "")) {
		// aceita http://localhost:8080/ e https://localhost:8080/
		return true
	}

	// qualquer subdomínio do seu domínio + root
	if host == "legitimatech.com.br" || strings.HasSuffix(host, ".legitimatech.com.br") {
		return true
	}
	return false
}

func CORSMiddleware() gin.HandlerFunc {
	cfg := cors.Config{
		AllowOriginFunc: isAllowedOrigin, // nada de ""

		AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:          []string{"Set-Cookie"},
		AllowCredentials:       true, // necessário para cookie/withCredentials
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          false,
		AllowBrowserExtensions: true, // opcional
	}

	return cors.New(cfg)
}
