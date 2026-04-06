package simbuildcore

import (
	"context"
	"slices"
	"sort"
	"strings"

	"github.com/google/uuid"
)

func (s *Service) TranslateSimBuild(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	if req == nil {
		return nil, invalidArgument("request body is required")
	}

	targetLang := normalizeLanguageCode(req.TargetLang)
	if targetLang == "" {
		return nil, invalidArgument("target_lang is required")
	}

	ok, err := s.repo.LanguageExists(ctx, targetLang)
	if err != nil {
		return nil, internalError("check target language support", err)
	}
	if !ok {
		return nil, unsupportedLanguage(targetLang)
	}

	parsed, err := parseAndValidateSimURL(req.URL)
	if err != nil {
		return nil, err
	}

	activeLanguages, err := s.repo.ListActiveLanguageCodes(ctx)
	if err != nil {
		return nil, internalError("list active languages", err)
	}
	if !slices.Contains(activeLanguages, "en") {
		activeLanguages = append(activeLanguages, "en")
	}

	skillRows, err := s.repo.ListSkillTranslations(ctx)
	if err != nil {
		return nil, internalError("load skill translations", err)
	}
	decorationRows, err := s.repo.ListDecorationTranslations(ctx)
	if err != nil {
		return nil, internalError("load decoration translations", err)
	}
	byLanguage, bySkillID, byEffect, byExternal := indexSkillTranslations(skillRows)
	sourceLang := detectSourceLanguage(parsed.Skills, byLanguage, activeLanguages)

	resolvedSkills := resolveSkills(parsed.Skills, sourceLang, targetLang, byLanguage, byEffect, byExternal, bySkillID)
	resolvedSkills = deduplicateResolvedSkills(resolvedSkills)
	originalSkills, translatedSkills, unmatched := buildSkillResponsePayload(resolvedSkills)
	translatedSkills = attachAssociatedJewels(resolvedSkills, translatedSkills, decorationRows, targetLang)
	setSkills, armorJewelSkills := splitTranslatedSkillsByKind(resolvedSkills, translatedSkills)

	mode := translationModeFull
	if len(unmatched) > 0 {
		mode = translationModePartial
	}

	return &TranslateResponse{
		SourceLangDetected:         sourceLang,
		TargetLang:                 targetLang,
		TranslationMode:            mode,
		SkillsOriginal:             originalSkills,
		SkillsTranslated:           translatedSkills,
		SetSkillsTranslated:        setSkills,
		ArmorJewelSkillsTranslated: armorJewelSkills,
		UnmatchedElements:          unmatched,
	}, nil
}

func indexSkillTranslations(rows []skillTranslationRow) (
	map[string]map[string][]skillTranslationRow,
	map[uuid.UUID]map[string]skillTranslationRow,
	map[string]map[string][]skillTranslationRow,
	map[string]map[string]skillTranslationRow,
) {
	byLanguage := make(map[string]map[string][]skillTranslationRow, 8)
	bySkillID := make(map[uuid.UUID]map[string]skillTranslationRow, 1024)
	byEffect := make(map[string]map[string][]skillTranslationRow, 8)
	byExternal := make(map[string]map[string]skillTranslationRow, 1024)

	for _, row := range rows {
		lang := normalizeLanguageCode(row.LanguageCode)
		nameKey := normalizeSkillName(row.Name)
		if nameKey != "" {
			if _, ok := byLanguage[lang]; !ok {
				byLanguage[lang] = make(map[string][]skillTranslationRow, 1024)
			}
			byLanguage[lang][nameKey] = append(byLanguage[lang][nameKey], row)
		}
		effectKey := normalizeSkillName(row.EffectSummary)
		if effectKey != "" {
			if _, ok := byEffect[lang]; !ok {
				byEffect[lang] = make(map[string][]skillTranslationRow, 1024)
			}
			byEffect[lang][effectKey] = append(byEffect[lang][effectKey], row)
		}

		if _, ok := bySkillID[row.SkillID]; !ok {
			bySkillID[row.SkillID] = map[string]skillTranslationRow{}
		}
		bySkillID[row.SkillID][lang] = row

		externalKey := strings.TrimSpace(strings.ToLower(row.ExternalKey))
		if externalKey != "" {
			if _, ok := byExternal[externalKey]; !ok {
				byExternal[externalKey] = map[string]skillTranslationRow{}
			}
			byExternal[externalKey][lang] = row
		}
	}

	return byLanguage, bySkillID, byEffect, byExternal
}

func resolveSkills(
	requested []requestedSkill,
	sourceLang string,
	targetLang string,
	byLanguage map[string]map[string][]skillTranslationRow,
	byEffect map[string]map[string][]skillTranslationRow,
	byExternal map[string]map[string]skillTranslationRow,
	bySkillID map[uuid.UUID]map[string]skillTranslationRow,
) []resolvedSkill {
	out := make([]resolvedSkill, 0, len(requested))
	for _, req := range requested {
		resolved := resolvedSkill{
			Requested: req,
		}

		nameKey := normalizeSkillName(req.BaseName)
		candidates := byLanguage[sourceLang][nameKey]
		if len(candidates) == 0 {
			candidates = byEffect[sourceLang][nameKey]
		}
		if len(candidates) == 0 && sourceLang != "en" {
			candidates = byLanguage["en"][nameKey]
		}
		if len(candidates) == 0 && sourceLang != "en" {
			candidates = byEffect["en"][nameKey]
		}
		if len(candidates) == 0 {
			if alias, ok := resolveSkillAlias(nameKey, sourceLang); ok {
				if strings.TrimSpace(alias.ExternalKey) != "" {
					if aliasRow, ok := pickTranslationByExternalKey(byExternal, alias.ExternalKey, sourceLang); ok {
						candidates = []skillTranslationRow{aliasRow}
					} else if aliasRow, ok := pickTranslationByExternalKey(byExternal, alias.ExternalKey, "en"); ok {
						candidates = []skillTranslationRow{aliasRow}
					}
				}
				if len(candidates) == 0 {
					canonicalKey := aliasCanonicalLookupKey(alias)
					if canonicalKey != "" {
						candidates = byLanguage[sourceLang][canonicalKey]
					}
				}
				if len(candidates) == 0 && sourceLang != "en" {
					canonicalKey := aliasCanonicalLookupKey(alias)
					if canonicalKey != "" {
						candidates = byLanguage["en"][canonicalKey]
					}
				}
				if len(candidates) == 0 {
					canonicalKey := aliasCanonicalLookupKey(alias)
					if canonicalKey != "" {
						for _, langMap := range byLanguage {
							if len(langMap[canonicalKey]) > 0 {
								candidates = langMap[canonicalKey]
								break
							}
						}
					}
				}
			}
		}
		if len(candidates) == 0 {
			out = append(out, resolved)
			continue
		}

		chosen := chooseBestSkillCandidate(candidates, req.RequestedLevel)
		resolved.SkillID = chosen.SkillID
		resolved.ExternalKey = chosen.ExternalKey
		resolved.MaxLevel = chosen.MaxLevel
		resolved.IsSetBonus = chosen.IsSetBonus
		resolved.SourceName = strings.TrimSpace(chosen.Name)
		resolved.Resolved = true

		byLang := bySkillID[chosen.SkillID]
		targetKey := normalizeLanguageCode(targetLang)
		if target, ok := byLang[targetKey]; ok && strings.TrimSpace(target.Name) != "" {
			resolved.TargetName = strings.TrimSpace(target.Name)
			resolved.Translated = true
		}

		english, enOK := chooseSkillTranslation(byLang, "en")
		if enOK {
			resolved.EnglishName = strings.TrimSpace(english.Name)
		}
		if !resolved.Translated && resolved.EnglishName != "" {
			resolved.TargetName = resolved.EnglishName
		}
		if resolved.TargetName == "" {
			resolved.TargetName = req.OriginalText
		}

		out = append(out, resolved)
	}
	return out
}

func chooseBestSkillCandidate(candidates []skillTranslationRow, requestedLevel int16) skillTranslationRow {
	best := candidates[0]
	bestScore := skillCandidateScore(best, requestedLevel)
	for _, candidate := range candidates[1:] {
		score := skillCandidateScore(candidate, requestedLevel)
		if score > bestScore {
			best = candidate
			bestScore = score
			continue
		}
		if score == bestScore && candidate.ExternalKey < best.ExternalKey {
			best = candidate
		}
	}
	return best
}

func skillCandidateScore(candidate skillTranslationRow, requestedLevel int16) int {
	diff := int(candidate.MaxLevel - requestedLevel)
	if diff < 0 {
		diff = -diff + 5
	}
	return 10_000 - diff
}

func chooseSkillTranslation(byLang map[string]skillTranslationRow, lang string) (skillTranslationRow, bool) {
	if len(byLang) == 0 {
		return skillTranslationRow{}, false
	}
	if row, ok := byLang[normalizeLanguageCode(lang)]; ok {
		if strings.TrimSpace(row.Name) != "" {
			return row, true
		}
	}
	if row, ok := byLang["en"]; ok {
		if strings.TrimSpace(row.Name) != "" {
			return row, true
		}
	}
	for _, row := range byLang {
		if strings.TrimSpace(row.Name) != "" {
			return row, true
		}
	}
	return skillTranslationRow{}, false
}

func deduplicateResolvedSkills(rows []resolvedSkill) []resolvedSkill {
	if len(rows) <= 1 {
		return rows
	}

	out := make([]resolvedSkill, 0, len(rows))
	indexByKey := make(map[string]int, len(rows))
	for _, row := range rows {
		key := resolvedSkillDedupeKey(row)
		if key == "" {
			out = append(out, row)
			continue
		}
		if idx, ok := indexByKey[key]; ok {
			out[idx] = mergeResolvedSkill(out[idx], row)
			continue
		}
		indexByKey[key] = len(out)
		out = append(out, row)
	}
	return out
}

func resolvedSkillDedupeKey(row resolvedSkill) string {
	if row.Resolved {
		externalKey := strings.TrimSpace(strings.ToLower(row.ExternalKey))
		if externalKey != "" {
			return "resolved:" + externalKey
		}
	}
	base := normalizeSkillName(row.Requested.BaseName)
	if base == "" {
		base = normalizeSkillName(row.Requested.OriginalText)
	}
	if base == "" {
		return ""
	}
	return "requested:" + base
}

func mergeResolvedSkill(existing resolvedSkill, incoming resolvedSkill) resolvedSkill {
	merged := existing

	if incoming.Requested.RequestedLevel > merged.Requested.RequestedLevel {
		merged.Requested.RequestedLevel = incoming.Requested.RequestedLevel
	}

	if !merged.Resolved && incoming.Resolved {
		merged.Resolved = true
		merged.SkillID = incoming.SkillID
		merged.ExternalKey = incoming.ExternalKey
		merged.MaxLevel = incoming.MaxLevel
		merged.IsSetBonus = incoming.IsSetBonus
		merged.SourceName = incoming.SourceName
		merged.TargetName = incoming.TargetName
		merged.EnglishName = incoming.EnglishName
		merged.Translated = incoming.Translated
		return merged
	}

	merged.IsSetBonus = merged.IsSetBonus || incoming.IsSetBonus
	merged.Translated = merged.Translated || incoming.Translated
	if merged.SkillID == uuid.Nil && incoming.SkillID != uuid.Nil {
		merged.SkillID = incoming.SkillID
	}
	if strings.TrimSpace(merged.ExternalKey) == "" && strings.TrimSpace(incoming.ExternalKey) != "" {
		merged.ExternalKey = incoming.ExternalKey
	}
	if incoming.MaxLevel > merged.MaxLevel {
		merged.MaxLevel = incoming.MaxLevel
	}
	if strings.TrimSpace(merged.SourceName) == "" && strings.TrimSpace(incoming.SourceName) != "" {
		merged.SourceName = incoming.SourceName
	}
	if strings.TrimSpace(merged.EnglishName) == "" && strings.TrimSpace(incoming.EnglishName) != "" {
		merged.EnglishName = incoming.EnglishName
	}

	if !existing.Translated && incoming.Translated && strings.TrimSpace(incoming.TargetName) != "" {
		merged.TargetName = incoming.TargetName
	} else if strings.TrimSpace(merged.TargetName) == "" && strings.TrimSpace(incoming.TargetName) != "" {
		merged.TargetName = incoming.TargetName
	}

	return merged
}

func buildSkillResponsePayload(resolved []resolvedSkill) (
	[]OriginalSkill,
	[]TranslatedSkill,
	[]UnmatchedElement,
) {
	original := make([]OriginalSkill, 0, len(resolved))
	translated := make([]TranslatedSkill, 0, len(resolved))
	unmatched := make([]UnmatchedElement, 0)

	for _, row := range resolved {
		original = append(original, OriginalSkill{
			Text:           row.Requested.OriginalText,
			Name:           row.Requested.BaseName,
			RequestedLevel: row.Requested.RequestedLevel,
		})

		outName := row.TargetName
		if strings.TrimSpace(outName) == "" {
			outName = row.Requested.OriginalText
		}
		translated = append(translated, TranslatedSkill{
			OriginalText:     row.Requested.OriginalText,
			OriginalName:     row.Requested.BaseName,
			RequestedLevel:   row.Requested.RequestedLevel,
			Name:             outName,
			Translated:       row.Resolved && row.Translated,
			SkillExternalKey: row.ExternalKey,
		})

		if !row.Resolved {
			unmatched = append(unmatched, UnmatchedElement{
				Kind:   "skill",
				Value:  row.Requested.OriginalText,
				Reason: "skill_not_found",
			})
			continue
		}
		if !row.Translated {
			unmatched = append(unmatched, UnmatchedElement{
				Kind:   "skill_translation",
				Value:  row.Requested.OriginalText,
				Reason: "target_translation_missing",
			})
		}

	}

	return original, translated, unmatched
}

func splitTranslatedSkillsByKind(
	resolved []resolvedSkill,
	translated []TranslatedSkill,
) ([]TranslatedSkill, []TranslatedSkill) {
	setSkills := make([]TranslatedSkill, 0, len(translated))
	armorJewelSkills := make([]TranslatedSkill, 0, len(translated))

	limit := len(translated)
	if len(resolved) < limit {
		limit = len(resolved)
	}
	for i := 0; i < limit; i++ {
		if !resolved[i].Resolved {
			continue
		}
		if isSetOrGroupResolvedSkill(resolved[i]) {
			setSkills = append(setSkills, translated[i])
			continue
		}
		armorJewelSkills = append(armorJewelSkills, translated[i])
	}

	return setSkills, armorJewelSkills
}

func attachAssociatedJewels(
	resolved []resolvedSkill,
	translated []TranslatedSkill,
	decorationRows []decorationTranslationRow,
	targetLang string,
) []TranslatedSkill {
	if len(translated) == 0 || len(decorationRows) == 0 {
		return translated
	}

	bySkillID := associatedJewelsBySkillID(decorationRows, targetLang)
	if len(bySkillID) == 0 {
		return translated
	}

	out := make([]TranslatedSkill, len(translated))
	copy(out, translated)

	limit := len(out)
	if len(resolved) < limit {
		limit = len(resolved)
	}
	for i := 0; i < limit; i++ {
		if !resolved[i].Resolved {
			continue
		}
		jewels := bySkillID[resolved[i].SkillID]
		if len(jewels) == 0 {
			continue
		}
		copied := make([]AssociatedJewel, len(jewels))
		copy(copied, jewels)
		out[i].AssociatedJewels = copied
	}
	return out
}

func associatedJewelsBySkillID(rows []decorationTranslationRow, targetLang string) map[uuid.UUID][]AssociatedJewel {
	type decorationAggregate struct {
		ExternalKey string
		SlotSize    int16
		Rarity      int16
		SkillLevel  int16
		NamesByLang map[string]string
	}

	bySkill := make(map[uuid.UUID]map[string]*decorationAggregate, 512)
	for _, row := range rows {
		if row.SkillID == uuid.Nil {
			continue
		}
		externalKey := strings.TrimSpace(strings.ToLower(row.DecorationExternalKey))
		if externalKey == "" {
			continue
		}
		lang := normalizeLanguageCode(row.LanguageCode)
		name := strings.TrimSpace(row.Name)

		decMap, ok := bySkill[row.SkillID]
		if !ok {
			decMap = make(map[string]*decorationAggregate, 8)
			bySkill[row.SkillID] = decMap
		}
		agg, ok := decMap[externalKey]
		if !ok {
			agg = &decorationAggregate{
				ExternalKey: externalKey,
				SlotSize:    row.SlotSize,
				Rarity:      row.Rarity,
				SkillLevel:  row.SkillLevel,
				NamesByLang: make(map[string]string, 4),
			}
			decMap[externalKey] = agg
		}
		if row.SkillLevel > agg.SkillLevel {
			agg.SkillLevel = row.SkillLevel
		}
		if lang != "" && name != "" {
			agg.NamesByLang[lang] = name
		}
	}

	out := make(map[uuid.UUID][]AssociatedJewel, len(bySkill))
	for skillID, decMap := range bySkill {
		jewels := make([]AssociatedJewel, 0, len(decMap))
		for _, agg := range decMap {
			name := pickDecorationNameForLanguage(agg.NamesByLang, targetLang)
			if name == "" {
				name = agg.ExternalKey
			}
			jewels = append(jewels, AssociatedJewel{
				DecorationExternalKey: agg.ExternalKey,
				Name:                  name,
				SlotSize:              agg.SlotSize,
				Rarity:                agg.Rarity,
				SkillLevel:            agg.SkillLevel,
			})
		}
		sort.Slice(jewels, func(i, j int) bool {
			if jewels[i].SlotSize != jewels[j].SlotSize {
				return jewels[i].SlotSize < jewels[j].SlotSize
			}
			if jewels[i].SkillLevel != jewels[j].SkillLevel {
				return jewels[i].SkillLevel > jewels[j].SkillLevel
			}
			if jewels[i].Name != jewels[j].Name {
				return jewels[i].Name < jewels[j].Name
			}
			return jewels[i].DecorationExternalKey < jewels[j].DecorationExternalKey
		})
		out[skillID] = jewels
	}
	return out
}

func pickDecorationNameForLanguage(namesByLang map[string]string, targetLang string) string {
	if len(namesByLang) == 0 {
		return ""
	}
	target := normalizeLanguageCode(targetLang)
	if target != "" {
		if value := strings.TrimSpace(namesByLang[target]); value != "" {
			return value
		}
	}
	if value := strings.TrimSpace(namesByLang["en"]); value != "" {
		return value
	}

	langs := make([]string, 0, len(namesByLang))
	for lang := range namesByLang {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	for _, lang := range langs {
		if value := strings.TrimSpace(namesByLang[lang]); value != "" {
			return value
		}
	}
	return ""
}

func isSetOrGroupResolvedSkill(skill resolvedSkill) bool {
	if skill.IsSetBonus {
		return true
	}
	alias, ok := resolveSkillAlias(skill.Requested.BaseName, "")
	if !ok {
		return false
	}
	switch strings.TrimSpace(strings.ToLower(alias.Category)) {
	case "set", "group":
		return true
	default:
		return false
	}
}
