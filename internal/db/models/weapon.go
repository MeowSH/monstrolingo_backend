package models

import (
	"time"

	"github.com/google/uuid"
)

type Weapon struct {
	BaseModel
	SoftDeleteModel

	ExternalKey string     `gorm:"type:text;not null;uniqueIndex:ux_weapons_external_key"`
	SyncedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index"`
	SyncRunID   *uuid.UUID `gorm:"type:uuid;index"`

	WeaponType      string `gorm:"type:text;not null;index"`
	Rarity          int16  `gorm:"not null;default:0;index"`
	Attack          int32  `gorm:"not null;default:0"`
	AffinityPercent int16  `gorm:"not null;default:0"`
	DefenseBonus    int16  `gorm:"not null;default:0"`
	Slot1Level      int16  `gorm:"not null;default:0"`
	Slot2Level      int16  `gorm:"not null;default:0"`
	Slot3Level      int16  `gorm:"not null;default:0"`
	ElementType     string `gorm:"type:text;index"`
	ElementValue    int32  `gorm:"not null;default:0"`
	AilmentType     string `gorm:"type:text;index"`
	AilmentValue    int32  `gorm:"not null;default:0"`
	SharpnessRed    int16  `gorm:"not null;default:0"`
	SharpnessOrange int16  `gorm:"not null;default:0"`
	SharpnessYellow int16  `gorm:"not null;default:0"`
	SharpnessGreen  int16  `gorm:"not null;default:0"`
	SharpnessBlue   int16  `gorm:"not null;default:0"`
	SharpnessWhite  int16  `gorm:"not null;default:0"`
	SharpnessPurple int16  `gorm:"not null;default:0"`
	CraftCost       *int32
	UpgradeCost     *int32
	TreeDepth       int16 `gorm:"not null;default:0"`
	IsFinalUpgrade  bool  `gorm:"not null;default:false"`
	SortOrder       int32 `gorm:"not null;default:0;index"`

	Translations []WeaponTranslation `gorm:"foreignKey:WeaponID"`
	Skills       []WeaponSkill       `gorm:"foreignKey:WeaponID"`
}

type WeaponTranslation struct {
	BaseModel

	WeaponID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_weapon_translations_weapon_language"`
	LanguageID  uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_weapon_translations_weapon_language"`
	Name        string    `gorm:"type:text;not null;index"`
	Description string    `gorm:"type:text"`
	Slug        string    `gorm:"type:text;index"`

	Weapon   Weapon   `gorm:"foreignKey:WeaponID;constraint:OnDelete:CASCADE"`
	Language Language `gorm:"foreignKey:LanguageID"`
}

type WeaponSkill struct {
	BaseModel

	WeaponID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_weapon_skills_weapon_skill_order"`
	SkillID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_weapon_skills_weapon_skill_order"`
	SortOrder  int16     `gorm:"not null;default:0;uniqueIndex:ux_weapon_skills_weapon_skill_order"`
	Level      int16     `gorm:"not null;default:1"`
	SourceType string    `gorm:"type:text;not null;default:'base'"`

	Weapon Weapon `gorm:"foreignKey:WeaponID;constraint:OnDelete:CASCADE"`
	Skill  Skill  `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}
