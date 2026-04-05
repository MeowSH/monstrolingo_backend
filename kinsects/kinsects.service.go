package kinsects

import "monstrolingo_backend/internal/catalogcore"

type kinsectsTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type kinsectsTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type kinsectsTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type kinsectsDetailResponse struct {
	Data catalogcore.KinsectDetailResponse `json:"data"`
}

func getKinsectsService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
