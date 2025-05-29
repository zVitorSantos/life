package models

import (
	"time"

	"gorm.io/gorm"
)

// GameProfile representa o perfil de jogo de um usuário
type GameProfile struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null;uniqueIndex"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	// Dados básicos do jogo
	Level int   `json:"level" gorm:"default:1"`
	XP    int64 `json:"xp" gorm:"default:0"`

	// Status do jogador
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	LastLogin *time.Time `json:"last_login"`

	// Estatísticas genéricas (flexível para qualquer tipo de jogo)
	Stats map[string]interface{} `json:"stats" gorm:"type:jsonb"`

	// Configurações do jogador
	Settings map[string]interface{} `json:"settings" gorm:"type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName especifica o nome da tabela
func (GameProfile) TableName() string {
	return "game_profiles"
}

// GetXPForNextLevel calcula XP necessário para o próximo level
func (gp *GameProfile) GetXPForNextLevel() int64 {
	// Fórmula genérica: level * 1000 (pode ser customizada depois)
	return int64(gp.Level * 1000)
}

// GetXPProgress retorna o progresso atual para o próximo level (0-100)
func (gp *GameProfile) GetXPProgress() float64 {
	if gp.Level == 1 && gp.XP == 0 {
		return 0
	}

	currentLevelXP := int64((gp.Level - 1) * 1000)
	nextLevelXP := gp.GetXPForNextLevel()
	progressXP := gp.XP - currentLevelXP

	if progressXP <= 0 {
		return 0
	}

	return float64(progressXP) / float64(nextLevelXP-currentLevelXP) * 100
}

// AddXP adiciona XP e verifica se subiu de level
func (gp *GameProfile) AddXP(amount int64) bool {
	gp.XP += amount

	// Verifica se subiu de level
	for gp.XP >= gp.GetXPForNextLevel() {
		gp.Level++
	}

	return true
}

// GetStat retorna uma estatística específica
func (gp *GameProfile) GetStat(key string) interface{} {
	if gp.Stats == nil {
		return nil
	}
	return gp.Stats[key]
}

// SetStat define uma estatística
func (gp *GameProfile) SetStat(key string, value interface{}) {
	if gp.Stats == nil {
		gp.Stats = make(map[string]interface{})
	}
	gp.Stats[key] = value
}

// GetSetting retorna uma configuração específica
func (gp *GameProfile) GetSetting(key string) interface{} {
	if gp.Settings == nil {
		return nil
	}
	return gp.Settings[key]
}

// SetSetting define uma configuração
func (gp *GameProfile) SetSetting(key string, value interface{}) {
	if gp.Settings == nil {
		gp.Settings = make(map[string]interface{})
	}
	gp.Settings[key] = value
}
