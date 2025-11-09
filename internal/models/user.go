package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Exercises      []Exercise      `json:"exercises,omitempty" gorm:"foreignKey:UserID"`
	WorkoutEntries []WorkoutEntry  `json:"workout_entries,omitempty" gorm:"foreignKey:UserID"`
	PortfolioSnapshots []PortfolioSnapshot `json:"portfolio_snapshots,omitempty" gorm:"foreignKey:UserID"`
}