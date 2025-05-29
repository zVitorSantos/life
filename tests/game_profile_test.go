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

// GameProfileResponse representa a resposta da API para GameProfile
type GameProfileResponse struct {
	ID        uint                   `json:"id"`
	UserID    uint                   `json:"user_id"`
	Level     int                    `json:"level"`
	XP        int                    `json:"xp"`
	Stats     map[string]interface{} `json:"stats"`
	Settings  map[string]interface{} `json:"settings"`
	IsActive  bool                   `json:"is_active"`
	LastLogin time.Time              `json:"last_login"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// TestGameProfileModel testa o modelo GameProfile
func TestGameProfileModel(t *testing.T) {
	// Teste de criação de perfil de jogo
	profile := models.GameProfile{
		UserID:   1,
		Level:    1,
		XP:       0,
		Stats:    map[string]interface{}{"games_played": 0, "wins": 0},
		Settings: map[string]interface{}{"theme": "dark", "sound": true},
		IsActive: true,
	}

	// Verifica se os campos obrigatórios estão preenchidos
	if profile.UserID == 0 {
		t.Error("UserID não pode ser zero")
	}
	if profile.Level < 1 {
		t.Error("Level deve ser pelo menos 1")
	}
	if profile.XP < 0 {
		t.Error("XP não pode ser negativo")
	}

	// Testa cálculo de XP necessário para próximo level
	nextLevelXP := profile.GetXPForNextLevel()
	expectedXP := int64(profile.Level * 1000)
	if nextLevelXP != expectedXP {
		t.Errorf("XP para próximo level esperado %d, recebido %d", expectedXP, nextLevelXP)
	}

	// Testa adição de XP
	initialLevel := profile.Level
	profile.AddXP(1500) // Suficiente para subir de level
	if profile.Level <= initialLevel {
		t.Errorf("Level deveria ter aumentado após adicionar XP")
	}

	// Testa métodos de estatísticas
	profile.SetStat("test_stat", 100)
	if profile.GetStat("test_stat") != 100 {
		t.Error("Falha ao definir/obter estatística")
	}

	// Testa métodos de configurações
	profile.SetSetting("test_setting", "value")
	if profile.GetSetting("test_setting") != "value" {
		t.Error("Falha ao definir/obter configuração")
	}
}

// TestGameProfileFlow testa o fluxo completo do GameProfile
func TestGameProfileFlow(t *testing.T) {
	setupTest(t)

	// 1. Registro e Login do usuário
	user := testRegister(t)
	if user == nil {
		t.Fatal("Falha no registro")
	}

	loginData := testLogin(t, user.Username, "senha123")
	if loginData == nil {
		t.Fatal("Falha no login")
	}

	// 2. Criar perfil de jogo
	profile := testCreateGameProfile(t, loginData.AccessToken)
	if profile == nil {
		t.Fatal("Falha ao criar perfil de jogo")
	}

	// 3. Obter perfil de jogo
	retrievedProfile := testGetGameProfile(t, loginData.AccessToken)
	if retrievedProfile == nil {
		t.Fatal("Falha ao obter perfil de jogo")
	}

	// 4. Adicionar XP
	updatedProfile := testAddXP(t, loginData.AccessToken, 500)
	if updatedProfile == nil {
		t.Fatal("Falha ao adicionar XP")
	}

	// 5. Atualizar estatísticas
	statsProfile := testUpdateStats(t, loginData.AccessToken)
	if statsProfile == nil {
		t.Fatal("Falha ao atualizar estatísticas")
	}

	// 6. Atualizar último login
	if !testUpdateLastLogin(t, loginData.AccessToken) {
		t.Fatal("Falha ao atualizar último login")
	}

	// 7. Obter leaderboard
	leaderboard := testGetLeaderboard(t, loginData.AccessToken)
	if leaderboard == nil {
		t.Fatal("Falha ao obter leaderboard")
	}
}

// testCreateGameProfile testa a criação de um perfil de jogo
func testCreateGameProfile(t *testing.T, accessToken string) *GameProfileResponse {
	url := fmt.Sprintf("%s/game-profile", baseURL)

	data := map[string]interface{}{
		"stats":    map[string]interface{}{"games_played": 0, "wins": 0},
		"settings": map[string]interface{}{"theme": "dark", "sound": true},
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

	var profile GameProfileResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&profile); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &profile
}

// testGetGameProfile testa a obtenção do perfil de jogo
func testGetGameProfile(t *testing.T, accessToken string) *GameProfileResponse {
	url := fmt.Sprintf("%s/game-profile", baseURL)

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

	var profile GameProfileResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&profile); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &profile
}

// testAddXP testa a adição de XP
func testAddXP(t *testing.T, accessToken string, xpAmount int) *GameProfileResponse {
	url := fmt.Sprintf("%s/game-profile/xp", baseURL)

	data := map[string]interface{}{
		"xp":     xpAmount,
		"reason": "Test XP addition",
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusOK, resp.StatusCode, string(body))
		return nil
	}

	var response map[string]interface{}
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&response); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	// Retorna o perfil atualizado
	return testGetGameProfile(t, accessToken)
}

// testUpdateStats testa a atualização de estatísticas
func testUpdateStats(t *testing.T, accessToken string) *GameProfileResponse {
	url := fmt.Sprintf("%s/game-profile/stats", baseURL)

	data := map[string]interface{}{
		"games_played": 5,
		"wins":         3,
		"losses":       2,
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

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Status code: %d", resp.StatusCode)
	t.Logf("Resposta: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code esperado %d, recebido %d. Resposta: %s", http.StatusOK, resp.StatusCode, string(body))
		return nil
	}

	var profile GameProfileResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&profile); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return &profile
}

// testUpdateLastLogin testa a atualização do último login
func testUpdateLastLogin(t *testing.T, accessToken string) bool {
	url := fmt.Sprintf("%s/game-profile/last-login", baseURL)

	req, err := http.NewRequest("PUT", url, nil)
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

// testGetLeaderboard testa a obtenção do leaderboard
func testGetLeaderboard(t *testing.T, accessToken string) []GameProfileResponse {
	url := fmt.Sprintf("%s/leaderboard?limit=10", baseURL)

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

	var leaderboard []GameProfileResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&leaderboard); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
		return nil
	}

	return leaderboard
}
