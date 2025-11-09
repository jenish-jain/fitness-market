package models

import (
	"time"

	"gorm.io/gorm"
)

type Exercise struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Category    string    `json:"category" gorm:"not null"`
	StockPrice  float64   `json:"stock_price" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User           User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	WorkoutEntries []WorkoutEntry `json:"workout_entries,omitempty" gorm:"foreignKey:ExerciseID"`
}