package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Exercise struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Ticker      string         `json:"ticker" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Category    string         `json:"category" gorm:"not null"`
	StockPrice  float64        `json:"stock_price" gorm:"not null;default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User           User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	WorkoutEntries []WorkoutEntry `json:"workout_entries,omitempty" gorm:"foreignKey:ExerciseID"`
}

// TableName specifies the table name for Exercise
func (Exercise) TableName() string {
	return "exercises"
}

// BeforeCreate hook to ensure unique ticker per user
func (e *Exercise) BeforeCreate(tx *gorm.DB) error {
	var count int64
	tx.Model(&Exercise{}).Where("user_id = ? AND ticker = ?", e.UserID, e.Ticker).Count(&count)
	if count > 0 {
		return errors.New("ticker symbol already exists for this user")
	}
	return nil
}

// BeforeUpdate hook to ensure unique ticker per user on update
func (e *Exercise) BeforeUpdate(tx *gorm.DB) error {
	var count int64
	tx.Model(&Exercise{}).Where("user_id = ? AND ticker = ? AND id != ?", e.UserID, e.Ticker, e.ID).Count(&count)
	if count > 0 {
		return errors.New("ticker symbol already exists for this user")
	}
	return nil
}
