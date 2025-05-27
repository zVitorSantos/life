package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"life/auth"
	"life/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler gerencia as operações de autenticação
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler cria uma nova instância do AuthHandler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// LoginResponse representa a resposta do login
// @Description Resposta de autenticação contendo os tokens
type LoginResponse struct {
	// Token de acesso JWT
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// Token de atualização
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// Tipo do token (sempre "Bearer")
	TokenType string `json:"token_type" example:"Bearer"`

	// Tempo de expiração em segundos
	ExpiresIn int64 `json:"expires_in" example:"3600"`
}

// Login autentica um usuário e retorna tokens
// @Summary Realiza login
// @Description Autentica um usuário e retorna tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "Credenciais de login"
// @Success 200 {object} handlers.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Gera access token
	accessToken, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// Gera refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar refresh token"})
		return
	}

	// Salva refresh token no banco
	rt := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour * 7), // 7 dias
	}

	if err := h.db.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar refresh token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hora
	})
}

// Refresh atualiza o access token usando o refresh token
// @Summary Atualiza access token
// @Description Atualiza o access token usando o refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body map[string]string true "Refresh token"
// @Success 200 {object} handlers.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var refreshData struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var rt models.RefreshToken
	if err := h.db.Where("token = ? AND is_revoked = ? AND expires_at > ?",
		refreshData.RefreshToken, false, time.Now()).First(&rt).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token inválido ou expirado"})
		return
	}

	// Gera novo access token
	accessToken, err := auth.GenerateToken(rt.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshData.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hora
	})
}

// Logout revoga um refresh token
// @Summary Realiza logout
// @Description Revoga um refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body map[string]string true "Refresh token"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var logoutData struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&logoutData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	result := h.db.Model(&models.RefreshToken{}).
		Where("token = ?", logoutData.RefreshToken).
		Update("is_revoked", true)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao revogar token"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token não encontrado"})
		return
	}

	c.Status(http.StatusNoContent)
}

// generateRefreshToken gera um token de atualização seguro
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
