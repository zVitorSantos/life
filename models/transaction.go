package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType define os tipos de transação
type TransactionType string

const (
	TransactionEarn     TransactionType = "earn"     // Ganhou dinheiro
	TransactionSpend    TransactionType = "spend"    // Gastou dinheiro
	TransactionTransfer TransactionType = "transfer" // Transferência entre usuários
	TransactionReward   TransactionType = "reward"   // Recompensa do sistema
	TransactionPenalty  TransactionType = "penalty"  // Penalidade/multa
	TransactionRefund   TransactionType = "refund"   // Reembolso
)

// TransactionStatus define o status da transação
type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "pending"   // Pendente
	TransactionCompleted TransactionStatus = "completed" // Concluída
	TransactionFailed    TransactionStatus = "failed"    // Falhou
	TransactionCancelled TransactionStatus = "cancelled" // Cancelada
	TransactionReversed  TransactionStatus = "reversed"  // Revertida
)

// Transaction representa uma movimentação financeira
type Transaction struct {
	ID uint `json:"id" gorm:"primaryKey"`

	// Relacionamentos
	WalletID uint   `json:"wallet_id" gorm:"not null"`
	Wallet   Wallet `json:"wallet" gorm:"foreignKey:WalletID"`

	// Dados da transação
	Type     TransactionType   `json:"type" gorm:"not null"`
	Status   TransactionStatus `json:"status" gorm:"default:pending"`
	Currency CurrencyType      `json:"currency" gorm:"not null"`
	Amount   int64             `json:"amount" gorm:"not null"`

	// Saldos (para auditoria)
	BalanceBefore int64 `json:"balance_before"`
	BalanceAfter  int64 `json:"balance_after"`

	// Metadados
	Description string `json:"description"`
	Reference   string `json:"reference"` // ID de referência (compra, quest, etc.)
	Category    string `json:"category"`  // Categoria da transação

	// Dados adicionais (flexível)
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`

	// Transferência (se aplicável)
	ToWalletID *uint   `json:"to_wallet_id"`
	ToWallet   *Wallet `json:"to_wallet" gorm:"foreignKey:ToWalletID"`

	// Controle de reversão
	ReversedByID *uint        `json:"reversed_by_id"`
	ReversedBy   *Transaction `json:"reversed_by" gorm:"foreignKey:ReversedByID"`
	ReversesID   *uint        `json:"reverses_id"`
	Reverses     *Transaction `json:"reverses" gorm:"foreignKey:ReversesID"`

	// Timestamps
	ProcessedAt *time.Time     `json:"processed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName especifica o nome da tabela
func (Transaction) TableName() string {
	return "transactions"
}

// IsCompleted verifica se a transação foi concluída
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionCompleted
}

// IsPending verifica se a transação está pendente
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionPending
}

// CanBeReversed verifica se a transação pode ser revertida
func (t *Transaction) CanBeReversed() bool {
	return t.IsCompleted() && t.ReversedByID == nil
}

// Complete marca a transação como concluída
func (t *Transaction) Complete() {
	t.Status = TransactionCompleted
	now := time.Now()
	t.ProcessedAt = &now
}

// Fail marca a transação como falhada
func (t *Transaction) Fail() {
	t.Status = TransactionFailed
	now := time.Now()
	t.ProcessedAt = &now
}

// Cancel marca a transação como cancelada
func (t *Transaction) Cancel() {
	t.Status = TransactionCancelled
	now := time.Now()
	t.ProcessedAt = &now
}

// Reverse marca a transação como revertida
func (t *Transaction) Reverse(reversalTransaction *Transaction) {
	t.Status = TransactionReversed
	t.ReversedByID = &reversalTransaction.ID
}

// GetMetadata retorna um metadado específico
func (t *Transaction) GetMetadata(key string) interface{} {
	if t.Metadata == nil {
		return nil
	}
	return t.Metadata[key]
}

// SetMetadata define um metadado
func (t *Transaction) SetMetadata(key string, value interface{}) {
	if t.Metadata == nil {
		t.Metadata = make(map[string]interface{})
	}
	t.Metadata[key] = value
}

// IsTransfer verifica se é uma transferência entre carteiras
func (t *Transaction) IsTransfer() bool {
	return t.Type == TransactionTransfer && t.ToWalletID != nil
}

// GetAbsoluteAmount retorna o valor absoluto da transação
func (t *Transaction) GetAbsoluteAmount() int64 {
	if t.Amount < 0 {
		return -t.Amount
	}
	return t.Amount
}
