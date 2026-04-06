package linkbuild

import "monstrolingo_backend/internal/simbuildcore"

type linkbuildTranslateRequest struct {
	URL        string `json:"url"`
	TargetLang string `json:"target_lang"`
}

type linkbuildTranslateResponse struct {
	SourceLangDetected string                          `json:"source_lang_detected"`
	TargetLang         string                          `json:"target_lang"`
	TranslationMode    string                          `json:"translation_mode"`
	SkillsOriginal     []simbuildcore.OriginalSkill    `json:"skills_original"`
	SkillsTranslated   []simbuildcore.TranslatedSkill  `json:"skills_translated"`
	UnmatchedElements  []simbuildcore.UnmatchedElement `json:"unmatched_elements"`
}

func getLinkbuildService() (*simbuildcore.Service, error) {
	return simbuildcore.GetService()
}
