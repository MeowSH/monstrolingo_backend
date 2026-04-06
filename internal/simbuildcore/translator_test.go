package simbuildcore

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildSkillResponsePayload_PartialWhenUnresolvedSkill(t *testing.T) {
	skillID := uuid.MustParse("00000000-0000-0000-0000-00000000000a")
	resolved := []resolvedSkill{
		{
			Requested: requestedSkill{
				OriginalText:   "Attack Boost Lv3",
				BaseName:       "Attack Boost",
				RequestedLevel: 3,
			},
			SkillID:     skillID,
			ExternalKey: "attack-boost",
			MaxLevel:    5,
			Resolved:    true,
			Translated:  true,
			TargetName:  "Boost d'Attaque",
		},
		{
			Requested: requestedSkill{
				OriginalText:   "Unknown Skill Lv1",
				BaseName:       "Unknown Skill",
				RequestedLevel: 1,
			},
			Resolved:   false,
			Translated: false,
			TargetName: "Unknown Skill Lv1",
		},
	}

	original, translated, unmatched := buildSkillResponsePayload(resolved)
	if len(original) != 2 || len(translated) != 2 {
		t.Fatalf("expected 2 output skills, got original=%d translated=%d", len(original), len(translated))
	}
	if len(unmatched) != 1 {
		t.Fatalf("expected 1 unmatched entry, got %d", len(unmatched))
	}
	if unmatched[0].Kind != "skill" {
		t.Fatalf("expected unmatched kind skill, got %s", unmatched[0].Kind)
	}
	if translated[0].SkillExternalKey != "attack-boost" {
		t.Fatalf("expected translated external key attack-boost, got %s", translated[0].SkillExternalKey)
	}
}

func TestResolveSkills_UsesAliasExternalKeyWhenNameNotFound(t *testing.T) {
	skillID := uuid.MustParse("00000000-0000-0000-0000-0000000000bb")
	requested := []requestedSkill{
		{
			OriginalText:   "Guts (Tenacity)",
			BaseName:       "Guts (Tenacity)",
			RequestedLevel: 1,
		},
	}

	byLanguage := map[string]map[string][]skillTranslationRow{
		"en": {},
	}
	byEffect := map[string]map[string][]skillTranslationRow{
		"en": {},
	}
	byExternal := map[string]map[string]skillTranslationRow{
		"lords-soul": {
			"en": {
				SkillID:      skillID,
				ExternalKey:  "lords-soul",
				MaxLevel:     3,
				LanguageCode: "en",
				Name:         "Lord's Soul",
			},
		},
	}
	bySkillID := map[uuid.UUID]map[string]skillTranslationRow{
		skillID: {
			"en": {
				SkillID:      skillID,
				ExternalKey:  "lords-soul",
				MaxLevel:     3,
				LanguageCode: "en",
				Name:         "Lord's Soul",
			},
		},
	}

	resolved := resolveSkills(requested, "en", "en", byLanguage, byEffect, byExternal, bySkillID)
	if len(resolved) != 1 {
		t.Fatalf("expected one resolved entry, got %d", len(resolved))
	}
	if !resolved[0].Resolved {
		t.Fatalf("expected alias resolution to resolve the skill")
	}
	if resolved[0].ExternalKey != "lords-soul" {
		t.Fatalf("expected external key lords-soul, got %s", resolved[0].ExternalKey)
	}
	if resolved[0].TargetName != "Lord's Soul" {
		t.Fatalf("expected target name Lord's Soul, got %s", resolved[0].TargetName)
	}
}

func TestResolveSkills_UsesScorcherAlias(t *testing.T) {
	skillID := uuid.MustParse("00000000-0000-0000-0000-0000000000bc")
	requested := []requestedSkill{
		{
			OriginalText:   "Scorcher I",
			BaseName:       "Scorcher",
			RequestedLevel: 1,
		},
	}

	byLanguage := map[string]map[string][]skillTranslationRow{
		"en": {},
	}
	byEffect := map[string]map[string][]skillTranslationRow{
		"en": {},
	}
	byExternal := map[string]map[string]skillTranslationRow{
		"rathaloss-flare": {
			"en": {
				SkillID:      skillID,
				ExternalKey:  "rathaloss-flare",
				MaxLevel:     1,
				LanguageCode: "en",
				Name:         "Rathalos's Flare",
			},
		},
	}
	bySkillID := map[uuid.UUID]map[string]skillTranslationRow{
		skillID: {
			"en": {
				SkillID:      skillID,
				ExternalKey:  "rathaloss-flare",
				MaxLevel:     1,
				LanguageCode: "en",
				Name:         "Rathalos's Flare",
			},
		},
	}

	resolved := resolveSkills(requested, "en", "en", byLanguage, byEffect, byExternal, bySkillID)
	if len(resolved) != 1 || !resolved[0].Resolved {
		t.Fatalf("expected Scorcher alias to resolve, got %+v", resolved)
	}
	if resolved[0].ExternalKey != "rathaloss-flare" {
		t.Fatalf("expected external key rathaloss-flare, got %s", resolved[0].ExternalKey)
	}
}

func TestResolveSkills_UsesJapaneseAliasCanonicalFallback(t *testing.T) {
	withAliasSnapshotForTest(t, []aliasSnapshotEntry{
		{
			LanguageCode:        "ja",
			Category:            "set",
			Alias:               "黒蝕一体Ⅰ",
			AliasNormalized:     normalizeSkillName("黒蝕一体Ⅰ"),
			AliasLookupKeys:     buildAliasLookupKeys("黒蝕一体Ⅰ"),
			CanonicalName:       "黒蝕竜の力",
			CanonicalNormalized: normalizeSkillName("黒蝕竜の力"),
			ActivationLevel:     2,
		},
	})

	skillID := uuid.MustParse("00000000-0000-0000-0000-0000000000bd")
	requested := []requestedSkill{
		{
			OriginalText:   "黒蝕一体Ⅰ",
			BaseName:       "黒蝕一体Ⅰ",
			RequestedLevel: 1,
		},
	}

	jaCanonical := skillTranslationRow{
		SkillID:      skillID,
		ExternalKey:  "gore-magalas-tyranny",
		MaxLevel:     3,
		LanguageCode: "ja",
		Name:         "黒蝕竜の力",
	}
	enCanonical := skillTranslationRow{
		SkillID:      skillID,
		ExternalKey:  "gore-magalas-tyranny",
		MaxLevel:     3,
		LanguageCode: "en",
		Name:         "Gore Magala's Tyranny",
	}

	byLanguage := map[string]map[string][]skillTranslationRow{
		"ja": {
			normalizeSkillName("黒蝕竜の力"): {jaCanonical},
		},
		"en": {
			normalizeSkillName("Gore Magala's Tyranny"): {enCanonical},
		},
	}
	byEffect := map[string]map[string][]skillTranslationRow{
		"ja": {},
		"en": {},
	}
	byExternal := map[string]map[string]skillTranslationRow{
		"gore-magalas-tyranny": {
			"ja": jaCanonical,
			"en": enCanonical,
		},
	}
	bySkillID := map[uuid.UUID]map[string]skillTranslationRow{
		skillID: {
			"ja": jaCanonical,
			"en": enCanonical,
		},
	}

	resolved := resolveSkills(requested, "ja", "en", byLanguage, byEffect, byExternal, bySkillID)
	if len(resolved) != 1 || !resolved[0].Resolved {
		t.Fatalf("expected Japanese alias to resolve, got %+v", resolved)
	}
	if resolved[0].ExternalKey != "gore-magalas-tyranny" {
		t.Fatalf("expected external key gore-magalas-tyranny, got %s", resolved[0].ExternalKey)
	}
	if resolved[0].TargetName != "Gore Magala's Tyranny" {
		t.Fatalf("expected translated target name in en fallback, got %s", resolved[0].TargetName)
	}
}

func TestSplitTranslatedSkillsByKind(t *testing.T) {
	resolved := []resolvedSkill{
		{
			Requested:  requestedSkill{OriginalText: "A", BaseName: "A", RequestedLevel: 1},
			Resolved:   true,
			IsSetBonus: true,
		},
		{
			Requested:  requestedSkill{OriginalText: "B", BaseName: "B", RequestedLevel: 1},
			Resolved:   true,
			IsSetBonus: false,
		},
		{
			Requested: requestedSkill{OriginalText: "C", BaseName: "C", RequestedLevel: 1},
			Resolved:  false,
		},
	}
	translated := []TranslatedSkill{
		{OriginalText: "A", Name: "A-fr"},
		{OriginalText: "B", Name: "B-fr"},
		{OriginalText: "C", Name: "C-fr"},
	}

	setSkills, armorJewelSkills := splitTranslatedSkillsByKind(resolved, translated)
	if len(setSkills) != 1 {
		t.Fatalf("expected 1 set skill, got %d", len(setSkills))
	}
	if len(armorJewelSkills) != 1 {
		t.Fatalf("expected 1 armor/jewel skill, got %d", len(armorJewelSkills))
	}
	if setSkills[0].OriginalText != "A" {
		t.Fatalf("unexpected set skill bucket content: %+v", setSkills[0])
	}
	if armorJewelSkills[0].OriginalText != "B" {
		t.Fatalf("unexpected armor/jewel bucket content: %+v", armorJewelSkills[0])
	}
}

func TestSplitTranslatedSkillsByKind_UsesAliasFallbackWhenFlagMissing(t *testing.T) {
	resolved := []resolvedSkill{
		{
			Requested: requestedSkill{
				OriginalText:   "Scorcher I",
				BaseName:       "Scorcher",
				RequestedLevel: 1,
			},
			Resolved:   true,
			IsSetBonus: false,
		},
	}
	translated := []TranslatedSkill{
		{OriginalText: "Scorcher I", Name: "Flamboiement I"},
	}

	setSkills, armorJewelSkills := splitTranslatedSkillsByKind(resolved, translated)
	if len(setSkills) != 1 {
		t.Fatalf("expected alias fallback to classify as set skill, got %d", len(setSkills))
	}
	if len(armorJewelSkills) != 0 {
		t.Fatalf("expected no armor/jewel skills when alias fallback matches set/group, got %d", len(armorJewelSkills))
	}
}

func TestAttachAssociatedJewels_TargetLanguageThenEnglishFallback(t *testing.T) {
	skillIDOne := uuid.MustParse("00000000-0000-0000-0000-0000000000d1")
	skillIDTwo := uuid.MustParse("00000000-0000-0000-0000-0000000000d2")

	resolved := []resolvedSkill{
		{
			Requested: requestedSkill{
				OriginalText:   "Attack Boost Lv3",
				BaseName:       "Attack Boost",
				RequestedLevel: 3,
			},
			Resolved: true,
			SkillID:  skillIDOne,
		},
		{
			Requested: requestedSkill{
				OriginalText:   "Unknown Skill Lv1",
				BaseName:       "Unknown Skill",
				RequestedLevel: 1,
			},
			Resolved: false,
			SkillID:  skillIDTwo,
		},
	}
	translated := []TranslatedSkill{
		{OriginalText: "Attack Boost Lv3", Name: "Boost d'Attaque"},
		{OriginalText: "Unknown Skill Lv1", Name: "Unknown Skill Lv1"},
	}

	decorationRows := []decorationTranslationRow{
		{
			SkillID:               skillIDOne,
			DecorationExternalKey: "attack-jewel-1",
			SlotSize:              1,
			Rarity:                5,
			SkillLevel:            1,
			LanguageCode:          "en",
			Name:                  "Attack Jewel [1]",
		},
		{
			SkillID:               skillIDOne,
			DecorationExternalKey: "attack-jewel-1",
			SlotSize:              1,
			Rarity:                5,
			SkillLevel:            1,
			LanguageCode:          "fr",
			Name:                  "Joyau attaque [1]",
		},
		{
			SkillID:               skillIDOne,
			DecorationExternalKey: "critical-jewel-2",
			SlotSize:              2,
			Rarity:                7,
			SkillLevel:            1,
			LanguageCode:          "en",
			Name:                  "Critical Jewel [2]",
		},
	}

	enriched := attachAssociatedJewels(resolved, translated, decorationRows, "fr")
	if len(enriched) != 2 {
		t.Fatalf("expected 2 translated rows, got %d", len(enriched))
	}

	if len(enriched[0].AssociatedJewels) != 2 {
		t.Fatalf("expected 2 associated jewels on first skill, got %d", len(enriched[0].AssociatedJewels))
	}
	first := enriched[0].AssociatedJewels[0]
	second := enriched[0].AssociatedJewels[1]
	if first.DecorationExternalKey != "attack-jewel-1" {
		t.Fatalf("expected first jewel to be attack-jewel-1, got %q", first.DecorationExternalKey)
	}
	if first.Name != "Joyau attaque [1]" {
		t.Fatalf("expected target-language jewel name, got %q", first.Name)
	}
	if second.DecorationExternalKey != "critical-jewel-2" {
		t.Fatalf("expected second jewel to be critical-jewel-2, got %q", second.DecorationExternalKey)
	}
	if second.Name != "Critical Jewel [2]" {
		t.Fatalf("expected english fallback jewel name, got %q", second.Name)
	}

	if len(enriched[1].AssociatedJewels) != 0 {
		t.Fatalf("expected no associated jewel for unresolved skill, got %+v", enriched[1].AssociatedJewels)
	}
}

func TestDeduplicateResolvedSkills_MergesAliasAndCanonicalEntries(t *testing.T) {
	skillID := uuid.MustParse("00000000-0000-0000-0000-0000000000d3")
	rows := []resolvedSkill{
		{
			Requested: requestedSkill{
				OriginalText:   "Guts (Tenacity)",
				BaseName:       "Guts (Tenacity)",
				RequestedLevel: 1,
			},
			SkillID:     skillID,
			ExternalKey: "lords-soul",
			MaxLevel:    3,
			IsSetBonus:  true,
			TargetName:  "Âme du seigneur",
			Resolved:    true,
			Translated:  true,
		},
		{
			Requested: requestedSkill{
				OriginalText:   "Lord's Soul",
				BaseName:       "Lord's Soul",
				RequestedLevel: 3,
			},
			SkillID:     skillID,
			ExternalKey: "lords-soul",
			MaxLevel:    3,
			IsSetBonus:  true,
			TargetName:  "Âme du seigneur",
			Resolved:    true,
			Translated:  true,
		},
	}

	out := deduplicateResolvedSkills(rows)
	if len(out) != 1 {
		t.Fatalf("expected 1 deduplicated row, got %d", len(out))
	}
	if out[0].ExternalKey != "lords-soul" {
		t.Fatalf("expected external key lords-soul, got %q", out[0].ExternalKey)
	}
	if out[0].Requested.RequestedLevel != 3 {
		t.Fatalf("expected highest requested level to be kept, got %d", out[0].Requested.RequestedLevel)
	}
	if out[0].Requested.OriginalText != "Guts (Tenacity)" {
		t.Fatalf("expected first original text to be preserved, got %q", out[0].Requested.OriginalText)
	}
}
