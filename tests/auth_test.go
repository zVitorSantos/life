package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// 4. Logout (usa o refresh token original do login)
	if !testLogout(t, loginData.RefreshToken) {
		t.Fatal("Falha no logout")
	}
}

// testRegister testa o registro de um novo usuário
func testRegister(t *testing.T) *User {
	url := fmt.Sprintf("%s/register", baseURL)

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status code esperado %d, recebido %d", http.StatusCreated, resp.StatusCode)
		return nil
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d", http.StatusOK, resp.StatusCode)
		return nil
	}

	var loginData LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d", http.StatusOK, resp.StatusCode)
		return nil
	}

	var loginData LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Erro na requisição: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Status code esperado %d, recebido %d", http.StatusNoContent, resp.StatusCode)
		return false
	}

	return true
}
