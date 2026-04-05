package models

import (
	"time"

	"github.com/google/uuid"
)

type Kinsect struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_kinsects_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	KinsectType           string `gorm:"type:text;not null;index"`
	AttackType            string `gorm:"type:text;index"`
	PowderType            string `gorm:"type:text;index"`
	KinsectBonusPrimary   string `gorm:"type:text;index"`
	KinsectBonusSecondary string `gorm:"type:text;index"`
	Rarity                int16  `gorm:"not null;default:0;index"`
	PowerValue            int16  `gorm:"not null;default:0"`
	SpeedValue            int16  `gorm:"not null;default:0"`
	HealValue             int16  `gorm:"not null;default:0"`
	StaminaValue          int16  `gorm:"not null;default:0"`
	ElementType           string `gorm:"type:text;index"`
	ElementValue          int32  `gorm:"not null;default:0"`
	SortOrder             int32  `gorm:"not null;default:0;index"`

	Translations []KinsectTranslation `gorm:"foreignKey:KinsectID"`
}

type KinsectTranslation struct {
	BaseModel

	KinsectID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_kinsect_translations_kinsect_language"`
	LanguageID  uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_kinsect_translations_kinsect_language"`
	Name        string    `gorm:"type:text;not null;index"`
	Description string    `gorm:"type:text"`
	Slug        string    `gorm:"type:text;index"`

	Kinsect  Kinsect  `gorm:"foreignKey:KinsectID;constraint:OnDelete:CASCADE"`
	Language Language `gorm:"foreignKey:LanguageID"`
}
