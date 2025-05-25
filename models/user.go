package models

import (
	"time"

	"gorm.io/gorm"
)

// User representa um usuário no sistema
// @Description Informações do usuário
type User struct {
	// ID único do usuário
	ID uint `json:"id" gorm:"primaryKey" example:"1"`

	// Nome de usuário único
	Username string `json:"username" gorm:"unique;not null" example:"johndoe"`

	// Nome de exibição
	DisplayName string `json:"display_name" gorm:"not null" example:"John Doe"`

	// Email do usuário
	Email string `json:"email" gorm:"unique;not null" example:"john@example.com"`

	// Senha do usuário (não serializada)
	Password string `json:"-" gorm:"not null"`

	// Data de criação
	CreatedAt time.Time `json:"created_at" example:"2024-05-25T20:00:00Z"`

	// Data da última atualização
	UpdatedAt time.Time `json:"updated_at" example:"2024-05-25T20:00:00Z"`

	// Data de exclusão (soft delete)
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
