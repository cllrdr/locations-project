package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gomodule/redigo/redis"
	"locations-project/internal/app/ds"
)

// Claims структура для JWT payload
type Claims struct {
	UserID      uint `json:"user_id"`
	IsModerator bool `json:"is_moderator"`
	jwt.RegisteredClaims
}

// AuthMiddleware проверяет JWT токен и блэклист в Redis
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Error: "authorization header required",
			})
			c.Abort()
			return
		}

		// ✅ Убираем "Bearer " если он есть (поддерживаем оба формата)
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Error: "invalid authorization header format",
			})
			c.Abort()
			return
		}

		// 1. Проверка блэклиста в Redis (используем пул соединений)
		conn := h.RedisPool.Get()
		defer conn.Close()

		isBlacklisted, err := redis.String(conn.Do("GET", tokenStr))
		if err == nil && isBlacklisted == "blacklisted" {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Error: "token revoked",
			})
			c.Abort()
			return
		}

		// 2. Парсинг JWT
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Error: "invalid token",
			})
			c.Abort()
			return
		}

		// 3. Сохраняем данные пользователя в контекст
		c.Set("user_id", claims.UserID)
		c.Set("is_moderator", claims.IsModerator)
		c.Next()
	}
}

// RequireModerator middleware для проверки роли модератора
func RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("is_moderator")
		if !exists {
			c.JSON(http.StatusForbidden, ds.ErrorResponse{
				Error: "moderator access required",
			})
			c.Abort()
			return
		}
		isModerator, ok := val.(bool)
		if !ok || !isModerator {
			c.JSON(http.StatusForbidden, ds.ErrorResponse{
				Error: "moderator access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IsOwnerOrModerator helper для проверки владения ресурсом
func (h *Handler) IsOwnerOrModerator(c *gin.Context, resourceOwnerID uint) bool {
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}

	isModerator, _ := c.Get("is_moderator")
	if isMod, ok := isModerator.(bool); ok && isMod {
		return true
	}

	if uid, ok := userID.(uint); ok && uid == resourceOwnerID {
		return true
	}

	return false
}