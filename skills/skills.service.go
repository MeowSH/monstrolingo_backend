package skills

import "monstrolingo_backend/internal/catalogcore"

type skillsTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type skillsTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type skillsTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type skillsDetailResponse struct {
	Data catalogcore.SkillDetailResponse `json:"data"`
}

func getSkillsService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
