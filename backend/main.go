package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

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
			auth.POST("/login", handleLogin)
		}

		token := api.Group("/token")
		{
			token.POST("/create", authMiddleware(), handleTokenCreate)
		}

		presale := api.Group("/presale")
		{
			presale.POST("/create", authMiddleware(), handlePresaleCreate)
			presale.GET("/:id", handlePresaleGet)
			presale.POST("/:id/participate", handlePresaleParticipate)
		}
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Placeholder handlers - will be implemented in separate files
func handleLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func handleTokenCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func handlePresaleCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func handlePresaleGet(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func handlePresaleParticipate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT verification
		c.Next()
	}
}