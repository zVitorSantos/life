package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"unique;not null"`
	DisplayName string         `json:"display_name" gorm:"not null"`
	Email       string         `json:"email" gorm:"unique;not null"`
	Password    string         `json:"-" gorm:"not null"` // "-" para não serializar a senha
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
