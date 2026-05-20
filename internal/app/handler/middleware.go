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

		// Убираем префикс "Bearer " если есть
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 1. Проверка блэклиста в Redis
		conn, err := redis.Dial("tcp", h.Config.RedisAddr)
		if err == nil {
			defer conn.Close()
			isBlacklisted, _ := redis.String(conn.Do("GET", tokenStr))
			if isBlacklisted == "blacklisted" {
				c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
					Error: "token revoked",
				})
				c.Abort()
				return
			}
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
		isModerator, exists := c.Get("is_moderator")
		if !exists || isModerator != true {
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
	if isModerator == true {
		return true
	}

	return userID == resourceOwnerID
}