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
