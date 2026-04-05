package models

import (
	"time"

	"github.com/google/uuid"
)

type Decoration struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_decorations_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	SlotSize    int16 `gorm:"not null;default:1;index"`
	Rarity      int16 `gorm:"not null;default:0;index"`
	IsCraftable bool  `gorm:"not null;default:true"`
	CraftCost   *int32
	SortOrder   int32 `gorm:"not null;default:0;index"`

	Translations []DecorationTranslation `gorm:"foreignKey:DecorationID"`
	Skills       []DecorationSkill       `gorm:"foreignKey:DecorationID"`
}

type DecorationTranslation struct {
	BaseModel

	DecorationID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_decorations_translations_decoration_language"`
	LanguageID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_decorations_translations_decoration_language"`
	Name         string    `gorm:"type:text;not null;index"`
	Description  string    `gorm:"type:text"`
	Slug         string    `gorm:"type:text;index"`

	Decoration Decoration `gorm:"foreignKey:DecorationID;constraint:OnDelete:CASCADE"`
	Language   Language   `gorm:"foreignKey:LanguageID"`
}

func (DecorationTranslation) TableName() string {
	return "decorations_translations"
}

type DecorationSkill struct {
	BaseModel

	DecorationID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_decoration_skills_decoration_skill_order"`
	SkillID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_decoration_skills_decoration_skill_order"`
	SortOrder    int16     `gorm:"not null;default:0;uniqueIndex:ux_decoration_skills_decoration_skill_order"`
	Level        int16     `gorm:"not null;default:1"`

	Decoration Decoration `gorm:"foreignKey:DecorationID;constraint:OnDelete:CASCADE"`
	Skill      Skill      `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}
