package models

// All returns all GORM entities in dependency-safe order.
func All() []any {
	return []any{
		&Language{},
		&SourceSyncRun{},

		&Skill{},
		&SkillTranslation{},
		&SkillLevel{},

		&Item{},
		&ItemTranslation{},

		&Weapon{},
		&WeaponTranslation{},
		&WeaponSkill{},

		&ArmorPiece{},
		&ArmorPieceTranslation{},
		&ArmorPieceSkill{},

		&Decoration{},
		&DecorationTranslation{},
		&DecorationSkill{},

		&Charm{},
		&CharmTranslation{},
		&CharmSkill{},

		&FoodSkill{},
		&FoodSkillTranslation{},
		&FoodSkillLevel{},

		&Kinsect{},
		&KinsectTranslation{},
	}
}
