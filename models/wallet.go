package models

import (
	"time"

	"gorm.io/gorm"
)

// CurrencyType define os tipos de moeda disponíveis
type CurrencyType string

const (
	CurrencyCoins  CurrencyType = "coins"  // Moeda principal do jogo
	CurrencyGems   CurrencyType = "gems"   // Moeda premium
	CurrencyTokens CurrencyType = "tokens" // Moeda especial/evento
)

// Wallet representa a carteira de um usuário
type Wallet struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	GameProfileID uint        `json:"game_profile_id" gorm:"not null"`
	GameProfile   GameProfile `json:"game_profile" gorm:"foreignKey:GameProfileID"`

	// Saldos por tipo de moeda
	CoinsBalance  int64 `json:"coins_balance" gorm:"default:0"`
	GemsBalance   int64 `json:"gems_balance" gorm:"default:0"`
	TokensBalance int64 `json:"tokens_balance" gorm:"default:0"`

	// Controle de segurança
	IsLocked   bool   `json:"is_locked" gorm:"default:false"`
	LockReason string `json:"lock_reason"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName especifica o nome da tabela
func (Wallet) TableName() string {
	return "wallets"
}

// GetBalance retorna o saldo de um tipo específico de moeda
func (w *Wallet) GetBalance(currency CurrencyType) int64 {
	switch currency {
	case CurrencyCoins:
		return w.CoinsBalance
	case CurrencyGems:
		return w.GemsBalance
	case CurrencyTokens:
		return w.TokensBalance
	default:
		return 0
	}
}

// SetBalance define o saldo de um tipo específico de moeda
func (w *Wallet) SetBalance(currency CurrencyType, amount int64) {
	switch currency {
	case CurrencyCoins:
		w.CoinsBalance = amount
	case CurrencyGems:
		w.GemsBalance = amount
	case CurrencyTokens:
		w.TokensBalance = amount
	}
}

// AddBalance adiciona valor ao saldo (pode ser negativo para subtrair)
func (w *Wallet) AddBalance(currency CurrencyType, amount int64) int64 {
	currentBalance := w.GetBalance(currency)
	newBalance := currentBalance + amount

	// Não permite saldo negativo
	if newBalance < 0 {
		newBalance = 0
	}

	w.SetBalance(currency, newBalance)
	return newBalance
}

// HasSufficientBalance verifica se tem saldo suficiente
func (w *Wallet) HasSufficientBalance(currency CurrencyType, amount int64) bool {
	return w.GetBalance(currency) >= amount
}

// CanSpend verifica se pode gastar (não está bloqueada e tem saldo)
func (w *Wallet) CanSpend(currency CurrencyType, amount int64) bool {
	return !w.IsLocked && w.HasSufficientBalance(currency, amount)
}

// Lock bloqueia a carteira
func (w *Wallet) Lock(reason string) {
	w.IsLocked = true
	w.LockReason = reason
}

// Unlock desbloqueia a carteira
func (w *Wallet) Unlock() {
	w.IsLocked = false
	w.LockReason = ""
}

// GetTotalValue retorna o valor total da carteira (para ranking/estatísticas)
// Usando uma conversão simples: 1 gem = 100 coins, 1 token = 10 coins
func (w *Wallet) GetTotalValue() int64 {
	return w.CoinsBalance + (w.GemsBalance * 100) + (w.TokensBalance * 10)
}
