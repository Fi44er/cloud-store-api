package service

import (
	"time"

	"github.com/gloude/store/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// Claims — кастомные JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// generateJWT генерирует JWT токен для пользователя
// Документация: https://github.com/golang-jwt/jwt
func generateJWT(user *models.User, secret string) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
