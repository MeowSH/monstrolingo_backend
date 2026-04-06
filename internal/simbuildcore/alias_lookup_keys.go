package simbuildcore

import (
	"regexp"
	"strings"
)

var (
	trailingUnicodeRomanRegex = regexp.MustCompile(`[ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩ]+$`)
	trailingASCIIRomanRegex   = regexp.MustCompile(`\s+(?i:(?:x|ix|iv|v?i{1,3}))$`)
)

func buildAliasLookupKeys(alias string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, 4)
	addKey := func(raw string) {
		key := normalizeSkillName(raw)
		if key == "" {
			return
		}
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}

	addKey(alias)

	parsed := parseRequestedSkill(alias)
	addKey(parsed.BaseName)

	addKey(stripTrailingUnicodeRoman(alias))
	addKey(stripTrailingASCIIRoman(alias))

	if len(out) == 0 {
		addKey(strings.TrimSpace(alias))
	}
	return out
}

func stripTrailingUnicodeRoman(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	return strings.TrimSpace(trailingUnicodeRomanRegex.ReplaceAllString(trimmed, ""))
}

func stripTrailingASCIIRoman(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	return strings.TrimSpace(trailingASCIIRomanRegex.ReplaceAllString(trimmed, ""))
}
