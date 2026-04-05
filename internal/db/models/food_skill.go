package models

import (
	"time"

	"github.com/google/uuid"
)

type FoodSkill struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_food_skills_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	FoodCategory          string `gorm:"type:text;not null;index"`
	MaxLevel              int16  `gorm:"not null;default:1"`
	BaseDurationSeconds   *int32
	BaseActivationPercent *int16
	SortOrder             int32 `gorm:"not null;default:0;index"`

	Translations []FoodSkillTranslation `gorm:"foreignKey:FoodSkillID"`
	Levels       []FoodSkillLevel       `gorm:"foreignKey:FoodSkillID"`
}

type FoodSkillTranslation struct {
	BaseModel

	FoodSkillID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_food_skill_translations_food_skill_language"`
	LanguageID  uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_food_skill_translations_food_skill_language"`
	Name        string    `gorm:"type:text;not null;index"`
	Description string    `gorm:"type:text"`
	Slug        string    `gorm:"type:text;index"`

	FoodSkill FoodSkill `gorm:"foreignKey:FoodSkillID;constraint:OnDelete:CASCADE"`
	Language  Language  `gorm:"foreignKey:LanguageID"`
}

type FoodSkillLevel struct {
	BaseModel

	FoodSkillID       uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_food_skill_levels_food_skill_level"`
	Level             int16     `gorm:"not null;default:1;uniqueIndex:ux_food_skill_levels_food_skill_level"`
	DurationSeconds   *int32
	ActivationPercent *int16
	EffectValueText   string `gorm:"type:text"`

	FoodSkill FoodSkill `gorm:"foreignKey:FoodSkillID;constraint:OnDelete:CASCADE"`
}
