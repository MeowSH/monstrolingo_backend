package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel is shared by all persisted entities.
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// SoftDeleteModel is embedded only on canonical tables.
type SoftDeleteModel struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
