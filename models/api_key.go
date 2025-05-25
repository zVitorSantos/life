package models

import (
	"time"

	"gorm.io/gorm"
)

// APIKey representa uma chave de API no sistema
// @Description Informações da chave de API
type APIKey struct {
	// ID único da chave
	ID uint `json:"id" gorm:"primaryKey" example:"1"`

	// Nome da chave para identificação
	Name string `json:"name" gorm:"not null" example:"Frontend App"`

	// Chave de API (hash)
	Key string `json:"key" gorm:"unique;not null"`

	// ID do usuário dono da chave
	UserID uint `json:"user_id" gorm:"not null"`

	// Data de expiração
	ExpiresAt time.Time `json:"expires_at" example:"2024-12-31T23:59:59Z"`

	// Último uso da chave
	LastUsedAt *time.Time `json:"last_used_at" example:"2024-05-25T20:00:00Z"`

	// Limite de requisições por minuto
	RateLimit int `json:"rate_limit" gorm:"default:60" example:"60"`

	// Status da chave (ativo/inativo)
	IsActive bool `json:"is_active" gorm:"default:true"`

	// Data de criação
	CreatedAt time.Time `json:"created_at" example:"2024-05-25T20:00:00Z"`

	// Data da última atualização
	UpdatedAt time.Time `json:"updated_at" example:"2024-05-25T20:00:00Z"`

	// Data de exclusão (soft delete)
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
