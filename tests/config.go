package tests

import (
	"net/http"
	"os"
	"testing"
	"time"
)

// baseURL é a URL base da API
var baseURL = "http://localhost:8080/api"

// setupTest configura o ambiente de teste
func setupTest(t *testing.T) {
	// Verifica se a API está rodando
	if os.Getenv("API_URL") != "" {
		baseURL = os.Getenv("API_URL")
	}

	// Configura o timeout para as requisições
	http.DefaultClient.Timeout = 10 * time.Second
}

// TestMain é a função principal de teste
func TestMain(m *testing.M) {
	// Executa os testes
	code := m.Run()

	// Limpa o ambiente após os testes
	os.Exit(code)
}
