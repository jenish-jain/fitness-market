package database

import (
	"fitness-market/internal/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./fitness_market.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")

	// Auto-migrate the schema
	AutoMigrate()
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Exercise{},
		&models.WorkoutEntry{},
		&models.PortfolioSnapshot{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Create indexes
	createIndexes()
}

func createIndexes() {
	// Create composite indexes for better query performance
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_workout_entries_user_date ON workout_entries(user_id, date)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_workout_entries_exercise_date ON workout_entries(exercise_id, date)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_portfolio_snapshots_user_date ON portfolio_snapshots(user_id, date)")

	log.Println("Database indexes created")
}

func GetDB() *gorm.DB {
	return DB
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}