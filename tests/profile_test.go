package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestProfileFlow testa o fluxo completo de perfil
func TestProfileFlow(t *testing.T) {
	// 1. Registro e Login
	user := testRegister(t)
	if user == nil {
		t.Fatal("Falha no registro")
	}

	loginData := testLogin(t, user.Username, "senha123")
	if loginData == nil {
		t.Fatal("Falha no login")
	}

	// 2. Obter perfil
	profile := testGetProfile(t, loginData.AccessToken)
	if profile == nil {
		t.Fatal("Falha ao obter perfil")
	}

	// 3. Atualizar perfil
	updatedProfile := testUpdateProfile(t, loginData.AccessToken)
	if updatedProfile == nil {
		t.Fatal("Falha ao atualizar perfil")
	}

	// 4. Logout
	if !testLogout(t, loginData.RefreshToken) {
		t.Fatal("Falha no logout")
	}
}

// testGetProfile testa a obtenção do perfil do usuário
func testGetProfile(t *testing.T, accessToken string) *User {
	url := fmt.Sprintf("%s/api/v1/profile", baseURL)

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

// testUpdateProfile testa a atualização do perfil do usuário
func testUpdateProfile(t *testing.T, accessToken string) *User {
	url := fmt.Sprintf("%s/api/v1/profile", baseURL)

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
