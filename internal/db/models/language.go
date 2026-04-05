package models

type Language struct {
	BaseModel

	Code      string `gorm:"type:varchar(16);not null;uniqueIndex"`
	Label     string `gorm:"type:varchar(64);not null"`
	IsActive  bool   `gorm:"not null;default:true;index"`
	SortOrder int16  `gorm:"not null;default:0"`
}
