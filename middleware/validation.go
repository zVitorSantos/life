package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestValidation é um middleware que valida e documenta o formato das requisições
func RequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Exemplos de formatos esperados por rota
		examples := map[string]interface{}{
			"/api/register": map[string]interface{}{
				"name":     "João Silva",
				"email":    "joao@email.com",
				"password": "senha123",
			},
			"/api/login": map[string]interface{}{
				"email":    "joao@email.com",
				"password": "senha123",
			},
			"/api/refresh": map[string]interface{}{
				"refresh_token": "seu_refresh_token_aqui",
			},
			"/api/logout": map[string]interface{}{
				"refresh_token": "seu_refresh_token_aqui",
			},
			"/api/profile": map[string]interface{}{
				"name":     "João Silva",
				"email":    "joao@email.com",
				"password": "nova_senha123",
			},
			"/api/api-keys": map[string]interface{}{
				"name":        "Minha API Key",
				"description": "API Key para integração",
			},
		}

		// Se for uma requisição POST/PUT e não tiver corpo
		if (c.Request.Method == "POST" || c.Request.Method == "PUT") && c.Request.Body == nil {
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

		// Se tiver corpo, valida se é JSON válido
		if c.Request.Body != nil {
			var jsonData interface{}
			if err := json.NewDecoder(c.Request.Body).Decode(&jsonData); err != nil {
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
