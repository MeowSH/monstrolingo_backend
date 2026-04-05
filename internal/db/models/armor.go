package models

import (
	"time"

	"github.com/google/uuid"
)

type ArmorPiece struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_armor_pieces_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	ArmorSetKey         string `gorm:"type:text;not null;index"`
	ArmorSetName        string `gorm:"type:text;index"`
	ArmorSetVariant     string `gorm:"type:text;index"`
	PieceType           string `gorm:"type:text;not null;index"`
	Rank                string `gorm:"type:text;not null;index"`
	Gender              string `gorm:"type:text;index"`
	Rarity              int16  `gorm:"not null;default:0;index"`
	DefenseBase         int16  `gorm:"not null;default:0"`
	DefenseMax          int16  `gorm:"not null;default:0"`
	DefenseAugmentedMax int16  `gorm:"not null;default:0"`
	FireRes             int16  `gorm:"not null;default:0"`
	WaterRes            int16  `gorm:"not null;default:0"`
	ThunderRes          int16  `gorm:"not null;default:0"`
	IceRes              int16  `gorm:"not null;default:0"`
	DragonRes           int16  `gorm:"not null;default:0"`
	Slot1Level          int16  `gorm:"not null;default:0"`
	Slot2Level          int16  `gorm:"not null;default:0"`
	Slot3Level          int16  `gorm:"not null;default:0"`
	IsLayered           bool   `gorm:"not null;default:false"`
	SortOrder           int32  `gorm:"not null;default:0;index"`

	Translations []ArmorPieceTranslation `gorm:"foreignKey:ArmorPieceID"`
	Skills       []ArmorPieceSkill       `gorm:"foreignKey:ArmorPieceID"`
}

type ArmorPieceTranslation struct {
	BaseModel

	ArmorPieceID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_armor_piece_translations_piece_language"`
	LanguageID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_armor_piece_translations_piece_language"`
	Name         string    `gorm:"type:text;not null;index"`
	Description  string    `gorm:"type:text"`
	Slug         string    `gorm:"type:text;index"`

	ArmorPiece ArmorPiece `gorm:"foreignKey:ArmorPieceID;constraint:OnDelete:CASCADE"`
	Language   Language   `gorm:"foreignKey:LanguageID"`
}

type ArmorPieceSkill struct {
	BaseModel

	ArmorPieceID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_armor_piece_skills_piece_skill_order"`
	SkillID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_armor_piece_skills_piece_skill_order"`
	SortOrder    int16     `gorm:"not null;default:0;uniqueIndex:ux_armor_piece_skills_piece_skill_order"`
	Level        int16     `gorm:"not null;default:1"`

	ArmorPiece ArmorPiece `gorm:"foreignKey:ArmorPieceID;constraint:OnDelete:CASCADE"`
	Skill      Skill      `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}
