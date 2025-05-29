package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"life/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TransactionHandler gerencia operações relacionadas às transações
type TransactionHandler struct {
	db *gorm.DB
}

// NewTransactionHandler cria uma nova instância do TransactionHandler
func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

// AddMoney adiciona dinheiro à carteira do usuário
func (h *TransactionHandler) AddMoney(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var request struct {
		Currency    string `json:"currency" binding:"required"`
		Amount      int64  `json:"amount" binding:"required,min=1"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Reference   string `json:"reference"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Valida o tipo de moeda
	currencyType, err := h.validateCurrency(request.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Busca a carteira
	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	// Verifica se a carteira não está bloqueada
	if wallet.IsLocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Carteira bloqueada: " + wallet.LockReason})
		return
	}

	// Executa a transação
	transaction, err := h.executeTransaction(wallet, models.TransactionEarn, currencyType, request.Amount, request.Description, request.Category, request.Reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar transação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Dinheiro adicionado com sucesso",
		"transaction": transaction,
		"new_balance": wallet.GetBalance(currencyType),
	})
}

// SpendMoney remove dinheiro da carteira do usuário
func (h *TransactionHandler) SpendMoney(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var request struct {
		Currency    string `json:"currency" binding:"required"`
		Amount      int64  `json:"amount" binding:"required,min=1"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Reference   string `json:"reference"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Valida o tipo de moeda
	currencyType, err := h.validateCurrency(request.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Busca a carteira
	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	// Verifica se pode gastar
	if !wallet.CanSpend(currencyType, request.Amount) {
		if wallet.IsLocked {
			c.JSON(http.StatusForbidden, gin.H{"error": "Carteira bloqueada: " + wallet.LockReason})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo insuficiente"})
		return
	}

	// Executa a transação (valor negativo para gasto)
	transaction, err := h.executeTransaction(wallet, models.TransactionSpend, currencyType, -request.Amount, request.Description, request.Category, request.Reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar transação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Dinheiro gasto com sucesso",
		"transaction": transaction,
		"new_balance": wallet.GetBalance(currencyType),
	})
}

// TransferMoney transfere dinheiro entre carteiras
func (h *TransactionHandler) TransferMoney(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var request struct {
		ToUserID    uint   `json:"to_user_id" binding:"required"`
		Currency    string `json:"currency" binding:"required"`
		Amount      int64  `json:"amount" binding:"required,min=1"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Não pode transferir para si mesmo
	if userID.(uint) == request.ToUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Não é possível transferir para si mesmo"})
		return
	}

	// Valida o tipo de moeda
	currencyType, err := h.validateCurrency(request.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Busca carteira de origem
	fromWallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carteira de origem não encontrada"})
		return
	}

	// Busca carteira de destino
	toWallet, err := h.getWalletByUserID(request.ToUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carteira de destino não encontrada"})
		return
	}

	// Verifica se pode transferir
	if !fromWallet.CanSpend(currencyType, request.Amount) {
		if fromWallet.IsLocked {
			c.JSON(http.StatusForbidden, gin.H{"error": "Sua carteira está bloqueada: " + fromWallet.LockReason})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo insuficiente"})
		return
	}

	if toWallet.IsLocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Carteira de destino está bloqueada"})
		return
	}

	// Executa a transferência em uma transação do banco
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Transação de saída
		fromTransaction := &models.Transaction{
			WalletID:      fromWallet.ID,
			Type:          models.TransactionTransfer,
			Status:        models.TransactionPending,
			Currency:      currencyType,
			Amount:        -request.Amount,
			BalanceBefore: fromWallet.GetBalance(currencyType),
			Description:   fmt.Sprintf("Transferência para usuário %d", request.ToUserID),
			Category:      "transfer_out",
			Reference:     fmt.Sprintf("transfer_%d_%d", userID, request.ToUserID),
			ToWalletID:    &toWallet.ID,
		}

		fromTransaction.BalanceAfter = fromTransaction.BalanceBefore + fromTransaction.Amount

		if err := tx.Create(fromTransaction).Error; err != nil {
			return err
		}

		// Atualiza saldo da carteira de origem
		fromWallet.AddBalance(currencyType, -request.Amount)
		if err := tx.Save(fromWallet).Error; err != nil {
			return err
		}

		// Transação de entrada
		toTransaction := &models.Transaction{
			WalletID:      toWallet.ID,
			Type:          models.TransactionTransfer,
			Status:        models.TransactionPending,
			Currency:      currencyType,
			Amount:        request.Amount,
			BalanceBefore: toWallet.GetBalance(currencyType),
			Description:   fmt.Sprintf("Transferência de usuário %d", userID),
			Category:      "transfer_in",
			Reference:     fmt.Sprintf("transfer_%d_%d", userID, request.ToUserID),
		}

		toTransaction.BalanceAfter = toTransaction.BalanceBefore + toTransaction.Amount

		if err := tx.Create(toTransaction).Error; err != nil {
			return err
		}

		// Atualiza saldo da carteira de destino
		toWallet.AddBalance(currencyType, request.Amount)
		if err := tx.Save(toWallet).Error; err != nil {
			return err
		}

		// Marca ambas as transações como concluídas
		fromTransaction.Complete()
		toTransaction.Complete()

		if err := tx.Save(fromTransaction).Error; err != nil {
			return err
		}

		if err := tx.Save(toTransaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar transferência"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transferência realizada com sucesso",
		"amount":      request.Amount,
		"currency":    request.Currency,
		"to_user_id":  request.ToUserID,
		"new_balance": fromWallet.GetBalance(currencyType),
	})
}

// GetTransactionHistory obtém o histórico de transações do usuário
func (h *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Parâmetros de filtro
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	transactionType := c.Query("type")
	currency := c.Query("currency")
	status := c.Query("status")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
		return
	}

	query := h.db.Where("wallet_id = ?", wallet.ID)

	// Aplica filtros
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}
	if currency != "" {
		query = query.Where("currency = ?", currency)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var transactions []models.Transaction
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar histórico"})
		return
	}

	// Conta total
	var total int64
	query.Model(&models.Transaction{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
		"filters": gin.H{
			"type":     transactionType,
			"currency": currency,
			"status":   status,
		},
	})
}

// GetTransaction obtém uma transação específica
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	transactionIDStr := c.Param("id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da transação inválido"})
		return
	}

	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
		return
	}

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND wallet_id = ?", transactionID, wallet.ID).
		Preload("ToWallet").
		Preload("ReversedBy").
		Preload("Reverses").
		First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transação não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar transação"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// executeTransaction executa uma transação de forma segura
func (h *TransactionHandler) executeTransaction(wallet *models.Wallet, transactionType models.TransactionType, currency models.CurrencyType, amount int64, description, category, reference string) (*models.Transaction, error) {
	var transaction *models.Transaction

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Cria a transação
		transaction = &models.Transaction{
			WalletID:      wallet.ID,
			Type:          transactionType,
			Status:        models.TransactionPending,
			Currency:      currency,
			Amount:        amount,
			BalanceBefore: wallet.GetBalance(currency),
			Description:   description,
			Category:      category,
			Reference:     reference,
		}

		transaction.BalanceAfter = transaction.BalanceBefore + amount

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// Atualiza o saldo da carteira
		wallet.AddBalance(currency, amount)
		if err := tx.Save(wallet).Error; err != nil {
			return err
		}

		// Marca a transação como concluída
		transaction.Complete()
		if err := tx.Save(transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// validateCurrency valida o tipo de moeda
func (h *TransactionHandler) validateCurrency(currency string) (models.CurrencyType, error) {
	switch currency {
	case "coins":
		return models.CurrencyCoins, nil
	case "gems":
		return models.CurrencyGems, nil
	case "tokens":
		return models.CurrencyTokens, nil
	default:
		return "", fmt.Errorf("tipo de moeda inválido: %s", currency)
	}
}

// getWalletByUserID busca carteira por user_id
func (h *TransactionHandler) getWalletByUserID(userID uint) (*models.Wallet, error) {
	var gameProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).First(&gameProfile).Error; err != nil {
		return nil, err
	}

	var wallet models.Wallet
	if err := h.db.Where("game_profile_id = ?", gameProfile.ID).First(&wallet).Error; err != nil {
		return nil, err
	}

	return &wallet, nil
}
