package main

import (
	"life/config"
	"life/logger"
	"os"

	"github.com/joho/godotenv"
)

// @title Life Game API
// @version 1.0
// @description API RESTful para o jogo Life
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API Key for authentication

// @tag.name auth
// @tag.description Operações de autenticação

// @tag.name users
// @tag.description Operações de usuário

// @tag.name profile
// @tag.description Operações de perfil

// @tag.name api-keys
// @tag.description Gerenciamento de chaves de API

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Erro ao carregar .env")
	}

	// Inicializa o container
	container, err := config.NewContainer()
	if err != nil {
		logger.Fatal("Erro ao inicializar container: " + err.Error())
	}

	// Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := container.Router.Engine.Run(":" + port); err != nil {
		logger.Fatal("Erro ao iniciar servidor: " + err.Error())
	}
}
