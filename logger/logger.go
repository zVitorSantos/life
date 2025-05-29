package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if gin.Mode() == gin.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})
}

// Fatal registra um erro fatal e encerra a aplicação
func Fatal(message string) {
	log.Fatal().Msg(message)
}

// Error registra um erro
func Error(message string) {
	log.Error().Msg(message)
}

// Info registra uma mensagem informativa
func Info(message string) {
	log.Info().Msg(message)
}

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Processa a requisição
		c.Next()

		// Coleta métricas após a requisição
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Cria evento de log
		event := log.Info()
		if statusCode >= 400 {
			event = log.Error()
		}

		// Adiciona campos ao log
		event.
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Str("query", raw).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent())

		if errorMessage != "" {
			event.Str("error", errorMessage)
		}

		// Registra o log
		event.Msg("Request completed")
	}
}
