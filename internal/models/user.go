package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	Email             string         `json:"email" gorm:"uniqueIndex;not null"`
	Password          string         `json:"-" gorm:"not null"`
	Name              string         `json:"name"`
	ResetToken        string         `json:"-" gorm:"index"`
	ResetTokenExpiry  *time.Time     `json:"-"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}
