package models

import (
	"time"

	"github.com/google/uuid"
)

type Charm struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_charms_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	Rarity    int16 `gorm:"not null;default:0;index"`
	MaxRank   int16 `gorm:"not null;default:1"`
	CraftCost *int32
	SortOrder int32 `gorm:"not null;default:0;index"`

	Translations []CharmTranslation `gorm:"foreignKey:CharmID"`
	Skills       []CharmSkill       `gorm:"foreignKey:CharmID"`
}

type CharmTranslation struct {
	BaseModel

	CharmID     uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_charm_translations_charm_language"`
	LanguageID  uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_charm_translations_charm_language"`
	Name        string    `gorm:"type:text;not null;index"`
	Description string    `gorm:"type:text"`
	Slug        string    `gorm:"type:text;index"`

	Charm    Charm    `gorm:"foreignKey:CharmID;constraint:OnDelete:CASCADE"`
	Language Language `gorm:"foreignKey:LanguageID"`
}

type CharmSkill struct {
	BaseModel

	CharmID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_charm_skills_charm_skill_rank"`
	SkillID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_charm_skills_charm_skill_rank"`
	RankRequired int16     `gorm:"not null;default:1;uniqueIndex:ux_charm_skills_charm_skill_rank"`
	Level        int16     `gorm:"not null;default:1"`
	SortOrder    int16     `gorm:"not null;default:0"`

	Charm Charm `gorm:"foreignKey:CharmID;constraint:OnDelete:CASCADE"`
	Skill Skill `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}
