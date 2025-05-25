package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Register(c *gin.Context) {
	// TODO: Implementar registro
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *UserHandler) Login(c *gin.Context) {
	// TODO: Implementar login
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: Implementar perfil
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// TODO: Implementar atualização de perfil
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
