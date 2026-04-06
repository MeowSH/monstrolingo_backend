package simbuildcore

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

const (
	aliasSnapshotPathEnv = "SIMBUILD_ALIAS_SNAPSHOT_PATH"
	// Local/private snapshot path (gitignored by default).
	defaultAliasSnapshotPath = "internal/simbuildcore/localbridge/data/local_alias_snapshot.json"
)

type skillAlias struct {
	ExternalKey         string
	CanonicalName       string
	CanonicalNormalized string
	Category            string
	ActivationLevel     int16
}

type skillAliasIndex struct {
	byLanguage map[string]map[string]skillAlias
	byAny      map[string]skillAlias
}

type aliasSnapshot struct {
	Entries []aliasSnapshotEntry `json:"entries"`
}

type aliasSnapshotEntry struct {
	LanguageCode        string   `json:"language_code"`
	Category            string   `json:"category"`
	Alias               string   `json:"alias"`
	AliasNormalized     string   `json:"alias_normalized"`
	AliasLookupKeys     []string `json:"alias_lookup_keys"`
	CanonicalName       string   `json:"canonical_name"`
	CanonicalNormalized string   `json:"canonical_normalized"`
	ActivationLevel     int16    `json:"activation_level"`
}

var (
	skillAliasIndexOnce sync.Once
	skillAliasIndexInst skillAliasIndex
)

var legacySkillAliasByName = map[string]skillAlias{
	"guts (tenacity)": {
		ExternalKey:         "lords-soul",
		CanonicalName:       "Lord's Soul",
		CanonicalNormalized: "lord's soul",
		Category:            "group",
		ActivationLevel:     3,
	},
	"binding counter": {
		ExternalKey:         "jin-dahaads-revolt",
		CanonicalName:       "Jin Dahaad's Revolt",
		CanonicalNormalized: "jin dahaad's revolt",
		Category:            "set",
		ActivationLevel:     2,
	},
	"black eclipse": {
		ExternalKey:         "gore-magalas-tyranny",
		CanonicalName:       "Gore Magala's Tyranny",
		CanonicalNormalized: "gore magala's tyranny",
		Category:            "set",
		ActivationLevel:     2,
	},
	"scorcher": {
		ExternalKey:         "rathaloss-flare",
		CanonicalName:       "Rathalos's Flare",
		CanonicalNormalized: "rathalos's flare",
		Category:            "set",
		ActivationLevel:     2,
	},
}

func resolveSkillAlias(normalizedName string, sourceLang string) (skillAlias, bool) {
	key := normalizeSkillName(normalizedName)
	if key == "" {
		return skillAlias{}, false
	}
	index := loadSkillAliasIndex()
	lang := normalizeLanguageCode(sourceLang)
	if lang != "" {
		if alias, ok := index.byLanguage[lang][key]; ok {
			return alias, true
		}
	}
	if lang != "en" {
		if alias, ok := index.byLanguage["en"][key]; ok {
			return alias, true
		}
	}
	alias, ok := index.byAny[key]
	if !ok {
		return skillAlias{}, false
	}
	return alias, true
}

func hasSkillAliasForLanguage(lang string, normalizedName string) bool {
	key := normalizeSkillName(normalizedName)
	if key == "" {
		return false
	}
	index := loadSkillAliasIndex()
	language := normalizeLanguageCode(lang)
	if language == "" {
		return false
	}
	_, ok := index.byLanguage[language][key]
	return ok
}

func aliasCanonicalLookupKey(alias skillAlias) string {
	key := normalizeSkillName(alias.CanonicalNormalized)
	if key != "" {
		return key
	}
	return normalizeSkillName(alias.CanonicalName)
}

func loadSkillAliasIndex() skillAliasIndex {
	skillAliasIndexOnce.Do(func() {
		skillAliasIndexInst = buildSkillAliasIndex()
	})
	return skillAliasIndexInst
}

func buildSkillAliasIndex() skillAliasIndex {
	index := skillAliasIndex{
		byLanguage: make(map[string]map[string]skillAlias, 8),
		byAny:      make(map[string]skillAlias, 512),
	}

	for key, alias := range legacySkillAliasByName {
		putAlias(index.byLanguage, "en", key, alias)
		putAliasAny(index.byAny, key, alias)
	}

	for _, entry := range readAliasSnapshotEntries() {
		lang := normalizeLanguageCode(entry.LanguageCode)
		if lang == "" {
			continue
		}
		alias := skillAlias{
			CanonicalName:       strings.TrimSpace(entry.CanonicalName),
			CanonicalNormalized: normalizeSkillName(entry.CanonicalNormalized),
			Category:            strings.TrimSpace(strings.ToLower(entry.Category)),
			ActivationLevel:     entry.ActivationLevel,
		}
		if alias.CanonicalNormalized == "" {
			alias.CanonicalNormalized = normalizeSkillName(entry.CanonicalName)
		}

		keys := entry.AliasLookupKeys
		if len(keys) == 0 && strings.TrimSpace(entry.AliasNormalized) != "" {
			keys = []string{entry.AliasNormalized}
		}
		if len(keys) == 0 && strings.TrimSpace(entry.Alias) != "" {
			keys = buildAliasLookupKeys(entry.Alias)
		}
		for _, key := range keys {
			putAlias(index.byLanguage, lang, key, alias)
		}
	}

	for _, lang := range preferredSiteLanguages {
		langAliases := index.byLanguage[lang]
		for key, alias := range langAliases {
			putAliasAny(index.byAny, key, alias)
		}
	}
	for _, langAliases := range index.byLanguage {
		for key, alias := range langAliases {
			putAliasAny(index.byAny, key, alias)
		}
	}
	return index
}

func readAliasSnapshotEntries() []aliasSnapshotEntry {
	path := strings.TrimSpace(os.Getenv(aliasSnapshotPathEnv))
	if path == "" {
		path = defaultAliasSnapshotPath
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var snapshot aliasSnapshot
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return nil
	}
	if len(snapshot.Entries) == 0 {
		return nil
	}
	return snapshot.Entries
}

func putAlias(target map[string]map[string]skillAlias, lang string, key string, alias skillAlias) {
	language := normalizeLanguageCode(lang)
	normalizedKey := normalizeSkillName(key)
	if language == "" || normalizedKey == "" {
		return
	}
	if _, ok := target[language]; !ok {
		target[language] = make(map[string]skillAlias, 256)
	}
	current, exists := target[language][normalizedKey]
	if !exists || shouldReplaceAlias(current, alias) {
		target[language][normalizedKey] = alias
	}
}

func putAliasAny(target map[string]skillAlias, key string, alias skillAlias) {
	normalizedKey := normalizeSkillName(key)
	if normalizedKey == "" {
		return
	}
	current, exists := target[normalizedKey]
	if !exists || shouldReplaceAlias(current, alias) {
		target[normalizedKey] = alias
	}
}

func shouldReplaceAlias(current skillAlias, next skillAlias) bool {
	if current.ExternalKey == "" && next.ExternalKey != "" {
		return true
	}
	if current.CanonicalNormalized == "" && next.CanonicalNormalized != "" {
		return true
	}
	if current.CanonicalNormalized == next.CanonicalNormalized {
		if next.ActivationLevel > 0 && (current.ActivationLevel == 0 || next.ActivationLevel < current.ActivationLevel) {
			return true
		}
	}
	return false
}

func pickTranslationByExternalKey(
	byExternal map[string]map[string]skillTranslationRow,
	externalKey string,
	preferredLang string,
) (skillTranslationRow, bool) {
	rows := byExternal[strings.TrimSpace(strings.ToLower(externalKey))]
	if len(rows) == 0 {
		return skillTranslationRow{}, false
	}
	if row, ok := rows[normalizeLanguageCode(preferredLang)]; ok {
		if strings.TrimSpace(row.Name) != "" {
			return row, true
		}
	}
	if row, ok := rows["en"]; ok {
		if strings.TrimSpace(row.Name) != "" {
			return row, true
		}
	}
	for _, row := range rows {
		if strings.TrimSpace(row.Name) != "" {
			return row, true
		}
	}
	return skillTranslationRow{}, false
}
