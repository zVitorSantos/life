package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken gera um novo token JWT
func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
