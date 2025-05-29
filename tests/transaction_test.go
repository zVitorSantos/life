package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"life/models"
)

// TransactionResponse representa a resposta da API para Transaction
type TransactionResponse struct {
	ID            uint                   `json:"id"`
	WalletID      uint                   `json:"wallet_id"`
	Type          string                 `json:"type"`
	Currency      string                 `json:"currency"`
	Amount        int64                  `json:"amount"`
	Status        string                 `json:"status"`
	Description   string                 `json:"description"`
	BalanceBefore int64                  `json:"balance_before"`
	BalanceAfter  int64                  `json:"balance_after"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// TestTransactionModel testa o modelo Transaction
func TestTransactionModel(t *testing.T) {
	// Teste de criação de transação
	transaction := models.Transaction{
		WalletID:      1,
		Type:          models.TransactionEarn,
		Currency:      models.CurrencyCoins,
		Amount:        500,
		Status:        models.TransactionCompleted,
		Description:   "Test transaction",
		BalanceBefore: 1000,
		BalanceAfter:  1500,
		Metadata:      map[string]interface{}{"source": "test"},
	}

	// Verifica se os campos obrigatórios estão preenchidos
	if transaction.WalletID == 0 {
		t.Error("WalletID não pode ser zero")
	}
	if transaction.Type == "" {
		t.Error("Type não pode ser vazio")
	}
	if transaction.Currency == "" {
		t.Error("Currency não pode ser vazio")
	}
	if transaction.Amount == 0 {
		t.Error("Amount não pode ser zero")
	}
	if transaction.Status == "" {
		t.Error("Status não pode ser vazio")
	}

	// Testa validação de tipos
	validTypes := []models.TransactionType{
		models.TransactionEarn,
		models.TransactionSpend,
		models.TransactionTransfer,
		models.TransactionReward,
		models.TransactionPenalty,
		models.TransactionRefund,
	}

	found := false
	for _, validType := range validTypes {
		if transaction.Type == validType {
			found = true
			break
		}
	}
	if !found {
		t.Error("Tipo de transação inválido")
	}

	// Testa validação de status
	validStatuses := []models.TransactionStatus{
		models.TransactionPending,
		models.TransactionCompleted,
		models.TransactionFailed,
		models.TransactionCancelled,
		models.TransactionReversed,
	}

	found = false
	for _, validStatus := range validStatuses {
		if transaction.Status == validStatus {
			found = true
			break
		}
	}
	if !found {
		t.Error("Status de transação inválido")
	}

	// Testa métodos de reversão
	if !transaction.CanBeReversed() {
		t.Error("Transação completada deveria poder ser revertida")
	}

	// Testa métodos de status
	if !transaction.IsCompleted() {
		t.Error("Transação deveria estar marcada como completada")
	}

	// Testa métodos de metadados
	transaction.SetMetadata("test_key", "test_value")
	if transaction.GetMetadata("test_key") != "test_value" {
		t.Error("Falha ao definir/obter metadado")
	}
}

// TestTransactionFlow testa o fluxo completo de Transaction
func TestTransactionFlow(t *testing.T) {
	setupTest(t)

	// 1. Setup: Registro, Login, perfil e carteira
	user := testRegister(t)
	if user == nil {
		t.Fatal("Falha no registro")
	}

	loginData := testLogin(t, user.Username, "senha123")
	if loginData == nil {
		t.Fatal("Falha no login")
	}

	profile := testCreateGameProfile(t, loginData.AccessToken)
	if profile == nil {
		t.Fatal("Falha ao criar perfil de jogo")
	}

	wallet := testCreateWallet(t, loginData.AccessToken)
	if wallet == nil {
		t.Fatal("Falha ao criar carteira")
	}

	// 2. Adicionar dinheiro
	addTransaction := testAddMoney(t, loginData.AccessToken, "coins", 1000, "Test add money")
	if addTransaction == nil {
		t.Fatal("Falha ao adicionar dinheiro")
	}

	// 3. Gastar dinheiro
	spendTransaction := testSpendMoney(t, loginData.AccessToken, "coins", 500, "Test spend money")
	if spendTransaction == nil {
		t.Fatal("Falha ao gastar dinheiro")
	}

	// 4. Obter histórico de transações
	history := testGetTransactionHistory(t, loginData.AccessToken)
	if history == nil {
		t.Fatal("Falha ao obter histórico de transações")
	}

	// 5. Obter transação específica
	if len(history) > 0 {
		transaction := testGetTransaction(t, loginData.AccessToken, fmt.Sprintf("%d", history[0]["id"]))
		if transaction == nil {
			t.Fatal("Falha ao obter transação específica")
		}
	}
}

// testAddMoney testa a adição de dinheiro
func testAddMoney(t *testing.T, accessToken string, currency string, amount int64, description string) *TransactionResponse {
	url := fmt.Sprintf("%s/transactions/add", baseURL)

	data := map[string]interface{}{
		"currency":    currency,
		"amount":      amount,
		"description": description,
		"metadata":    map[string]interface{}{"source": "test"},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusCreated, resp.StatusCode, string(body))
		return nil
	}

	var transaction TransactionResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&transaction); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &transaction
}

// testSpendMoney testa o gasto de dinheiro
func testSpendMoney(t *testing.T, accessToken string, currency string, amount int64, description string) *TransactionResponse {
	url := fmt.Sprintf("%s/transactions/spend", baseURL)

	data := map[string]interface{}{
		"currency":    currency,
		"amount":      amount,
		"description": description,
		"metadata":    map[string]interface{}{"purpose": "test"},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusCreated, resp.StatusCode, string(body))
		return nil
	}

	var transaction TransactionResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&transaction); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &transaction
}

// testGetTransactionHistory testa a obtenção do histórico de transações
func testGetTransactionHistory(t *testing.T, accessToken string) []map[string]interface{} {
	url := fmt.Sprintf("%s/transactions/history?limit=10", baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusOK, resp.StatusCode, string(body))
		return nil
	}

	var history []map[string]interface{}
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&history); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return history
}

// testGetTransaction testa a obtenção de uma transação específica
func testGetTransaction(t *testing.T, accessToken string, transactionID string) *TransactionResponse {
	url := fmt.Sprintf("%s/transactions/%s", baseURL, transactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusOK, resp.StatusCode, string(body))
		return nil
	}

	var transaction TransactionResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&transaction); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &transaction
}
