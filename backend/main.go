package main

import (
	"log"
	"net/http"

	"launchpad-backend/config"
	"launchpad-backend/handlers"
	"launchpad-backend/middleware"
	"launchpad-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := services.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db.DB, cfg.JWTSecret)
	tokenHandler := handlers.NewTokenHandler(db.DB)
	presaleHandler := handlers.NewPresaleHandler(db.DB)

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "launchpad-backend",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		token := api.Group("/token")
		{
			token.POST("/create", middleware.AuthMiddleware(cfg.JWTSecret), tokenHandler.CreateToken)
		}

		presale := api.Group("/presale")
		{
			presale.POST("/create", middleware.AuthMiddleware(cfg.JWTSecret), presaleHandler.CreatePresale)
			presale.GET("/:id", presaleHandler.GetPresale)
			presale.POST("/:id/participate", presaleHandler.Participate)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}