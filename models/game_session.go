package models

import (
	"time"

	"gorm.io/gorm"
)

// SessionStatus define o status da sessão
type SessionStatus string

const (
	SessionActive     SessionStatus = "active"     // Sessão ativa
	SessionInactive   SessionStatus = "inactive"   // Sessão inativa
	SessionExpired    SessionStatus = "expired"    // Sessão expirada
	SessionTerminated SessionStatus = "terminated" // Sessão terminada forçadamente
)

// GameSession representa uma sessão de jogo ativa
type GameSession struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	GameProfileID uint        `json:"game_profile_id" gorm:"not null"`
	GameProfile   GameProfile `json:"game_profile" gorm:"foreignKey:GameProfileID"`

	// Dados da sessão
	Status       SessionStatus `json:"status" gorm:"default:active"`
	StartedAt    time.Time     `json:"started_at"`
	LastActivity time.Time     `json:"last_activity"`
	EndedAt      *time.Time    `json:"ended_at"`

	// Informações técnicas
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Platform  string `json:"platform"`

	// Estatísticas da sessão
	Duration     int64 `json:"duration"` // em segundos
	ActionsCount int   `json:"actions_count"`

	// Dados da sessão (flexível)
	SessionData map[string]interface{} `json:"session_data" gorm:"type:jsonb"`

	// Controle de segurança
	IsValid       bool   `json:"is_valid" gorm:"default:true"`
	InvalidReason string `json:"invalid_reason"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName especifica o nome da tabela
func (GameSession) TableName() string {
	return "game_sessions"
}

// IsActive verifica se a sessão está ativa
func (gs *GameSession) IsActive() bool {
	return gs.Status == SessionActive && gs.IsValid
}

// IsExpired verifica se a sessão expirou (mais de 30 minutos sem atividade)
func (gs *GameSession) IsExpired() bool {
	if gs.Status != SessionActive {
		return true
	}

	// Considera expirada se não há atividade há mais de 30 minutos
	return time.Since(gs.LastActivity) > 30*time.Minute
}

// UpdateActivity atualiza o timestamp da última atividade
func (gs *GameSession) UpdateActivity() {
	gs.LastActivity = time.Now()
	gs.ActionsCount++
}

// End termina a sessão
func (gs *GameSession) End() {
	now := time.Now()
	gs.Status = SessionInactive
	gs.EndedAt = &now
	gs.Duration = int64(now.Sub(gs.StartedAt).Seconds())
}

// Expire marca a sessão como expirada
func (gs *GameSession) Expire() {
	now := time.Now()
	gs.Status = SessionExpired
	gs.EndedAt = &now
	gs.Duration = int64(now.Sub(gs.StartedAt).Seconds())
}

// Terminate termina a sessão forçadamente
func (gs *GameSession) Terminate(reason string) {
	now := time.Now()
	gs.Status = SessionTerminated
	gs.EndedAt = &now
	gs.Duration = int64(now.Sub(gs.StartedAt).Seconds())
	gs.IsValid = false
	gs.InvalidReason = reason
}

// GetDurationMinutes retorna a duração em minutos
func (gs *GameSession) GetDurationMinutes() int64 {
	if gs.EndedAt != nil {
		return gs.Duration / 60
	}
	// Se ainda está ativa, calcula duração atual
	return int64(time.Since(gs.StartedAt).Minutes())
}

// GetSessionData retorna um dado específico da sessão
func (gs *GameSession) GetSessionData(key string) interface{} {
	if gs.SessionData == nil {
		return nil
	}
	return gs.SessionData[key]
}

// SetSessionData define um dado da sessão
func (gs *GameSession) SetSessionData(key string, value interface{}) {
	if gs.SessionData == nil {
		gs.SessionData = make(map[string]interface{})
	}
	gs.SessionData[key] = value
}

// CanPerformAction verifica se pode realizar ações (sessão ativa e válida)
func (gs *GameSession) CanPerformAction() bool {
	return gs.IsActive() && !gs.IsExpired()
}

// GetActivityStatus retorna o status de atividade da sessão
func (gs *GameSession) GetActivityStatus() string {
	if !gs.IsActive() {
		return "offline"
	}

	timeSinceActivity := time.Since(gs.LastActivity)

	if timeSinceActivity < 5*time.Minute {
		return "online"
	} else if timeSinceActivity < 15*time.Minute {
		return "away"
	} else {
		return "idle"
	}
}
