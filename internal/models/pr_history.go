package models

import (
	"time"

	"gorm.io/gorm"
)

// PRHistory maintains the history of personal records per exercise
type PRHistory struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	UserID         uint           `json:"user_id" gorm:"index;not null"`
	ExerciseID     uint           `json:"exercise_id" gorm:"index;not null"`
	WorkoutEntryID uint           `json:"workout_entry_id" gorm:"index;not null"`
	Score          float64        `json:"score" gorm:"not null"`
	Weight         float64        `json:"weight" gorm:"not null"`
	Reps           int            `json:"reps" gorm:"not null"`
	Sets           int            `json:"sets" gorm:"not null"`
	AchievedAt     time.Time      `json:"achieved_at" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User         User         `json:"-" gorm:"foreignKey:UserID"`
	Exercise     Exercise     `json:"exercise,omitempty" gorm:"foreignKey:ExerciseID"`
	WorkoutEntry WorkoutEntry `json:"workout_entry,omitempty" gorm:"foreignKey:WorkoutEntryID"`
}

func (PRHistory) TableName() string {
	return "pr_history"
}
