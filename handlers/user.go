package handlers

import (
	"life/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler gerencia as operações de usuário
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler cria uma nova instância do UserHandler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// Register registra um novo usuário
// @Summary Registra um novo usuário
// @Description Cria uma nova conta de usuário
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Dados do usuário"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Verifica se o usuário já existe
	var existingUser models.User
	if err := h.db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Usuário ou email já existe"})
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}
	user.Password = string(hashedPassword)

	// Cria o usuário
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	// Remove a senha do response
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// GetProfile retorna o perfil do usuário autenticado
// @Summary Obtém perfil do usuário
// @Description Retorna os dados do perfil do usuário autenticado
// @Tags profile
// @Security Bearer
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Remove a senha do response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateProfile atualiza o perfil do usuário autenticado
// @Summary Atualiza perfil do usuário
// @Description Atualiza os dados do perfil do usuário autenticado
// @Tags profile
// @Security Bearer
// @Accept json
// @Produce json
// @Param user body models.User true "Dados do usuário"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Atualiza apenas campos permitidos
	user.DisplayName = updateData.DisplayName
	user.Email = updateData.Email

	// Se a senha foi fornecida, atualiza
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}

	// Remove a senha do response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}
