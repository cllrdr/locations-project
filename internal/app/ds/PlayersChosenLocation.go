package ds

type PlayersChosenLocation struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	RequestID  uint `gorm:"not null;uniqueIndex:idx_request_location"`
	LocationID uint `gorm:"not null;uniqueIndex:idx_request_location"`
	Priority   int  `gorm:"not null"`

	Request  PlayersLocationRequest `gorm:"foreignKey:RequestID"`
	Location Location               `gorm:"foreignKey:LocationID"`
}
