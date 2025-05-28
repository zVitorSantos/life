package tests

import (
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
