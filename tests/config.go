package tests

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

// baseURL é a URL base da API
var baseURL = "http://localhost:8080/api/v1"

// setupTest configura o ambiente de teste
func setupTest(t *testing.T) {
	// Verifica se a API está rodando
	if os.Getenv("API_URL") != "" {
		baseURL = os.Getenv("API_URL")
	}

	// Configura o timeout para as requisições
	http.DefaultClient.Timeout = 10 * time.Second

	// Inicia a API se não estiver rodando
	if err := startAPI(); err != nil {
		t.Fatalf("Erro ao iniciar API: %v", err)
	}

	// Aguarda a API iniciar
	time.Sleep(5 * time.Second)
}

// startAPI inicia a API em background
func startAPI() error {
	// Configura as variáveis de ambiente
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "life_test")
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("PORT", "8080")

	// Inicia a API em background
	cmd := exec.Command("go", "run", "main.go")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar API: %v", err)
	}

	return nil
}

// TestMain é a função principal de teste
func TestMain(m *testing.M) {
	// Executa os testes
	code := m.Run()

	// Limpa o ambiente após os testes
	os.Exit(code)
}
