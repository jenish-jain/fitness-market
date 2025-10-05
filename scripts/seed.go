package main

import (
	"fitness-market/internal/database"
	"fitness-market/internal/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	database.Init()
	db := database.GetDB()

	// Seed test users
	testUsers := []models.User{
		{
			Email:      "test@example.com",
			Name:       "Test User",
			SupabaseID: "test-supabase-id-1",
			Role:       "user",
			Active:     true,
		},
		{
			Email:      "admin@example.com",
			Name:       "Admin User",
			SupabaseID: "test-supabase-id-2",
			Role:       "admin",
			Active:     true,
		},
	}

	for _, user := range testUsers {
		var existingUser models.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error != nil {
			// User doesn't exist, create it
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", user.Email, err)
			} else {
				log.Printf("Created user: %s", user.Email)
			}
		} else {
			log.Printf("User %s already exists, skipping", user.Email)
		}
	}

	log.Println("Seeding completed")
}