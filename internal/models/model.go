package models

import (
	"time"

	"gorm.io/gorm"
)

// Model base model
type Model struct {
	ID        int32          `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
