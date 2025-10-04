package main

import (
	"fitness-market/internal/auth"
	"fitness-market/internal/database"
	"fitness-market/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	database.InitDatabase()

	// Initialize Supabase
	auth.InitSupabase()

	// Setup router
	r := gin.Default()

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, gin.H{"user": user})
		})
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}