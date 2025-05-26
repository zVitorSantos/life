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

// Router gerencia as rotas da aplicação
type Router struct {
	Engine *gin.Engine
	db     *gorm.DB
	user   *handlers.UserHandler
	auth   *handlers.AuthHandler
	apiKey *handlers.APIKeyHandler
	health *handlers.HealthHandler
}

// NewRouter cria uma nova instância do Router
func NewRouter(db *gorm.DB, user *handlers.UserHandler, auth *handlers.AuthHandler, apiKey *handlers.APIKeyHandler, health *handlers.HealthHandler) *Router {
	return &Router{
		Engine: gin.Default(),
		db:     db,
		user:   user,
		auth:   auth,
		apiKey: apiKey,
		health: health,
	}
}

// SetupRouter configura todas as rotas da aplicação
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Inicializa handlers
	userHandler := handlers.NewUserHandler(db)
	apiKeyHandler := handlers.NewAPIKeyHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	healthHandler := handlers.NewHealthHandler(db)

	// Middleware global
	r.Use(gin.Recovery())
	r.Use(logger.LogRequest())

	// Documentação Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health checks
	setupHealthRoutes(r, healthHandler)

	// Rotas públicas
	public := r.Group("/api")
	{
		setupPublicRoutes(public, userHandler, authHandler)
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

// setupHealthRoutes configura as rotas de health check
func setupHealthRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler) {
	// @Summary Health check
	// @Description Verifica a saúde da aplicação
	// @Tags health
	// @Produce json
	// @Success 200 {object} handlers.HealthResponse
	// @Failure 503 {object} handlers.HealthResponse
	// @Router /health [get]
	router.GET("/health", healthHandler.HealthCheck)

	// @Summary Readiness check
	// @Description Verifica se a aplicação está pronta para receber tráfego
	// @Tags health
	// @Produce json
	// @Success 200 {object} handlers.HealthResponse
	// @Failure 503 {object} handlers.HealthResponse
	// @Router /ready [get]
	router.GET("/ready", healthHandler.ReadinessCheck)

	// @Summary Liveness check
	// @Description Verifica se a aplicação está viva
	// @Tags health
	// @Produce json
	// @Success 200 {object} handlers.HealthResponse
	// @Router /live [get]
	router.GET("/live", healthHandler.LivenessCheck)
}

// setupPublicRoutes configura as rotas públicas
func setupPublicRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) {
	// Middleware para rotas de autenticação
	auth := router.Group("")
	auth.Use(middleware.MethodNotAllowed())
	auth.Use(middleware.RequestValidation())
	{
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
		auth.POST("/register", userHandler.Register)

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
		auth.POST("/login", authHandler.Login)

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
		auth.POST("/refresh", authHandler.Refresh)

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
		auth.POST("/logout", authHandler.Logout)
	}
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
