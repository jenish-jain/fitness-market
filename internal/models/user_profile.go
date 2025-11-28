package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type BodyweightEntry struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Weight    float64        `json:"weight" gorm:"not null"`
	Unit      string         `json:"unit" gorm:"default:'kg'"`
	RecordedAt time.Time    `json:"recorded_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ExercisePR struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	UserID       uint           `json:"user_id" gorm:"index;not null"`
	User         User           `json:"-" gorm:"foreignKey:UserID"`
	ExerciseName string         `json:"exercise_name" gorm:"not null"`
	Weight       float64        `json:"weight" gorm:"not null"`
	Unit         string         `json:"unit" gorm:"default:'kg'"`
	Reps         int            `json:"reps" gorm:"default:1"`
	RecordedAt   time.Time      `json:"recorded_at" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
