package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MethodNotAllowed é um middleware que retorna uma mensagem mais descritiva quando o método HTTP está incorreto
func MethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":           "Método não permitido",
			"message":         "Este endpoint requer um método HTTP diferente. Verifique a documentação em /swagger/index.html para mais detalhes.",
			"path":            c.Request.URL.Path,
			"method":          c.Request.Method,
			"allowed_methods": []string{"POST"}, // Você pode personalizar isso baseado na rota
		})
		c.Abort()
	}
}
