package main

import (
	"log"
	"os"

	"life/config"
	"life/logger"
	"life/routes"

	_ "life/docs"

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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Inicializa o logger
	logger.InitLogger()

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configura o router
	r := routes.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
