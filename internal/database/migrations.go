package database

import (
	"fitness-market/internal/models"
	"log"
)

func RunMigrations() {
	if DB == nil {
		log.Println("Database not initialized, skipping migrations")
		return
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.BodyweightEntry{},
		&models.ExercisePR{},
	)
	if err != nil {
		log.Printf("Migration error: %v", err)
	}
}
