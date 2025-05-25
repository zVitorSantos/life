package errors

import (
	"errors"
	"net/http"
)

// AppError representa um erro da aplicação
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implementa a interface error
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// New cria um novo AppError
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Erros comuns
var (
	ErrInvalidRequest     = New(http.StatusBadRequest, "requisição inválida", nil)
	ErrUnauthorized       = New(http.StatusUnauthorized, "não autorizado", nil)
	ErrForbidden          = New(http.StatusForbidden, "acesso negado", nil)
	ErrNotFound           = New(http.StatusNotFound, "recurso não encontrado", nil)
	ErrInternalServer     = New(http.StatusInternalServerError, "erro interno do servidor", nil)
	ErrServiceUnavailable = New(http.StatusServiceUnavailable, "serviço indisponível", nil)
)

// IsAppError verifica se um erro é um AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError retorna o AppError de um erro
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return ErrInternalServer
}
