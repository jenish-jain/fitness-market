package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutEntry struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	ExerciseID uint           `json:"exercise_id" gorm:"index;not null"`
	Weight     float64        `json:"weight" gorm:"not null"`
	Reps       int            `json:"reps" gorm:"not null"`
	Sets       int            `json:"sets" gorm:"not null"`
	Notes      string         `json:"notes"`
	Date       time.Time      `json:"date" gorm:"not null"`
	Score      float64        `json:"score" gorm:"not null;default:0"`
	IsPR       bool           `json:"is_pr" gorm:"not null;default:false"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User     `json:"-" gorm:"foreignKey:UserID"`
	Exercise Exercise `json:"exercise,omitempty" gorm:"foreignKey:ExerciseID"`
}

func (WorkoutEntry) TableName() string {
	return "workout_entries"
}
