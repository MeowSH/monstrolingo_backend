package simbuildcore

import (
	"slices"
	"sort"
	"strings"
)

var preferredSiteLanguages = []string{
	"en",
	"ja",
	"ko",
	"zh-hans",
	"zh-hant",
}

func detectSourceLanguage(
	skills []requestedSkill,
	byLanguage map[string]map[string][]skillTranslationRow,
	available []string,
) string {
	candidates := detectionLanguageCandidates(available)
	if len(candidates) == 0 {
		candidates = sortedLanguageKeys(byLanguage)
	}
	if len(candidates) == 0 {
		return "en"
	}

	bestLang := candidates[0]
	bestScore := -1
	bestNameHits := -1
	for _, lang := range candidates {
		nameMap := byLanguage[lang]
		nameHits := 0
		aliasHits := 0
		for _, skill := range skills {
			key := normalizeSkillName(skill.BaseName)
			if key == "" {
				continue
			}
			if len(nameMap[key]) > 0 {
				nameHits++
				continue
			}
			if hasSkillAliasForLanguage(lang, key) {
				aliasHits++
			}
		}
		// Prioritize direct DB name matches first, then alias-only hints.
		score := (nameHits * 10) + (aliasHits * 3)
		if score > bestScore {
			bestLang = lang
			bestScore = score
			bestNameHits = nameHits
			continue
		}
		if score == bestScore {
			if nameHits > bestNameHits {
				bestLang = lang
				bestNameHits = nameHits
				continue
			}
			if nameHits == bestNameHits && languageTieBreak(lang, bestLang) {
				bestLang = lang
			}
		}
	}

	if bestScore <= 0 && slices.Contains(candidates, "en") {
		return "en"
	}
	return bestLang
}

func detectionLanguageCandidates(available []string) []string {
	if len(available) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(available))
	seen := map[string]struct{}{}
	for _, code := range available {
		n := normalizeLanguageCode(code)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		normalized = append(normalized, n)
	}

	siteFirst := make([]string, 0, len(preferredSiteLanguages))
	for _, preferred := range preferredSiteLanguages {
		if _, ok := seen[preferred]; ok {
			siteFirst = append(siteFirst, preferred)
		}
	}
	if len(siteFirst) > 0 {
		return siteFirst
	}
	sort.Strings(normalized)
	return normalized
}

func languageTieBreak(candidate string, current string) bool {
	candidateRank := languagePriority(candidate)
	currentRank := languagePriority(current)
	if candidateRank != currentRank {
		return candidateRank < currentRank
	}
	return candidate < current
}

func languagePriority(code string) int {
	n := normalizeLanguageCode(code)
	for idx, preferred := range preferredSiteLanguages {
		if n == preferred {
			return idx
		}
	}
	return 100
}

func normalizeLanguageCode(raw string) string {
	code := strings.TrimSpace(strings.ToLower(raw))
	code = strings.ReplaceAll(code, "_", "-")
	return code
}

func sortedLanguageKeys(byLanguage map[string]map[string][]skillTranslationRow) []string {
	keys := make([]string, 0, len(byLanguage))
	for code := range byLanguage {
		keys = append(keys, code)
	}
	sort.Strings(keys)
	return keys
}
