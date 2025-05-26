package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"life/models"
)

func TestUserModel(t *testing.T) {
	// Teste de criação de usuário
	user := models.User{
		Username:    "testuser",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Password:    "hashedpassword",
	}

	// Verifica se os campos obrigatórios estão preenchidos
	if user.Username == "" {
		t.Error("Username não pode ser vazio")
	}
	if user.DisplayName == "" {
		t.Error("DisplayName não pode ser vazio")
	}
	if user.Email == "" {
		t.Error("Email não pode ser vazio")
	}
	if user.Password == "" {
		t.Error("Password não pode ser vazio")
	}

	// Verifica se os timestamps são definidos após a criação
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt deve ser definido")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("UpdatedAt deve ser definido")
	}
}

// TestUserFlow testa o fluxo completo de usuário
func TestUserFlow(t *testing.T) {
	// 1. Registro e Login
	user := testRegister(t)
	if user == nil {
		t.Fatal("Falha no registro")
	}

	loginData := testLogin(t, user.Username, "senha123")
	if loginData == nil {
		t.Fatal("Falha no login")
	}

	// 2. Obter usuário
	retrievedUser := testGetUser(t, loginData.AccessToken, strconv.FormatUint(uint64(user.ID), 10))
	if retrievedUser == nil {
		t.Fatal("Falha ao obter usuário")
	}

	// 3. Atualizar usuário
	updatedUser := testUpdateUser(t, loginData.AccessToken, strconv.FormatUint(uint64(user.ID), 10))
	if updatedUser == nil {
		t.Fatal("Falha ao atualizar usuário")
	}

	// 4. Listar usuários
	users := testListUsers(t, loginData.AccessToken)
	if users == nil {
		t.Fatal("Falha ao listar usuários")
	}

	// 5. Logout
	if !testLogout(t, loginData.RefreshToken) {
		t.Fatal("Falha no logout")
	}
}

// testGetUser testa a obtenção de um usuário específico
func testGetUser(t *testing.T, accessToken string, userID string) *User {
	url := fmt.Sprintf("%s/api/v1/users/%s", baseURL, userID)

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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d", http.StatusOK, resp.StatusCode)
		return nil
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &user
}

// testUpdateUser testa a atualização de um usuário
func testUpdateUser(t *testing.T, accessToken string, userID string) *User {
	url := fmt.Sprintf("%s/api/v1/users/%s", baseURL, userID)

	timestamp := time.Now().Format("20060102150405")
	data := map[string]string{
		"display_name": "Usuário Teste Atualizado",
		"email":        fmt.Sprintf("test_updated_%s@example.com", timestamp),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Erro ao criar JSON: %v", err)
		return nil
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d", http.StatusOK, resp.StatusCode)
		return nil
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &user
}

// testListUsers testa a listagem de usuários
func testListUsers(t *testing.T, accessToken string) []User {
	url := fmt.Sprintf("%s/api/v1/users", baseURL)

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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d", http.StatusOK, resp.StatusCode)
		return nil
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return users
}
