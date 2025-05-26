package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken representa um token de atualização
// @Description Informações do token de atualização
type RefreshToken struct {
	// ID único do token
	ID uint `json:"id" gorm:"primaryKey" example:"1"`

	// Token de atualização
	Token string `json:"token" gorm:"uniqueIndex;not null" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// ID do usuário dono do token
	UserID uint `json:"user_id" gorm:"not null" example:"1"`

	// Data de expiração do token
	ExpiresAt time.Time `json:"expires_at" gorm:"not null" example:"2024-12-31T23:59:59Z"`

	// Indica se o token foi revogado
	IsRevoked bool `json:"is_revoked" gorm:"default:false" example:"false"`

	// Data de criação
	CreatedAt time.Time `json:"created_at" example:"2024-05-25T20:00:00Z"`

	// Data da última atualização
	UpdatedAt time.Time `json:"updated_at" example:"2024-05-25T20:00:00Z"`

	// Data de exclusão (soft delete)
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
