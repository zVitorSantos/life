package routes

import (
	"life/handlers"
	"life/logger"
	"life/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// SetupRouter configura todas as rotas da aplicação
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Inicializa handlers
	userHandler := handlers.NewUserHandler(db)
	apiKeyHandler := handlers.NewAPIKeyHandler(db)

	// Middleware global
	r.Use(gin.Recovery())
	r.Use(logger.LogRequest())

	// Documentação Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas públicas
	public := r.Group("/api")
	{
		setupPublicRoutes(public, userHandler)
	}

	// Rotas protegidas por JWT
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		setupProtectedRoutes(protected, userHandler, apiKeyHandler)
	}

	// Rotas protegidas por API Key
	apiProtected := r.Group("/api")
	apiProtected.Use(middleware.APIKeyAuth(db))
	{
		setupAPIProtectedRoutes(apiProtected)
	}

	return r
}

// setupPublicRoutes configura as rotas públicas
func setupPublicRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler) {
	// @Summary Registra um novo usuário
	// @Description Cria uma nova conta de usuário
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param user body models.User true "Dados do usuário"
	// @Success 201 {object} models.User
	// @Failure 400 {object} map[string]string
	// @Failure 409 {object} map[string]string
	// @Router /register [post]
	router.POST("/register", userHandler.Register)

	// @Summary Realiza login
	// @Description Autentica um usuário e retorna um token JWT
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param credentials body map[string]string true "Credenciais de login"
	// @Success 200 {object} map[string]string "JWT token"
	// @Failure 400 {object} map[string]string
	// @Failure 401 {object} map[string]string
	// @Router /login [post]
	router.POST("/login", userHandler.Login)
}

// setupProtectedRoutes configura as rotas protegidas por JWT
func setupProtectedRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, apiKeyHandler *handlers.APIKeyHandler) {
	// Rotas de perfil
	// @Summary Obtém perfil do usuário
	// @Description Retorna os dados do perfil do usuário autenticado
	// @Tags profile
	// @Security Bearer
	// @Produce json
	// @Success 200 {object} models.User
	// @Failure 401 {object} map[string]string
	// @Failure 404 {object} map[string]string
	// @Router /profile [get]
	router.GET("/profile", userHandler.GetProfile)

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
	router.PUT("/profile", userHandler.UpdateProfile)

	// Rotas de API Key
	apiKeys := router.Group("/api-keys")
	{
		// @Summary Cria uma nova chave de API
		// @Description Cria uma nova chave de API para o usuário autenticado
		// @Tags api-keys
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Param apiKey body models.APIKey true "Dados da chave de API"
		// @Success 201 {object} models.APIKey
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /api-keys [post]
		apiKeys.POST("", apiKeyHandler.CreateAPIKey)

		// @Summary Lista chaves de API
		// @Description Retorna todas as chaves de API do usuário autenticado
		// @Tags api-keys
		// @Security Bearer
		// @Produce json
		// @Success 200 {array} models.APIKey
		// @Failure 401 {object} map[string]string
		// @Router /api-keys [get]
		apiKeys.GET("", apiKeyHandler.ListAPIKeys)

		// @Summary Remove chave de API
		// @Description Remove uma chave de API específica
		// @Tags api-keys
		// @Security Bearer
		// @Param id path int true "ID da chave de API"
		// @Success 204 "No Content"
		// @Failure 401 {object} map[string]string
		// @Failure 404 {object} map[string]string
		// @Router /api-keys/{id} [delete]
		apiKeys.DELETE("/:id", apiKeyHandler.DeleteAPIKey)

		// @Summary Atualiza chave de API
		// @Description Atualiza os dados de uma chave de API específica
		// @Tags api-keys
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Param id path int true "ID da chave de API"
		// @Param apiKey body models.APIKey true "Dados da chave de API"
		// @Success 200 {object} models.APIKey
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Failure 404 {object} map[string]string
		// @Router /api-keys/{id} [put]
		apiKeys.PUT("/:id", apiKeyHandler.UpdateAPIKey)
	}
}

// setupAPIProtectedRoutes configura as rotas protegidas por API Key
func setupAPIProtectedRoutes(router *gin.RouterGroup) {
	// Aqui você pode adicionar rotas que podem ser acessadas via API Key
	// Por exemplo, endpoints públicos que precisam de rate limiting
}
