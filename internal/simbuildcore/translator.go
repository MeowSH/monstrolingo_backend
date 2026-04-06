package simbuildcore

import (
	"context"
	"slices"
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
	byLanguage, bySkillID, byEffect, byExternal := indexSkillTranslations(skillRows)
	sourceLang := detectSourceLanguage(parsed.Skills, byLanguage, activeLanguages)

	resolvedSkills := resolveSkills(parsed.Skills, sourceLang, targetLang, byLanguage, byEffect, byExternal, bySkillID)
	originalSkills, translatedSkills, unmatched := buildSkillResponsePayload(resolvedSkills)

	mode := translationModeFull
	if len(unmatched) > 0 {
		mode = translationModePartial
	}

	return &TranslateResponse{
		SourceLangDetected: sourceLang,
		TargetLang:         targetLang,
		TranslationMode:    mode,
		SkillsOriginal:     originalSkills,
		SkillsTranslated:   translatedSkills,
		UnmatchedElements:  unmatched,
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
