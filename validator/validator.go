package validator

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// Erros de validação
	ErrInvalidUsername     = errors.New("nome de usuário inválido")
	ErrInvalidPassword     = errors.New("senha inválida")
	ErrInvalidEmail        = errors.New("email inválido")
	ErrInvalidDisplayName  = errors.New("nome de exibição inválido")
	ErrPasswordTooShort    = errors.New("senha muito curta")
	ErrPasswordTooLong     = errors.New("senha muito longa")
	ErrPasswordNoNumber    = errors.New("senha deve conter pelo menos um número")
	ErrPasswordNoSpecial   = errors.New("senha deve conter pelo menos um caractere especial")
	ErrPasswordNoUppercase = errors.New("senha deve conter pelo menos uma letra maiúscula")
	ErrPasswordNoLowercase = errors.New("senha deve conter pelo menos uma letra minúscula")
)

// ValidateUsername valida um nome de usuário
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return ErrInvalidUsername
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	if err != nil || !matched {
		return ErrInvalidUsername
	}

	return nil
}

// ValidatePassword valida uma senha
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if len(password) > 100 {
		return ErrPasswordTooLong
	}

	if !strings.ContainsAny(password, "0123456789") {
		return ErrPasswordNoNumber
	}

	if !strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?") {
		return ErrPasswordNoSpecial
	}

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return ErrPasswordNoUppercase
	}

	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return ErrPasswordNoLowercase
	}

	return nil
}

// ValidateEmail valida um email
func ValidateEmail(email string) error {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if err != nil || !matched {
		return ErrInvalidEmail
	}

	return nil
}

// ValidateDisplayName valida um nome de exibição
func ValidateDisplayName(name string) error {
	if len(name) < 2 || len(name) > 50 {
		return ErrInvalidDisplayName
	}

	return nil
}
