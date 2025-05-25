package config

import (
	"life/handlers"
	"life/logger"
	"life/routes"

	"gorm.io/gorm"
)

// Container gerencia as dependências da aplicação
type Container struct {
	DB            *gorm.DB
	UserHandler   *handlers.UserHandler
	AuthHandler   *handlers.AuthHandler
	APIKeyHandler *handlers.APIKeyHandler
	HealthHandler *handlers.HealthHandler
	Router        *routes.Router
}

// NewContainer cria uma nova instância do container
func NewContainer() (*Container, error) {
	// Inicializa o logger
	logger.InitLogger()

	// Inicializa o banco de dados
	db, err := InitDB()
	if err != nil {
		return nil, err
	}

	// Inicializa os handlers
	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	apiKeyHandler := handlers.NewAPIKeyHandler(db)
	healthHandler := handlers.NewHealthHandler(db)

	// Inicializa o router
	router := routes.NewRouter(db, userHandler, authHandler, apiKeyHandler, healthHandler)

	return &Container{
		DB:            db,
		UserHandler:   userHandler,
		AuthHandler:   authHandler,
		APIKeyHandler: apiKeyHandler,
		HealthHandler: healthHandler,
		Router:        router,
	}, nil
}
