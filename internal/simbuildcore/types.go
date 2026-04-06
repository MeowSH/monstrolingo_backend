package simbuildcore

import "github.com/google/uuid"

const (
	translationModeFull    = "full"
	translationModePartial = "partial"
)

type TranslateRequest struct {
	URL        string
	TargetLang string
}

type TranslateResponse struct {
	SourceLangDetected         string             `json:"source_lang_detected"`
	TargetLang                 string             `json:"target_lang"`
	TranslationMode            string             `json:"translation_mode"`
	SkillsOriginal             []OriginalSkill    `json:"skills_original"`
	SkillsTranslated           []TranslatedSkill  `json:"skills_translated"`
	SetSkillsTranslated        []TranslatedSkill  `json:"set_skills_translated"`
	ArmorJewelSkillsTranslated []TranslatedSkill  `json:"armor_jewel_skills_translated"`
	UnmatchedElements          []UnmatchedElement `json:"unmatched_elements"`
}

type OriginalSkill struct {
	Text           string `json:"text"`
	Name           string `json:"name"`
	RequestedLevel int16  `json:"requested_level"`
}

type TranslatedSkill struct {
	OriginalText     string            `json:"original_text"`
	OriginalName     string            `json:"original_name"`
	RequestedLevel   int16             `json:"requested_level"`
	Name             string            `json:"name"`
	Translated       bool              `json:"translated"`
	SkillExternalKey string            `json:"skill_external_key,omitempty"`
	AssociatedJewels []AssociatedJewel `json:"associated_jewels,omitempty"`
}

type AssociatedJewel struct {
	DecorationExternalKey string `json:"decoration_external_key"`
	Name                  string `json:"name"`
	SlotSize              int16  `json:"slot_size"`
	Rarity                int16  `json:"rarity"`
	SkillLevel            int16  `json:"skill_level"`
}

type UnmatchedElement struct {
	Kind   string `json:"kind"`
	Value  string `json:"value"`
	Reason string `json:"reason"`
}

type parsedSimURL struct {
	RawURL            string
	Skills            []requestedSkill
	WeaponSkillsText  string
	WeaponGroupText   string
	WeaponSetText     string
	WeaponSetSkills   []requestedSkill
	WeaponExtraSkills []requestedSkill
}

type requestedSkill struct {
	OriginalText   string
	BaseName       string
	RequestedLevel int16
}

type resolvedSkill struct {
	Requested   requestedSkill
	SkillID     uuid.UUID
	ExternalKey string
	MaxLevel    int16
	IsSetBonus  bool
	SourceName  string
	TargetName  string
	EnglishName string
	Resolved    bool
	Translated  bool
}

type skillTranslationRow struct {
	SkillID       uuid.UUID
	ExternalKey   string
	MaxLevel      int16
	IsSetBonus    bool
	LanguageCode  string
	Name          string
	EffectSummary string
}

type decorationTranslationRow struct {
	SkillID               uuid.UUID
	DecorationExternalKey string
	SlotSize              int16
	Rarity                int16
	SkillLevel            int16
	LanguageCode          string
	Name                  string
}
