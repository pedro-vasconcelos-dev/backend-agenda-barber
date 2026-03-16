package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/auth_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/barbershop_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/user_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/middlewares"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// Health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// r.POST("/stripe/webhook", stripeWebhook.Handle)

	// Auth
	auth := r.Group("/")
	{
		auth.POST("/auth/register", func(ctx *gin.Context) { auth_handler.Register(ctx, db) })
		auth.POST("/auth/login", func(ctx *gin.Context) { auth_handler.Login(ctx, db) })
		auth.POST("/auth/logout", middlewares.JWTAuthMiddleware(), func(ctx *gin.Context) { auth_handler.Logout(ctx) })
	}

	// Barbershop
	barbershops := r.Group("/barbershops")
	{
		barbershops.Use(middlewares.JWTAuthMiddleware())
		// Exemplo de uso do middleware de acesso por barbearia:
		// barbershops.Use(middlewares.BarbershopAccessMiddleware(db))
		barbershops.POST("/", func(ctx *gin.Context) { barbershop_handler.CreateBarbershop(ctx, db) })
	}

	// User
	user := r.Group("/users")
	{
		user.Use(middlewares.JWTAuthMiddleware())
		user.GET("/me", func(ctx *gin.Context) { user_handler.GetMe(ctx, db) })
	}
}
