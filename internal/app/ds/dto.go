package ds

import (
	"time"
)

// RegisterRequest структура для регистрации
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"Иван Петров"`
}

// LoginRequest структура для входа
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// UserResponse безопасная структура пользователя для ответа
type UserResponse struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	IsModerator bool   `json:"is_moderator"`
}

// LoginResponse ответ при успешном входе/регистрации
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse ответ с сообщением
type MessageResponse struct {
	Message string `json:"message"`
}

// GameResponse ответ с информацией о заявке
type GameResponse struct {
	ID            uint                `json:"id"`
	Nickname      string              `json:"nickname"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"created_at"`
	FormedAt      *time.Time          `json:"formed_at,omitempty"`
	CompletedAt   *time.Time          `json:"completed_at,omitempty"`
	CreatorName   string              `json:"creator_name,omitempty"`
	ModeratorName string              `json:"moderator_name,omitempty"`
	RandomPool    int                 `json:"random_pool,omitempty"`
	Locations     []LocationResponse  `json:"locations,omitempty"`
}

// LocationResponse ответ с информацией о локации в заявке
type LocationResponse struct {
	ID            uint   `json:"id"`
	LocationID    uint   `json:"location_id"`
	Priority      int    `json:"priority"`
	LocationName  string `json:"location_name"`
	LocationImage string `json:"location_image"`
	IsRandomed    bool   `json:"is_randomed"`
}

// GamesListResponse список заявок
type GamesListResponse struct {
	ID            uint       `json:"id"`
	Nickname      string     `json:"nickname"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	FormedAt      *time.Time `json:"formed_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatorName   string     `json:"creator_name,omitempty"`
	ModeratorName string     `json:"moderator_name,omitempty"`
	RandomPool    int        `json:"random_pool,omitempty"`
}

// DraftInfoResponse информация о черновике
type DraftInfoResponse struct {
	DraftID      uint `json:"draft_id"`
	LocationsCnt int  `json:"locations_cnt"`
}

// MessageWithIDResponse ответ с сообщением и ID
type MessageWithIDResponse struct {
	Message   string `json:"message"`
	RequestID uint   `json:"request_id,omitempty"`
}