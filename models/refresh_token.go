package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken representa um token de atualização
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
