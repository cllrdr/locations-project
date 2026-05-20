package handler

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateToken создаёт JWT токен с user_id и is_moderator
// Важно: функция с маленькой буквы — доступна только внутри пакета handler
func (h *Handler) generateToken(userID uint, isModerator bool) (string, error) {
	if h.Config.JWTSecret == "" {
		return "", fmt.Errorf("JWT secret is empty, check config.toml or .env")
	}

	// Используем ту же структуру Claims, что и в middleware.go
	claims := Claims{
		UserID:      userID,
		IsModerator: isModerator,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.Config.JWTSecret))
}