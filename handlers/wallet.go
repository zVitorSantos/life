package handlers

import (
	"net/http"
	"strconv"

	"life/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WalletHandler gerencia operações relacionadas às carteiras
type WalletHandler struct {
	db *gorm.DB
}

// NewWalletHandler cria uma nova instância do WalletHandler
func NewWalletHandler(db *gorm.DB) *WalletHandler {
	return &WalletHandler{db: db}
}

// GetWallet obtém a carteira do usuário autenticado
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Busca o perfil de jogo primeiro
	var gameProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).First(&gameProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil de jogo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perfil de jogo"})
		return
	}

	var wallet models.Wallet
	if err := h.db.Where("game_profile_id = ?", gameProfile.ID).Preload("GameProfile.User").First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// CreateWallet cria uma carteira para o perfil de jogo do usuário
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Busca o perfil de jogo
	var gameProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).First(&gameProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil de jogo não encontrado. Crie um perfil primeiro."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perfil de jogo"})
		return
	}

	// Verifica se já existe uma carteira
	var existingWallet models.Wallet
	if err := h.db.Where("game_profile_id = ?", gameProfile.ID).First(&existingWallet).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Carteira já existe"})
		return
	}

	// Cria a carteira
	wallet := models.Wallet{
		GameProfileID: gameProfile.ID,
		CoinsBalance:  0,
		GemsBalance:   0,
		TokensBalance: 0,
		IsLocked:      false,
	}

	if err := h.db.Create(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar carteira"})
		return
	}

	// Carrega a carteira criada com as relações
	if err := h.db.Preload("GameProfile.User").First(&wallet, wallet.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar carteira criada"})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// GetBalance obtém o saldo de uma moeda específica
func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	currency := c.Param("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de moeda é obrigatório"})
		return
	}

	// Valida o tipo de moeda
	var currencyType models.CurrencyType
	switch currency {
	case "coins":
		currencyType = models.CurrencyCoins
	case "gems":
		currencyType = models.CurrencyGems
	case "tokens":
		currencyType = models.CurrencyTokens
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de moeda inválido"})
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

	balance := wallet.GetBalance(currencyType)

	c.JSON(http.StatusOK, gin.H{
		"currency": currency,
		"balance":  balance,
	})
}

// GetAllBalances obtém todos os saldos da carteira
func (h *WalletHandler) GetAllBalances(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	balances := gin.H{
		"coins":       wallet.CoinsBalance,
		"gems":        wallet.GemsBalance,
		"tokens":      wallet.TokensBalance,
		"total_value": wallet.GetTotalValue(),
		"is_locked":   wallet.IsLocked,
		"lock_reason": wallet.LockReason,
	}

	c.JSON(http.StatusOK, balances)
}

// LockWallet bloqueia a carteira do usuário
func (h *WalletHandler) LockWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var request struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Motivo é obrigatório"})
		return
	}

	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	wallet.Lock(request.Reason)

	if err := h.db.Save(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao bloquear carteira"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Carteira bloqueada com sucesso",
		"is_locked":   wallet.IsLocked,
		"lock_reason": wallet.LockReason,
	})
}

// UnlockWallet desbloqueia a carteira do usuário
func (h *WalletHandler) UnlockWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	wallet.Unlock()

	if err := h.db.Save(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao desbloquear carteira"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Carteira desbloqueada com sucesso",
		"is_locked": wallet.IsLocked,
	})
}

// GetWalletStatus obtém o status da carteira
func (h *WalletHandler) GetWalletStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	wallet, err := h.getWalletByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	status := gin.H{
		"wallet_id":   wallet.ID,
		"is_locked":   wallet.IsLocked,
		"lock_reason": wallet.LockReason,
		"created_at":  wallet.CreatedAt,
		"updated_at":  wallet.UpdatedAt,
		"balances": gin.H{
			"coins":  wallet.CoinsBalance,
			"gems":   wallet.GemsBalance,
			"tokens": wallet.TokensBalance,
		},
		"total_value": wallet.GetTotalValue(),
	}

	c.JSON(http.StatusOK, status)
}

// getWalletByUserID é um método auxiliar para buscar carteira por user_id
func (h *WalletHandler) getWalletByUserID(userID uint) (*models.Wallet, error) {
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

// GetWalletHistory obtém o histórico de transações da carteira
func (h *WalletHandler) GetWalletHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Parâmetros de paginação
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

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
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carteira não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar carteira"})
		return
	}

	var transactions []models.Transaction
	if err := h.db.Where("wallet_id = ?", wallet.ID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar histórico"})
		return
	}

	// Conta total de transações
	var total int64
	h.db.Model(&models.Transaction{}).Where("wallet_id = ?", wallet.ID).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}
