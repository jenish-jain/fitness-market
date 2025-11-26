package models

import (
	"time"

	"gorm.io/gorm"
)

type PortfolioSnapshot struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	Date          time.Time `json:"date" gorm:"not null;index"`
	TotalValue    float64   `json:"total_value" gorm:"not null;default:0"`
	WorkoutCount  int       `json:"workout_count" gorm:"not null;default:0"`
	ActiveStreaks int       `json:"active_streaks" gorm:"not null;default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}