package ds

import (
	"database/sql"
)

type Location struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(50);not null"`
	Description string         `gorm:"type:text;not null"`
	ImagePath   sql.NullString `gorm:"type:text;default:null"`
	VideoPath   sql.NullString `gorm:"type:text;default:null"`
	Players     string         `gorm:"type:varchar(50);not null"`
}