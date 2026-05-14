package ds

import (
	"time"
)

type GameStatus string

const (
	GameStatusDraft     GameStatus = "черновик"
	GameStatusDeleted   GameStatus = "удалён"
	GameStatusFormed    GameStatus = "сформирован"
	GameStatusCompleted GameStatus = "завершён"
	GameStatusRejected  GameStatus = "отклонён"
)

type PlayersLocationGame struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nickname    string     `gorm:"type:varchar(100);not null" json:"nickname"`
	Status      GameStatus `gorm:"type:varchar(20);not null;check:status IN ('черновик','удалён','сформирован','завершён','отклонён')" json:"status"`
	CreatedAt   time.Time  `gorm:"not null" json:"created_at"`
	FormedAt    *time.Time `gorm:"default:null" json:"formed_at,omitempty"`
	CompletedAt *time.Time `gorm:"default:null" json:"completed_at,omitempty"`
	CreatorID   uint       `gorm:"not null" json:"creator_id"`
	ModeratorID *uint      `gorm:"default:null" json:"moderator_id,omitempty"`

	Creator   User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Moderator User `gorm:"foreignKey:ModeratorID" json:"moderator,omitempty"`
}
