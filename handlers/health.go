package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler gerencia os health checks
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler cria uma nova instância do HealthHandler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse representa a resposta do health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  struct {
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	} `json:"database"`
}

// HealthCheck verifica a saúde da aplicação
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	// Verifica conexão com o banco
	sqlDB, err := h.db.DB()
	if err != nil {
		response.Status = "error"
		response.Database.Status = "error"
		response.Database.Error = err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.Status = "error"
		response.Database.Status = "error"
		response.Database.Error = err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	response.Database.Status = "ok"
	c.JSON(http.StatusOK, response)
}

// ReadinessCheck verifica se a aplicação está pronta para receber tráfego
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	// Verifica conexão com o banco
	sqlDB, err := h.db.DB()
	if err != nil {
		response.Status = "error"
		response.Database.Status = "error"
		response.Database.Error = err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.Status = "error"
		response.Database.Status = "error"
		response.Database.Error = err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	response.Database.Status = "ok"
	c.JSON(http.StatusOK, response)
}

// LivenessCheck verifica se a aplicação está viva
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}
