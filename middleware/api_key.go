package middleware

import (
	"net/http"
	"time"

	"life/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RateLimiter é uma estrutura simples para controle de rate limiting
type RateLimiter struct {
	requests map[string][]time.Time
}

var limiter = &RateLimiter{
	requests: make(map[string][]time.Time),
}

// APIKeyAuth é um middleware para autenticação via API Key
func APIKeyAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key não fornecida"})
			c.Abort()
			return
		}

		var key models.APIKey
		if err := db.Where("key = ? AND is_active = ?", apiKey, true).First(&key).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key inválida"})
			c.Abort()
			return
		}

		// Verifica expiração
		if time.Now().After(key.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key expirada"})
			c.Abort()
			return
		}

		// Rate limiting
		if !checkRateLimit(apiKey, key.RateLimit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Limite de requisições excedido"})
			c.Abort()
			return
		}

		// Atualiza último uso
		now := time.Now()
		db.Model(&key).Update("last_used_at", now)

		// Adiciona informações ao contexto
		c.Set("user_id", key.UserID)
		c.Set("api_key_id", key.ID)

		c.Next()
	}
}

// checkRateLimit verifica se a requisição está dentro do limite
func checkRateLimit(key string, limit int) bool {
	now := time.Now()
	window := now.Add(-time.Minute)

	// Limpa requisições antigas
	var validRequests []time.Time
	for _, t := range limiter.requests[key] {
		if t.After(window) {
			validRequests = append(validRequests, t)
		}
	}
	limiter.requests[key] = validRequests

	// Verifica limite
	if len(validRequests) >= limit {
		return false
	}

	// Adiciona nova requisição
	limiter.requests[key] = append(limiter.requests[key], now)
	return true
}
