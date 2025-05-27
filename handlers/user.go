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
		// Adiciona mais detalhes ao erro de validação
		validationErrors := make(map[string]string)
		if err.Error() == "EOF" {
			validationErrors["body"] = "O corpo da requisição é obrigatório"
		} else {
			validationErrors["validation"] = err.Error()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": validationErrors,
		})
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

// UpdateProfileData representa os dados para atualização de perfil
type UpdateProfileData struct {
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

// UpdateProfile atualiza o perfil do usuário autenticado
// @Summary Atualiza perfil do usuário
// @Description Atualiza os dados do perfil do usuário autenticado
// @Tags profile
// @Security Bearer
// @Accept json
// @Produce json
// @Param user body handlers.UpdateProfileData true "Dados do usuário"
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

	var updateData UpdateProfileData
	if err := c.ShouldBindJSON(&updateData); err != nil {
		validationErrors := make(map[string]string)
		if err.Error() == "EOF" {
			validationErrors["body"] = "O corpo da requisição é obrigatório"
		} else {
			validationErrors["validation"] = err.Error()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": validationErrors,
		})
		return
	}

	// Atualiza apenas campos permitidos
	user.DisplayName = updateData.DisplayName
	user.Email = updateData.Email

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}

	// Remove a senha do response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// GetUser retorna um usuário específico
// @Summary Obtém um usuário específico
// @Description Retorna os dados de um usuário específico
// @Tags users
// @Security Bearer
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Remove a senha do response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUserData representa os dados para atualização de usuário
type UpdateUserData struct {
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

// UpdateUser atualiza um usuário específico
// @Summary Atualiza um usuário específico
// @Description Atualiza os dados de um usuário específico
// @Tags users
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param user body handlers.UpdateUserData true "Dados do usuário"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	var updateData UpdateUserData
	if err := c.ShouldBindJSON(&updateData); err != nil {
		validationErrors := make(map[string]string)
		if err.Error() == "EOF" {
			validationErrors["body"] = "O corpo da requisição é obrigatório"
		} else {
			validationErrors["validation"] = err.Error()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": validationErrors,
		})
		return
	}

	// Atualiza apenas campos permitidos
	user.DisplayName = updateData.DisplayName
	user.Email = updateData.Email

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}

	// Remove a senha do response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// ListUsers retorna todos os usuários
// @Summary Lista todos os usuários
// @Description Retorna uma lista de todos os usuários
// @Tags users
// @Security Bearer
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar usuários"})
		return
	}

	// Remove a senha de todos os usuários
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}
