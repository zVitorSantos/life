package main

import (
	"log"
	"os"

	"life/config"
	"life/handlers"
	"life/logger"
	"life/middleware"

	_ "life/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Life Game API
// @version 1.0
// @description API RESTful para o jogo Life
// @host localhost:8080
// @BasePath /api
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

	// Inicializa handlers
	userHandler := handlers.NewUserHandler(db)

	r := gin.Default()

	// Middleware global
	r.Use(gin.Recovery())
	r.Use(logger.LogRequest())

	// Documentação Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas públicas
	public := r.Group("/api")
	{
		// @Summary Registra um novo usuário
		// @Description Cria uma nova conta de usuário
		// @Tags auth
		// @Accept json
		// @Produce json
		// @Success 201 {object} models.User
		// @Router /register [post]
		public.POST("/register", userHandler.Register)

		// @Summary Realiza login
		// @Description Autentica um usuário e retorna um token JWT
		// @Tags auth
		// @Accept json
		// @Produce json
		// @Success 200 {string} string "JWT token"
		// @Router /login [post]
		public.POST("/login", userHandler.Login)
	}

	// Rotas protegidas
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// @Summary Obtém perfil do usuário
		// @Description Retorna os dados do perfil do usuário autenticado
		// @Tags profile
		// @Security Bearer
		// @Produce json
		// @Success 200 {object} models.User
		// @Router /profile [get]
		protected.GET("/profile", userHandler.GetProfile)

		// @Summary Atualiza perfil do usuário
		// @Description Atualiza os dados do perfil do usuário autenticado
		// @Tags profile
		// @Security Bearer
		// @Accept json
		// @Produce json
		// @Success 200 {object} models.User
		// @Router /profile [put]
		protected.PUT("/profile", userHandler.UpdateProfile)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
