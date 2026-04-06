package simbuildcore

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	simURLHostsEnv = "SIMBUILD_ALLOWED_SIM_HOSTS"
)

func parseAndValidateSimURL(raw string) (parsedSimURL, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return parsedSimURL{}, invalidArgument("url is required")
	}

	parsed, err := parseURLWithOptionalScheme(trimmed)
	if err != nil {
		return parsedSimURL{}, invalidArgument("url is invalid")
	}
	if !isSupportedSimURL(parsed) {
		return parsedSimURL{}, invalidArgument("url must target a supported simulator host and /sim/ path")
	}
	if strings.TrimSpace(parsed.Fragment) == "" {
		return parsedSimURL{}, invalidArgument("url fragment is required")
	}

	fragmentValues, err := url.ParseQuery(parsed.Fragment)
	if err != nil {
		return parsedSimURL{}, invalidArgument("url fragment is invalid")
	}

	skillsRaw := strings.TrimSpace(fragmentValues.Get("skills"))
	if skillsRaw == "" {
		return parsedSimURL{}, invalidArgument("url fragment must include skills")
	}

	out := parsedSimURL{
		RawURL:           trimmed,
		WeaponSkillsText: strings.TrimSpace(fragmentValues.Get("ws")),
		WeaponGroupText:  strings.TrimSpace(fragmentValues.Get("wgs")),
		WeaponSetText:    strings.TrimSpace(fragmentValues.Get("wss")),
	}

	tokens := splitSkillTokens(skillsRaw)
	out.Skills = make([]requestedSkill, 0, len(tokens)+2)
	for _, token := range tokens {
		skill := parseRequestedSkill(token)
		if skill.BaseName == "" {
			continue
		}
		out.Skills = append(out.Skills, skill)
	}

	if out.WeaponGroupText != "" {
		parsedSkill := parseRequestedSkill(out.WeaponGroupText)
		if parsedSkill.BaseName != "" {
			out.WeaponSetSkills = append(out.WeaponSetSkills, parsedSkill)
			out.Skills = append(out.Skills, parsedSkill)
		}
	}
	if out.WeaponSetText != "" {
		parsedSkill := parseRequestedSkill(out.WeaponSetText)
		if parsedSkill.BaseName != "" {
			out.WeaponSetSkills = append(out.WeaponSetSkills, parsedSkill)
			out.Skills = append(out.Skills, parsedSkill)
		}
	}
	for _, token := range splitSkillTokens(out.WeaponSkillsText) {
		parsedSkill := parseRequestedSkill(token)
		if parsedSkill.BaseName == "" {
			continue
		}
		out.WeaponExtraSkills = append(out.WeaponExtraSkills, parsedSkill)
		out.Skills = append(out.Skills, parsedSkill)
	}

	if len(out.Skills) == 0 {
		return parsedSimURL{}, invalidArgument("url does not contain usable skills")
	}

	return out, nil
}

func parseURLWithOptionalScheme(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err == nil && parsed.Host != "" {
		return parsed, nil
	}
	trimmed := strings.TrimSpace(raw)
	if looksLikeHostPath(trimmed) {
		return url.Parse("https://" + strings.TrimPrefix(trimmed, "//"))
	}
	return parsed, err
}

func isSupportedSimURL(parsed *url.URL) bool {
	if parsed == nil {
		return false
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" {
		return false
	}

	allowedHosts, configured := configuredSimulatorHosts()
	if !configured {
		return false
	}
	if _, ok := allowedHosts[host]; !ok {
		return false
	}

	path := strings.TrimSpace(parsed.Path)
	if path == "/sim" || path == "/sim/" {
		return true
	}
	return strings.HasPrefix(path, "/sim/")
}

func looksLikeHostPath(raw string) bool {
	if raw == "" {
		return false
	}
	if strings.Contains(raw, "://") {
		return false
	}
	if strings.HasPrefix(raw, "/") {
		return false
	}
	idx := strings.Index(raw, "/")
	if idx <= 0 {
		return false
	}
	host := raw[:idx]
	if strings.ContainsAny(host, " ?#") {
		return false
	}
	return strings.Contains(host, ".")
}

func configuredSimulatorHosts() (map[string]struct{}, bool) {
	raw := strings.TrimSpace(os.Getenv(simURLHostsEnv))
	if raw == "" {
		return nil, false
	}
	out := make(map[string]struct{}, 4)
	for _, part := range strings.Split(raw, ",") {
		host := normalizeConfiguredHost(part)
		if host == "" {
			continue
		}
		out[host] = struct{}{}
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

func normalizeConfiguredHost(raw string) string {
	value := strings.TrimSpace(strings.ToLower(raw))
	if value == "" {
		return ""
	}
	value = strings.TrimPrefix(value, "https://")
	value = strings.TrimPrefix(value, "http://")
	if idx := strings.IndexAny(value, "/?#"); idx >= 0 {
		value = value[:idx]
	}

	if strings.Contains(value, ":") {
		if host, _, err := net.SplitHostPort(value); err == nil {
			value = host
		}
	}
	value = strings.Trim(value, "[]")
	return strings.TrimSpace(value)
}

func splitSkillTokens(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func parseRequestedSkill(raw string) requestedSkill {
	original := strings.TrimSpace(raw)
	if original == "" {
		return requestedSkill{}
	}

	working := original
	level := int16(1)

	if idx := strings.LastIndex(strings.ToLower(working), " lv"); idx >= 0 {
		lvPart := strings.TrimSpace(working[idx+1:])
		if strings.HasPrefix(strings.ToLower(lvPart), "lv") {
			n, err := strconv.Atoi(strings.TrimSpace(lvPart[2:]))
			if err == nil && n > 0 {
				level = int16(n)
				working = strings.TrimSpace(working[:idx])
			}
		}
	}

	fields := strings.Fields(working)
	if len(fields) >= 2 {
		if roman := romanToInt(fields[len(fields)-1]); roman > 0 {
			level = int16(roman)
			working = strings.Join(fields[:len(fields)-1], " ")
		}
	}

	baseName := strings.TrimSpace(working)
	if baseName == "" {
		baseName = original
	}

	return requestedSkill{
		OriginalText:   original,
		BaseName:       baseName,
		RequestedLevel: level,
	}
}

func romanToInt(raw string) int {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "I":
		return 1
	case "II":
		return 2
	case "III":
		return 3
	case "IV":
		return 4
	case "V":
		return 5
	case "VI":
		return 6
	case "VII":
		return 7
	case "VIII":
		return 8
	case "IX":
		return 9
	case "X":
		return 10
	default:
		return 0
	}
}

func normalizeSkillName(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	replacer := strings.NewReplacer("’", "'", "‘", "'", "`", "'")
	value = replacer.Replace(value)
	value = strings.ToLower(value)
	value = strings.Join(strings.Fields(value), " ")
	return value
}
