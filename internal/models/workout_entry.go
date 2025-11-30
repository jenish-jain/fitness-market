package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutEntry struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	ExerciseID uint           `json:"exercise_id" gorm:"not null;index"`
	Date       time.Time      `json:"date" gorm:"not null;index"`
	Sets       int            `json:"sets" gorm:"not null"`
	Reps       int            `json:"reps" gorm:"not null"`
	Weight     float64        `json:"weight"`
	Duration   int            `json:"duration"`
	Notes      string         `json:"notes"`
	Score      float64        `json:"score"`
	IsPR       bool           `json:"is_pr" gorm:"default:false"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Exercise Exercise `json:"exercise,omitempty" gorm:"foreignKey:ExerciseID"`
}
