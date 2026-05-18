package ds

type PlayersChosenLocation struct {
	ID         uint `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID  uint `gorm:"not null;uniqueIndex:idx_request_location" json:"request_id"`
	LocationID uint `gorm:"not null;uniqueIndex:idx_request_location" json:"location_id"`
	Priority   int  `json:"priority"`
	IsRandomed bool `json:"is_randomed"`

	Request  PlayersLocationGame `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE" json:"request,omitempty"`
	Location Location            `gorm:"foreignKey:LocationID;constraint:OnDelete:CASCADE" json:"location,omitempty"`
}
