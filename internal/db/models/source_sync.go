package models

import (
	"time"
)

type SourceSyncRun struct {
	BaseModel

	Category       string     `gorm:"type:text;not null;index"`
	Status         string     `gorm:"type:text;not null;index"`
	TriggeredBy    string     `gorm:"type:text;not null;default:'manual'"`
	DryRun         bool       `gorm:"not null;default:false"`
	CacheMode      string     `gorm:"type:text;not null;default:'live'"`
	LanguagesCount int16      `gorm:"not null;default:0"`
	StartedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	FinishedAt     *time.Time `gorm:"type:timestamptz"`
	ErrorMessage   string     `gorm:"type:text"`
}
