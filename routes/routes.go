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
	gameProfileHandler := handlers.NewGameProfileHandler(db)
	walletHandler := handlers.NewWalletHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)

	// Middleware global
	r.Use(gin.Recovery())
	r.Use(logger.LogRequest())

	// Documentação Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health checks
	setupHealthRoutes(r, healthHandler)

	// Rotas públicas
	public := r.Group("/api/v1")
	{
		setupPublicRoutes(public, userHandler, authHandler)
	}

	// Rotas protegidas por JWT
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		setupProtectedRoutes(protected, userHandler, apiKeyHandler, gameProfileHandler, walletHandler, transactionHandler)
	}

	// Rotas protegidas por API Key
	apiProtected := r.Group("/api/v1")
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
	router.Use(middleware.MethodNotAllowed())
	router.Use(middleware.RequestValidation())
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
		router.POST("/register", userHandler.Register)

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
		router.POST("/login", authHandler.Login)

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
		router.POST("/refresh", authHandler.Refresh)

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
		router.POST("/logout", authHandler.Logout)
	}
}

// setupProtectedRoutes configura as rotas protegidas por JWT
func setupProtectedRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, apiKeyHandler *handlers.APIKeyHandler, gameProfileHandler *handlers.GameProfileHandler, walletHandler *handlers.WalletHandler, transactionHandler *handlers.TransactionHandler) {
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

	// Rotas de usuário
	router.GET("/users", userHandler.ListUsers)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("/users/:id", userHandler.UpdateUser)

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

	// Rotas de Game Profile
	gameProfile := router.Group("/game-profile")
	{
		// @Summary Cria um novo perfil de jogo
		// @Description Cria um novo perfil de jogo para o usuário autenticado
		// @Tags game-profile
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 201 {object} models.GameProfile
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /game-profile [post]
		gameProfile.POST("", gameProfileHandler.CreateGameProfile)

		// @Summary Obtém perfil de jogo
		// @Description Retorna os dados do perfil de jogo do usuário autenticado
		// @Tags game-profile
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} models.GameProfile
		// @Failure 401 {object} map[string]string
		// @Failure 404 {object} map[string]string
		// @Router /game-profile [get]
		gameProfile.GET("", gameProfileHandler.GetGameProfile)

		// @Summary Atualiza perfil de jogo
		// @Description Atualiza os dados do perfil de jogo do usuário autenticado
		// @Tags game-profile
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} models.GameProfile
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Failure 404 {object} map[string]string
		// @Router /game-profile [put]
		gameProfile.PUT("", gameProfileHandler.UpdateGameProfile)

		// @Summary Adiciona XP ao perfil
		// @Description Adiciona XP ao perfil de jogo do usuário autenticado
		// @Tags game-profile
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /game-profile/xp [post]
		gameProfile.POST("/xp", gameProfileHandler.AddXP)

		// @Summary Obtém estatísticas do perfil
		// @Description Retorna as estatísticas do perfil de jogo
		// @Tags game-profile
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /game-profile/stats [get]
		gameProfile.GET("/stats", gameProfileHandler.GetStats)

		// @Summary Atualiza último login
		// @Description Atualiza o timestamp do último login
		// @Tags game-profile
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /game-profile/last-login [put]
		gameProfile.PUT("/last-login", gameProfileHandler.UpdateLastLogin)
	}

	// @Summary Obtém ranking de jogadores
	// @Description Retorna o ranking dos melhores jogadores
	// @Tags game-profile
	// @Security Bearer
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /leaderboard [get]
	router.GET("/leaderboard", gameProfileHandler.GetLeaderboard)

	// Rotas de Wallet
	wallet := router.Group("/wallet")
	{
		// @Summary Cria uma nova carteira
		// @Description Cria uma nova carteira para o usuário autenticado
		// @Tags wallet
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 201 {object} models.Wallet
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /wallet [post]
		wallet.POST("", walletHandler.CreateWallet)

		// @Summary Obtém carteira
		// @Description Retorna os dados da carteira do usuário autenticado
		// @Tags wallet
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} models.Wallet
		// @Failure 401 {object} map[string]string
		// @Router /wallet [get]
		wallet.GET("", walletHandler.GetWallet)

		// @Summary Obtém saldo específico
		// @Description Retorna o saldo de uma moeda específica
		// @Tags wallet
		// @Security Bearer
		// @Param currency path string true "Tipo de moeda (coins, gems, tokens)"
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /wallet/balance/{currency} [get]
		wallet.GET("/balance/:currency", walletHandler.GetBalance)

		// @Summary Obtém todos os saldos
		// @Description Retorna todos os saldos da carteira
		// @Tags wallet
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /wallet/balances [get]
		wallet.GET("/balances", walletHandler.GetAllBalances)

		// @Summary Bloqueia carteira
		// @Description Bloqueia a carteira do usuário
		// @Tags wallet
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /wallet/lock [post]
		wallet.POST("/lock", walletHandler.LockWallet)

		// @Summary Desbloqueia carteira
		// @Description Desbloqueia a carteira do usuário
		// @Tags wallet
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /wallet/unlock [post]
		wallet.POST("/unlock", walletHandler.UnlockWallet)

		// @Summary Status da carteira
		// @Description Obtém o status completo da carteira
		// @Tags wallet
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /wallet/status [get]
		wallet.GET("/status", walletHandler.GetWalletStatus)

		// @Summary Histórico da carteira
		// @Description Obtém o histórico de transações da carteira
		// @Tags wallet
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /wallet/history [get]
		wallet.GET("/history", walletHandler.GetWalletHistory)
	}

	// Rotas de Transaction
	transactions := router.Group("/transactions")
	{
		// @Summary Adiciona dinheiro
		// @Description Adiciona dinheiro à carteira do usuário
		// @Tags transactions
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /transactions/add [post]
		transactions.POST("/add", transactionHandler.AddMoney)

		// @Summary Gasta dinheiro
		// @Description Remove dinheiro da carteira do usuário
		// @Tags transactions
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /transactions/spend [post]
		transactions.POST("/spend", transactionHandler.SpendMoney)

		// @Summary Transfere dinheiro
		// @Description Transfere dinheiro entre usuários
		// @Tags transactions
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]string
		// @Failure 401 {object} map[string]string
		// @Router /transactions/transfer [post]
		transactions.POST("/transfer", transactionHandler.TransferMoney)

		// @Summary Histórico de transações
		// @Description Obtém o histórico de transações do usuário
		// @Tags transactions
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 401 {object} map[string]string
		// @Router /transactions/history [get]
		transactions.GET("/history", transactionHandler.GetTransactionHistory)

		// @Summary Obtém transação
		// @Description Retorna os dados de uma transação específica
		// @Tags transactions
		// @Security Bearer
		// @Param id path int true "ID da transação"
		// @Success 200 {object} models.Transaction
		// @Failure 401 {object} map[string]string
		// @Failure 404 {object} map[string]string
		// @Router /transactions/{id} [get]
		transactions.GET("/:id", transactionHandler.GetTransaction)
	}
}

// setupAPIProtectedRoutes configura as rotas protegidas por API Key
func setupAPIProtectedRoutes(router *gin.RouterGroup) {
	// Aqui você pode adicionar rotas que podem ser acessadas via API Key
	// Por exemplo, endpoints públicos que precisam de rate limiting
}
