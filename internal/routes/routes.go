package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/auth_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/barbershop_handler"
	billing_handler "github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/billing_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/professional_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/stripe_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/handlers/user_handler"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/middlewares"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// Health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	r.POST("/stripe/webhook", func(ctx *gin.Context) { stripe_handler.HandleStripeWebhook(ctx, db) })

	// Auth
	auth := r.Group("/")
	{
		auth.POST("/auth/register", func(ctx *gin.Context) { auth_handler.Register(ctx, db) })
		auth.POST("/auth/login", func(ctx *gin.Context) { auth_handler.Login(ctx, db) })
		auth.POST("/auth/logout", middlewares.JWTAuthMiddleware(), func(ctx *gin.Context) { auth_handler.Logout(ctx) })
	}

	private := r.Group("/")
	{
		private.Use(middlewares.JWTAuthMiddleware())

		// Barbershops
		private.POST("/barbershops", func(ctx *gin.Context) { barbershop_handler.CreateBarbershop(ctx, db) })

		// Users
		private.GET("/users/me", func(ctx *gin.Context) { user_handler.GetMe(ctx, db) })

		// Professionals
		private.POST("/professionals", func(ctx *gin.Context) { professional_handler.CreateProfessional(ctx, db) })

		// Stripe/Billing
		private.POST("/billing/checkout-session", func(ctx *gin.Context) { billing_handler.CheckoutSession(ctx, db) })
		private.GET("/billing/status", func(ctx *gin.Context) { billing_handler.GetBillingStatus(ctx, db) })
		private.POST("/billing/cancel", func(ctx *gin.Context) { billing_handler.CancelSubscription(ctx, db) })
	}
}
