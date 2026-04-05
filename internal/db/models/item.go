package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_items_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	Category           string `gorm:"type:text;not null;index"`
	Subcategory        string `gorm:"type:text;index"`
	Rarity             int16  `gorm:"not null;default:0;index"`
	CarryLimit         *int16
	BuyPrice           *int32
	SellPrice          *int32
	Points             *int32
	IsCraftingMaterial bool   `gorm:"not null;default:false"`
	IsConsumable       bool   `gorm:"not null;default:false"`
	IconKey            string `gorm:"type:text"`
	SortOrder          int32  `gorm:"not null;default:0;index"`

	Translations []ItemTranslation `gorm:"foreignKey:ItemID"`
}

type ItemTranslation struct {
	BaseModel

	ItemID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_item_translations_item_language"`
	LanguageID  uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_item_translations_item_language"`
	Name        string    `gorm:"type:text;not null;index"`
	Description string    `gorm:"type:text"`
	FlavorText  string    `gorm:"type:text"`
	Slug        string    `gorm:"type:text;index"`

	Item     Item     `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Language Language `gorm:"foreignKey:LanguageID"`
}
