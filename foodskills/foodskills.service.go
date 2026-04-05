package foodskills

import "monstrolingo_backend/internal/catalogcore"

type foodskillsTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type foodskillsTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type foodskillsTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type foodskillsDetailResponse struct {
	Data catalogcore.FoodSkillDetailResponse `json:"data"`
}

func getFoodskillsService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
