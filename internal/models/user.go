package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Name      string         `gorm:"not null" json:"name"`
	SupabaseID string        `gorm:"uniqueIndex;not null" json:"supabase_id"`
	Role      string         `gorm:"default:user" json:"role"`
	Active    bool           `gorm:"default:true" json:"active"`
}