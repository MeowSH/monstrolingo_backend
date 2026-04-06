package simbuildcore

import "testing"

func TestDetectSourceLanguage_EnglishBestMatch(t *testing.T) {
	byLanguage := map[string]map[string][]skillTranslationRow{
		"en": {
			normalizeSkillName("Maximum Might"): {{Name: "Maximum Might", LanguageCode: "en"}},
			normalizeSkillName("Agitator"):      {{Name: "Agitator", LanguageCode: "en"}},
		},
		"ja": {
			normalizeSkillName("渾身"):  {{Name: "渾身", LanguageCode: "ja"}},
			normalizeSkillName("挑戦者"): {{Name: "挑戦者", LanguageCode: "ja"}},
		},
	}
	skills := []requestedSkill{
		{BaseName: "Maximum Might", RequestedLevel: 3},
		{BaseName: "Agitator", RequestedLevel: 2},
	}

	lang := detectSourceLanguage(skills, byLanguage, []string{"en", "ja"})
	if lang != "en" {
		t.Fatalf("expected en, got %s", lang)
	}
}

func TestDetectSourceLanguage_FallbackToEnglish(t *testing.T) {
	byLanguage := map[string]map[string][]skillTranslationRow{
		"en": {normalizeSkillName("Attack Boost"): {{Name: "Attack Boost", LanguageCode: "en"}}},
		"ja": {normalizeSkillName("攻撃"): {{Name: "攻撃", LanguageCode: "ja"}}},
	}
	skills := []requestedSkill{
		{BaseName: "Unknown Skill", RequestedLevel: 1},
	}

	lang := detectSourceLanguage(skills, byLanguage, []string{"ja", "en"})
	if lang != "en" {
		t.Fatalf("expected en fallback, got %s", lang)
	}
}

func TestDetectSourceLanguage_UsesAliasHitWhenNamesMissing(t *testing.T) {
	withAliasSnapshotForTest(t, []aliasSnapshotEntry{
		{
			LanguageCode:    "ja",
			Category:        "set",
			Alias:           "黒蝕一体Ⅰ",
			AliasNormalized: normalizeSkillName("黒蝕一体Ⅰ"),
			AliasLookupKeys: buildAliasLookupKeys("黒蝕一体Ⅰ"),
			CanonicalName:   "黒蝕竜の力",
			ActivationLevel: 2,
		},
	})

	byLanguage := map[string]map[string][]skillTranslationRow{
		"en": {},
		"ja": {},
	}
	skills := []requestedSkill{
		{BaseName: "黒蝕一体Ⅰ", RequestedLevel: 1},
	}

	lang := detectSourceLanguage(skills, byLanguage, []string{"en", "ja"})
	if lang != "ja" {
		t.Fatalf("expected ja based on alias hit, got %s", lang)
	}
}
