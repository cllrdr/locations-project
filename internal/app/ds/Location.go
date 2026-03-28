package ds

import (
	"database/sql"
)

type Location struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(50);not null" json:"name"`
	Description string         `gorm:"type:text;not null" json:"description"`
	ImagePath   sql.NullString `gorm:"type:text;default:null" json:"image_path,omitempty"`
	VideoPath   sql.NullString `gorm:"type:text;default:null" json:"video_path,omitempty"`
	Players     string         `gorm:"type:varchar(50);not null" json:"players"`
	IsDeleted   bool           `gorm:"type:boolean;not null;default:false" json:"-"`
}