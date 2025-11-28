package main

import (
	"fitness-market/internal/database"
	"fitness-market/internal/handlers"
	"fitness-market/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	database.Init()

	// Run migrations
	database.RunMigrations()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware for frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", handlers.Logout)
		auth.POST("/reset-password", handlers.RequestPasswordReset)
		auth.POST("/reset-password/confirm", handlers.ResetPassword)
	}

	// Protected routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", handlers.GetUserProfile)

		// Bodyweight routes
		api.POST("/profile/bodyweight", handlers.AddBodyweight)
		api.GET("/profile/bodyweight", handlers.GetBodyweightHistory)

		// Exercise PR routes
		api.POST("/profile/exercise-prs", handlers.AddExercisePR)
		api.GET("/profile/exercise-prs", handlers.GetExercisePRs)
		api.PUT("/profile/exercise-prs/:id", handlers.UpdateExercisePR)
		api.DELETE("/profile/exercise-prs/:id", handlers.DeleteExercisePR)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
