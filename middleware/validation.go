package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestValidation é um middleware que valida e documenta o formato das requisições
func RequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Exemplos de formatos esperados por rota
		examples := map[string]interface{}{
			"/register": map[string]interface{}{
				"username":     "joaosilva",
				"display_name": "João Silva",
				"email":        "joao@email.com",
				"password":     "senha123",
			},
			"/login": map[string]interface{}{
				"username": "joaosilva",
				"password": "senha123",
			},
			"/refresh": map[string]interface{}{
				"refresh_token": "seu_refresh_token_aqui",
			},
			"/logout": map[string]interface{}{
				"refresh_token": "seu_refresh_token_aqui",
			},
			"/profile": map[string]interface{}{
				"display_name": "João Silva",
				"email":        "joao@email.com",
				"password":     "nova_senha123",
			},
			"/api-keys": map[string]interface{}{
				"name":        "Minha API Key",
				"description": "API Key para integração",
			},
		}

		// Se for uma requisição POST/PUT e não tiver corpo
		if (c.Request.Method == "POST" || c.Request.Method == "PUT") && c.Request.ContentLength == 0 {
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
		if c.Request.ContentLength > 0 {
			contentType := c.GetHeader("Content-Type")
			if contentType != "application/json" {
				example, exists := examples[c.Request.URL.Path]
				if exists {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Content-Type inválido",
						"message": "O Content-Type deve ser application/json",
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
