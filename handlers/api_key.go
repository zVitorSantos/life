package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"life/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIKeyHandler struct {
	db *gorm.DB
}

func NewAPIKeyHandler(db *gorm.DB) *APIKeyHandler {
	return &APIKeyHandler{db: db}
}

// generateAPIKey gera uma nova chave de API segura
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateAPIKey cria uma nova chave de API
// @Summary Cria uma nova chave de API
// @Description Cria uma nova chave de API para o usuário autenticado
// @Tags api-keys
// @Accept json
// @Produce json
// @Param apiKey body models.APIKey true "Dados da chave de API"
// @Success 201 {object} models.APIKey
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api-keys [post]
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	var apiKey models.APIKey

	if err := c.ShouldBindJSON(&apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gera uma nova chave
	key, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar chave de API"})
		return
	}

	apiKey.Key = key
	apiKey.UserID = userID
	apiKey.CreatedAt = time.Now()
	apiKey.UpdatedAt = time.Now()

	if err := h.db.Create(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar chave de API"})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}

// ListAPIKeys lista todas as chaves de API do usuário
// @Summary Lista chaves de API
// @Description Retorna todas as chaves de API do usuário autenticado
// @Tags api-keys
// @Produce json
// @Success 200 {array} models.APIKey
// @Failure 401 {object} map[string]string
// @Router /api-keys [get]
func (h *APIKeyHandler) ListAPIKeys(c *gin.Context) {
	userID := c.GetUint("user_id")
	var apiKeys []models.APIKey

	if err := h.db.Where("user_id = ?", userID).Find(&apiKeys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar chaves de API"})
		return
	}

	c.JSON(http.StatusOK, apiKeys)
}

// DeleteAPIKey remove uma chave de API
// @Summary Remove chave de API
// @Description Remove uma chave de API específica
// @Tags api-keys
// @Param id path int true "ID da chave de API"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api-keys/{id} [delete]
func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.APIKey{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar chave de API"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chave de API não encontrada"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateAPIKey atualiza uma chave de API
// @Summary Atualiza chave de API
// @Description Atualiza os dados de uma chave de API específica
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path int true "ID da chave de API"
// @Param apiKey body models.APIKey true "Dados da chave de API"
// @Success 200 {object} models.APIKey
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api-keys/{id} [put]
func (h *APIKeyHandler) UpdateAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var apiKey models.APIKey

	if err := c.ShouldBindJSON(&apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Model(&models.APIKey{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"name":       apiKey.Name,
			"expires_at": apiKey.ExpiresAt,
			"rate_limit": apiKey.RateLimit,
			"is_active":  apiKey.IsActive,
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar chave de API"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chave de API não encontrada"})
		return
	}

	c.JSON(http.StatusOK, apiKey)
}
