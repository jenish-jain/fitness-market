package scripts

import (
	"fitness-market/internal/database"
	"fitness-market/internal/models"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func Seed() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	log.Println("Seeding database...")

	// Initialize database
	database.Init()
	db := database.GetDB()

	// Create seed users
	users := []models.User{
		{
			Email: "john.doe@example.com",
			Name:  "John Doe",
		},
		{
			Email: "jane.smith@example.com",
			Name:  "Jane Smith",
		},
		{
			Email: "mike.wilson@example.com",
			Name:  "Mike Wilson",
		},
	}

	for _, user := range users {
		var existingUser models.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error != nil {
			db.Create(&user)
			log.Printf("Created user: %s", user.Email)
		} else {
			log.Printf("User already exists: %s", user.Email)
		}
	}

	// Create seed exercises
	exercises := []models.Exercise{
		{
			UserID:      1,
			Name:        "Push-ups",
			Description: "Classic bodyweight exercise for chest and triceps",
			Category:    "Strength",
			StockPrice:  10.50,
		},
		{
			UserID:      1,
			Name:        "Squats",
			Description: "Lower body strength exercise",
			Category:    "Strength",
			StockPrice:  15.75,
		},
		{
			UserID:      2,
			Name:        "Running",
			Description: "Cardiovascular endurance exercise",
			Category:    "Cardio",
			StockPrice:  8.25,
		},
		{
			UserID:      2,
			Name:        "Plank",
			Description: "Core strengthening exercise",
			Category:    "Core",
			StockPrice:  12.00,
		},
		{
			UserID:      3,
			Name:        "Deadlifts",
			Description: "Full body compound movement",
			Category:    "Strength",
			StockPrice:  25.50,
		},
	}

	for _, exercise := range exercises {
		var existingExercise models.Exercise
		result := db.Where("user_id = ? AND name = ?", exercise.UserID, exercise.Name).First(&existingExercise)
		if result.Error != nil {
			db.Create(&exercise)
			log.Printf("Created exercise: %s", exercise.Name)
		} else {
			log.Printf("Exercise already exists: %s", exercise.Name)
		}
	}

	// Create seed workout entries
	now := time.Now()
	workoutEntries := []models.WorkoutEntry{
		{
			UserID:     1,
			ExerciseID: 1,
			Date:       now.AddDate(0, 0, -2),
			Sets:       3,
			Reps:       15,
			Weight:     0,
			Notes:      "Good form, felt strong",
		},
		{
			UserID:     1,
			ExerciseID: 2,
			Date:       now.AddDate(0, 0, -2),
			Sets:       3,
			Reps:       20,
			Weight:     0,
			Notes:      "Legs feeling good",
		},
		{
			UserID:     2,
			ExerciseID: 3,
			Date:       now.AddDate(0, 0, -1),
			Sets:       1,
			Reps:       1,
			Duration:   1800, // 30 minutes
			Notes:      "5K run in the park",
		},
		{
			UserID:     3,
			ExerciseID: 5,
			Date:       now.AddDate(0, 0, -1),
			Sets:       4,
			Reps:       8,
			Weight:     135.0,
			Notes:      "PR! Felt great",
		},
	}

	for _, entry := range workoutEntries {
		var existingEntry models.WorkoutEntry
		result := db.Where("user_id = ? AND exercise_id = ? AND date = ?", entry.UserID, entry.ExerciseID, entry.Date).First(&existingEntry)
		if result.Error != nil {
			db.Create(&entry)
			log.Printf("Created workout entry for user %d, exercise %d", entry.UserID, entry.ExerciseID)
		} else {
			log.Printf("Workout entry already exists for user %d, exercise %d", entry.UserID, entry.ExerciseID)
		}
	}

	// Create seed portfolio snapshots
	portfolioSnapshots := []models.PortfolioSnapshot{
		{
			UserID:        1,
			Date:          now.AddDate(0, 0, -7),
			TotalValue:    850.75,
			WorkoutCount:  12,
			ActiveStreaks: 3,
		},
		{
			UserID:        2,
			Date:          now.AddDate(0, 0, -7),
			TotalValue:    675.50,
			WorkoutCount:  8,
			ActiveStreaks: 2,
		},
		{
			UserID:        3,
			Date:          now.AddDate(0, 0, -7),
			TotalValue:    1250.25,
			WorkoutCount:  15,
			ActiveStreaks: 5,
		},
	}

	for _, snapshot := range portfolioSnapshots {
		var existingSnapshot models.PortfolioSnapshot
		result := db.Where("user_id = ? AND date = ?", snapshot.UserID, snapshot.Date).First(&existingSnapshot)
		if result.Error != nil {
			db.Create(&snapshot)
			log.Printf("Created portfolio snapshot for user %d", snapshot.UserID)
		} else {
			log.Printf("Portfolio snapshot already exists for user %d", snapshot.UserID)
		}
	}

	log.Println("Database seeding completed successfully!")

	// Close database connection
	database.Close()
}