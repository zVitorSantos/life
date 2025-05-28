package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestValidation é um middleware que valida e documenta o formato das requisições
func RequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Exemplos de formatos esperados por rota
		examples := map[string]interface{}{
			"/api/register": map[string]interface{}{
				"username":     "test_user",
				"display_name": "Usuário Teste",
				"email":        "test@email.com",
				"password":     "senha123",
			},
			"/api/login": map[string]interface{}{
				"username": "test_user",
				"password": "senha123",
			},
			"/api/refresh": map[string]interface{}{
				"refresh_token": "seu_refresh_token_aqui",
			},
			"/api/logout": map[string]interface{}{
				"refresh_token": "seu_refresh_token_aqui",
			},
			"/api/profile": map[string]interface{}{
				"display_name": "Usuário Teste",
				"email":        "test@email.com",
			},
			"/api/api-keys": map[string]interface{}{
				"name":        "Minha API Key",
				"description": "API Key para integração",
			},
		}

		// Se for uma requisição POST/PUT, valida o corpo
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			// Lê o corpo da requisição
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Erro ao ler requisição",
					"message": "Não foi possível ler o corpo da requisição",
				})
				c.Abort()
				return
			}

			// Restaura o corpo para o handler usar
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// Se não tiver corpo
			if len(body) == 0 {
				example, exists := examples[c.Request.URL.Path]
				if exists {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Dados inválidos",
						"message": "O corpo da requisição é obrigatório",
						"example": example,
						"format":  "application/json",
					})
					c.Abort()
					return
				}
			}

			// Valida se é JSON válido
			var jsonData interface{}
			if err := json.Unmarshal(body, &jsonData); err != nil {
				example, exists := examples[c.Request.URL.Path]
				if exists {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "JSON inválido",
						"message": "O corpo da requisição deve ser um JSON válido",
						"example": example,
						"format":  "application/json",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
