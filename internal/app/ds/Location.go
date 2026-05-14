package ds

type Location struct {
	ID               uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string `gorm:"type:varchar(50);not null" json:"name"`
	Description      string `gorm:"type:text;not null" json:"description"`
	ShortDescription string `gorm:"type:text;not null" json:"short_description"`
	ImagePath        string `gorm:"type:text;not null" json:"image_path,omitempty"`
	VideoPath        string `gorm:"type:text;not null" json:"video_path,omitempty"`
	Players          string `gorm:"type:varchar(50);not null" json:"players"`
	IsDeleted        bool   `gorm:"type:boolean;not null;default:false" json:"-"`
}
