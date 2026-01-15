package database

import (
	"fitness-market/internal/models"
	"log"
)

func RunMigrations() {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.BodyweightEntry{},
		&models.ExercisePR{},
		&models.Exercise{},
		&models.WorkoutEntry{},
		&models.PRHistory{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
