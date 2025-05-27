package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// User representa um usuário nos testes
type User struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TestAuthFlow testa o fluxo completo de autenticação
func TestAuthFlow(t *testing.T) {
	setupTest(t)
	// 1. Registro
	user := testRegister(t)
	if user == nil {
		t.Fatal("Falha no registro")
	}

	// 2. Login
	loginData := testLogin(t, user.Username, "senha123")
	if loginData == nil {
		t.Fatal("Falha no login")
	}

	// 3. Refresh Token
	newLoginData := testRefreshToken(t, loginData.RefreshToken)
	if newLoginData == nil {
		t.Fatal("Falha no refresh token")
	}

	// 4. Logout
	if !testLogout(t, newLoginData.RefreshToken) {
		t.Fatal("Falha no logout")
	}
}

// testRegister testa o registro de um novo usuário
func testRegister(t *testing.T) *User {
	url := fmt.Sprintf("%s/register", baseURL)

	timestamp := time.Now().Format("20060102150405")
	data := map[string]string{
		"username":     fmt.Sprintf("test_user_%s", timestamp),
		"password":     "senha123",
		"display_name": "Usuário Teste",
		"email":        fmt.Sprintf("test_%s@example.com", timestamp),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return nil
	}

	// Log do corpo da requisição
	t.Logf("Corpo da requisição: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	// Log da resposta
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusCreated, resp.StatusCode, string(body))
		return nil
	}

	var user User
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&user); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &user
}

// testLogin testa o login de um usuário
func testLogin(t *testing.T, username, password string) *LoginResponse {
	url := fmt.Sprintf("%s/login", baseURL)

	data := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return nil
	}

	// Log do corpo da requisição
	t.Logf("Corpo da requisição: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	// Log da resposta
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusOK, resp.StatusCode, string(body))
		return nil
	}

	var loginData LoginResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&loginData); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &loginData
}

// testRefreshToken testa o refresh do token
func testRefreshToken(t *testing.T, refreshToken string) *LoginResponse {
	url := fmt.Sprintf("%s/refresh", baseURL)

	data := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return nil
	}

	// Log do corpo da requisição
	t.Logf("Corpo da requisição: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	// Log da resposta
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusOK, resp.StatusCode, string(body))
		return nil
	}

	var loginData LoginResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&loginData); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &loginData
}

// testLogout testa o logout
func testLogout(t *testing.T, refreshToken string) bool {
	url := fmt.Sprintf("%s/logout", baseURL)

	data := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return false
	}

	// Log do corpo da requisição
	t.Logf("Corpo da requisição: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro ao criar requisição: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Log da resposta
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusNoContent, resp.StatusCode, string(body))
		return false
	}

	return true
}
