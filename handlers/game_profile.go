package handlers

import (
	"net/http"
	"strconv"
	"time"

	"life/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GameProfileHandler gerencia operações relacionadas aos perfis de jogo
type GameProfileHandler struct {
	db *gorm.DB
}

// NewGameProfileHandler cria uma nova instância do GameProfileHandler
func NewGameProfileHandler(db *gorm.DB) *GameProfileHandler {
	return &GameProfileHandler{db: db}
}

// GetGameProfile obtém o perfil de jogo do usuário autenticado
func (h *GameProfileHandler) GetGameProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var gameProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).Preload("User").First(&gameProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil de jogo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perfil de jogo"})
		return
	}

	c.JSON(http.StatusOK, gameProfile)
}

// CreateGameProfile cria um perfil de jogo para o usuário
func (h *GameProfileHandler) CreateGameProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Verifica se já existe um perfil de jogo
	var existingProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).First(&existingProfile).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Perfil de jogo já existe"})
		return
	}

	// Cria o perfil de jogo
	gameProfile := models.GameProfile{
		UserID:    userID.(uint),
		Level:     1,
		XP:        0,
		IsActive:  true,
		LastLogin: nil,
		Stats:     make(map[string]interface{}),
		Settings:  make(map[string]interface{}),
	}

	if err := h.db.Create(&gameProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar perfil de jogo"})
		return
	}

	// Retorna o perfil criado sem carregar a relação User
	c.JSON(http.StatusCreated, gameProfile)
}

// UpdateGameProfile atualiza o perfil de jogo do usuário
func (h *GameProfileHandler) UpdateGameProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var gameProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).First(&gameProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil de jogo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perfil de jogo"})
		return
	}

	var updateData struct {
		Settings map[string]interface{} `json:"settings"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Atualiza apenas as configurações (por segurança)
	if updateData.Settings != nil {
		gameProfile.Settings = updateData.Settings
	}

	if err := h.db.Save(&gameProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar perfil de jogo"})
		return
	}

	// Retorna o perfil atualizado sem carregar a relação User
	c.JSON(http.StatusOK, gameProfile)
}

// AddXP adiciona XP ao perfil do usuário
func (h *GameProfileHandler) AddXP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var request struct {
		Amount int64  `json:"amount" binding:"required,min=1"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var gameProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).First(&gameProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil de jogo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perfil de jogo"})
		return
	}

	// Salva dados antes da mudança
	oldLevel := gameProfile.Level
	oldXP := gameProfile.XP

	// Adiciona XP
	gameProfile.AddXP(request.Amount)

	// Salva no banco
	if err := h.db.Save(&gameProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar XP"})
		return
	}

	// Resposta com informações da mudança
	response := gin.H{
		"old_level":     oldLevel,
		"new_level":     gameProfile.Level,
		"old_xp":        oldXP,
		"new_xp":        gameProfile.XP,
		"xp_added":      request.Amount,
		"level_up":      gameProfile.Level > oldLevel,
		"next_level_xp": gameProfile.GetXPForNextLevel(),
		"progress":      gameProfile.GetXPProgress(),
	}

	if request.Reason != "" {
		response["reason"] = request.Reason
	}

	c.JSON(http.StatusOK, response)
}

// GetStats obtém estatísticas específicas do perfil
func (h *GameProfileHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var gameProfile models.GameProfile
	if err := h.db.Where("user_id = ?", userID).First(&gameProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil de jogo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perfil de jogo"})
		return
	}

	stats := gin.H{
		"level":         gameProfile.Level,
		"xp":            gameProfile.XP,
		"next_level_xp": gameProfile.GetXPForNextLevel(),
		"progress":      gameProfile.GetXPProgress(),
		"is_active":     gameProfile.IsActive,
		"last_login":    gameProfile.LastLogin,
		"custom_stats":  gameProfile.Stats,
	}

	c.JSON(http.StatusOK, stats)
}

// UpdateLastLogin atualiza o timestamp do último login
func (h *GameProfileHandler) UpdateLastLogin(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	now := time.Now()
	if err := h.db.Model(&models.GameProfile{}).Where("user_id = ?", userID).Update("last_login", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar último login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Último login atualizado",
		"last_login": now,
	})
}

// GetLeaderboard obtém o ranking de jogadores
func (h *GameProfileHandler) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	var profiles []models.GameProfile
	if err := h.db.Preload("User").
		Where("is_active = ?", true).
		Order("level DESC, xp DESC").
		Limit(limit).
		Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ranking"})
		return
	}

	// Formata a resposta do leaderboard
	leaderboard := make([]gin.H, len(profiles))
	for i, profile := range profiles {
		leaderboard[i] = gin.H{
			"rank":     i + 1,
			"user_id":  profile.UserID,
			"username": profile.User.Username,
			"level":    profile.Level,
			"xp":       profile.XP,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
		"total":       len(leaderboard),
	})
}
