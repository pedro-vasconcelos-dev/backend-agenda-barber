package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/configurations"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/middlewares"
	"github.com/legitimatech-rpa/backend-agenda-barber/internal/routes"
)

func main() {
	_ = godotenv.Load()

	db, err := configurations.NewGORMPostgresConnection()
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	gin.SetMode(gin.ReleaseMode)

	server := gin.Default()

	server.Use(middlewares.CORSMiddleware())

	routes.RegisterRoutes(server, db)

	server.Run("localhost:0235")
}
