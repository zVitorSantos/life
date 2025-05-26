package main

import (
	"fmt"
	"life/config"
	_ "life/docs"
	"life/logger"
	"life/routes"
	"os"
	"strings"

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

// requiredEnvVars lista todas as variáveis de ambiente necessárias
var requiredEnvVars = []string{
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
	"JWT_SECRET",
}

func validateEnvVars() error {
	var missing []string

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missing = append(missing, envVar)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("variáveis de ambiente obrigatórias não encontradas: %s", strings.Join(missing, ", "))
	}

	return nil
}

func main() {
	// Carrega variáveis de ambiente do arquivo .env se existir
	if err := godotenv.Load(); err != nil {
		logger.Info("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Valida se todas as variáveis necessárias estão presentes
	if err := validateEnvVars(); err != nil {
		logger.Fatal("Erro ao validar variáveis de ambiente: " + err.Error())
	}

	// Inicializa o container
	container, err := config.NewContainer()
	if err != nil {
		logger.Fatal("Erro ao inicializar container: " + err.Error())
	}

	// Configura o router
	r := routes.SetupRouter(container.DB)

	// Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Erro ao iniciar servidor: " + err.Error())
	}
}
