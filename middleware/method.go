package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MethodNotAllowed é um middleware que retorna uma mensagem mais descritiva quando o método HTTP está incorreto
func MethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define os métodos permitidos para cada rota
		allowedMethods := map[string][]string{
			"/api/v1/register":  {"POST"},
			"/api/v1/login":     {"POST"},
			"/api/v1/refresh":   {"POST"},
			"/api/v1/logout":    {"POST"},
			"/api/v1/profile":   {"GET", "PUT"},
			"/api/v1/users":     {"GET"},
			"/api/v1/users/:id": {"GET", "PUT"},
		}

		// Obtém os métodos permitidos para a rota atual
		methods, exists := allowedMethods[c.Request.URL.Path]
		if !exists {
			// Se a rota não estiver na lista, permite todos os métodos
			c.Next()
			return
		}

		// Verifica se o método atual está permitido
		for _, method := range methods {
			if method == c.Request.Method {
				c.Next()
				return
			}
		}

		// Se o método não estiver permitido, retorna erro
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":           "Método não permitido",
			"message":         "Este endpoint requer um método HTTP diferente. Verifique a documentação em /swagger/index.html para mais detalhes.",
			"path":            c.Request.URL.Path,
			"method":          c.Request.Method,
			"allowed_methods": methods,
		})
		c.Abort()
	}
}
