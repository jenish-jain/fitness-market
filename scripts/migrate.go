package main

import (
	"fitness-market/internal/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	log.Println("Running database migrations...")

	// Initialize database connection and run migrations
	database.Init()

	log.Println("Migration completed successfully!")

	// Close database connection
	database.Close()
}