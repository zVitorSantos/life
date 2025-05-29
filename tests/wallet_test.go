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

// WalletResponse representa a resposta da API para Wallet
type WalletResponse struct {
	ID            uint      `json:"id"`
	GameProfileID uint      `json:"game_profile_id"`
	Coins         int64     `json:"coins"`
	Gems          int64     `json:"gems"`
	Tokens        int64     `json:"tokens"`
	IsLocked      bool      `json:"is_locked"`
	LockReason    string    `json:"lock_reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BalanceResponse representa a resposta de saldo
type BalanceResponse struct {
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

// AllBalancesResponse representa todas as moedas
type AllBalancesResponse struct {
	Coins  int64 `json:"coins"`
	Gems   int64 `json:"gems"`
	Tokens int64 `json:"tokens"`
	Total  int64 `json:"total_value"`
}

// TestWalletModel testa o modelo Wallet
func TestWalletModel(t *testing.T) {
	// Teste de criação de carteira
	wallet := models.Wallet{
		GameProfileID: 1,
		CoinsBalance:  1000,
		GemsBalance:   50,
		TokensBalance: 10,
		IsLocked:      false,
	}

	// Verifica se os campos obrigatórios estão preenchidos
	if wallet.GameProfileID == 0 {
		t.Error("GameProfileID não pode ser zero")
	}
	if wallet.CoinsBalance < 0 {
		t.Error("CoinsBalance não pode ser negativo")
	}
	if wallet.GemsBalance < 0 {
		t.Error("GemsBalance não pode ser negativo")
	}
	if wallet.TokensBalance < 0 {
		t.Error("TokensBalance não pode ser negativo")
	}

	// Testa cálculo de valor total
	expectedTotal := wallet.CoinsBalance + (wallet.GemsBalance * 100) + (wallet.TokensBalance * 10)
	totalValue := wallet.GetTotalValue()
	if totalValue != expectedTotal {
		t.Errorf("Valor total esperado %d, recebido %d", expectedTotal, totalValue)
	}

	// Testa validação de saldo suficiente
	if !wallet.HasSufficientBalance(models.CurrencyCoins, 500) {
		t.Error("Deveria ter saldo suficiente para 500 coins")
	}
	if wallet.HasSufficientBalance(models.CurrencyCoins, 2000) {
		t.Error("Não deveria ter saldo suficiente para 2000 coins")
	}

	// Testa bloqueio/desbloqueio
	wallet.Lock("Test lock")
	if !wallet.IsLocked {
		t.Error("Carteira deveria estar bloqueada")
	}
	if wallet.LockReason != "Test lock" {
		t.Error("Motivo do bloqueio incorreto")
	}

	wallet.Unlock()
	if wallet.IsLocked {
		t.Error("Carteira deveria estar desbloqueada")
	}
}

// TestWalletFlow testa o fluxo completo da Wallet
func TestWalletFlow(t *testing.T) {
	setupTest(t)

	// 1. Registro, Login e criação de perfil
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

	// 2. Criar carteira
	wallet := testCreateWallet(t, loginData.AccessToken)
	if wallet == nil {
		t.Fatal("Falha ao criar carteira")
	}

	// 3. Obter carteira
	retrievedWallet := testGetWallet(t, loginData.AccessToken)
	if retrievedWallet == nil {
		t.Fatal("Falha ao obter carteira")
	}

	// 4. Verificar saldos
	coinsBalance := testGetBalance(t, loginData.AccessToken, "coins")
	if coinsBalance == nil {
		t.Fatal("Falha ao obter saldo de coins")
	}

	allBalances := testGetAllBalances(t, loginData.AccessToken)
	if allBalances == nil {
		t.Fatal("Falha ao obter todos os saldos")
	}

	// 5. Bloquear carteira
	if !testLockWallet(t, loginData.AccessToken, "Test lock") {
		t.Fatal("Falha ao bloquear carteira")
	}

	// 6. Verificar status
	status := testGetWalletStatus(t, loginData.AccessToken)
	if status == nil {
		t.Fatal("Falha ao obter status da carteira")
	}

	// 7. Desbloquear carteira
	if !testUnlockWallet(t, loginData.AccessToken) {
		t.Fatal("Falha ao desbloquear carteira")
	}

	// 8. Obter histórico
	history := testGetWalletHistory(t, loginData.AccessToken)
	if history == nil {
		t.Fatal("Falha ao obter histórico da carteira")
	}
}

// testCreateWallet testa a criação de uma carteira
func testCreateWallet(t *testing.T, accessToken string) *WalletResponse {
	url := fmt.Sprintf("%s/wallet", baseURL)

	data := map[string]interface{}{
		"initial_coins":  1000,
		"initial_gems":   50,
		"initial_tokens": 10,
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

	var wallet WalletResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&wallet); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &wallet
}

// testGetWallet testa a obtenção da carteira
func testGetWallet(t *testing.T, accessToken string) *WalletResponse {
	url := fmt.Sprintf("%s/wallet", baseURL)

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

	var wallet WalletResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&wallet); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &wallet
}

// testGetBalance testa a obtenção de saldo específico
func testGetBalance(t *testing.T, accessToken string, currency string) *BalanceResponse {
	url := fmt.Sprintf("%s/wallet/balance/%s", baseURL, currency)

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

	var balance BalanceResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&balance); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &balance
}

// testGetAllBalances testa a obtenção de todos os saldos
func testGetAllBalances(t *testing.T, accessToken string) *AllBalancesResponse {
	url := fmt.Sprintf("%s/wallet/balances", baseURL)

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

	var balances AllBalancesResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&balances); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &balances
}

// testLockWallet testa o bloqueio da carteira
func testLockWallet(t *testing.T, accessToken string, reason string) bool {
	url := fmt.Sprintf("%s/wallet/lock", baseURL)

	data := map[string]string{
		"reason": reason,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	return resp.StatusCode == http.StatusOK
}

// testUnlockWallet testa o desbloqueio da carteira
func testUnlockWallet(t *testing.T, accessToken string) bool {
	url := fmt.Sprintf("%s/wallet/unlock", baseURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return false
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	return resp.StatusCode == http.StatusOK
}

// testGetWalletStatus testa a obtenção do status da carteira
func testGetWalletStatus(t *testing.T, accessToken string) map[string]interface{} {
	url := fmt.Sprintf("%s/wallet/status", baseURL)

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

	var status map[string]interface{}
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&status); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return status
}

// testGetWalletHistory testa a obtenção do histórico da carteira
func testGetWalletHistory(t *testing.T, accessToken string) []map[string]interface{} {
	url := fmt.Sprintf("%s/wallet/history?limit=10", baseURL)

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
