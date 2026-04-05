package models

import (
	"time"

	"github.com/google/uuid"
)

type Skill struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_skills_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	SkillKind       string `gorm:"type:text;not null;index"`
	MaxLevel        int16  `gorm:"not null;default:1"`
	IsBinary        bool   `gorm:"not null;default:false"`
	IsSetBonusSkill bool   `gorm:"not null;default:false"`
	SortOrder       int32  `gorm:"not null;default:0;index"`

	Translations []SkillTranslation `gorm:"foreignKey:SkillID"`
	Levels       []SkillLevel       `gorm:"foreignKey:SkillID"`
}

type SkillTranslation struct {
	BaseModel

	SkillID       uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_skill_translations_skill_language"`
	LanguageID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_skill_translations_skill_language"`
	Name          string    `gorm:"type:text;not null;index"`
	Description   string    `gorm:"type:text"`
	EffectSummary string    `gorm:"type:text"`
	Slug          string    `gorm:"type:text;index"`

	Skill    Skill    `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
	Language Language `gorm:"foreignKey:LanguageID"`
}

type SkillLevel struct {
	BaseModel

	SkillID         uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_skill_levels_skill_level"`
	Level           int16     `gorm:"not null;default:1;uniqueIndex:ux_skill_levels_skill_level"`
	EffectValueText string    `gorm:"type:text"`
	Description     string    `gorm:"type:text"`

	Skill Skill `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}
