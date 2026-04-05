package decorations

import "monstrolingo_backend/internal/catalogcore"

type decorationsTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type decorationsTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type decorationsTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type decorationsDetailResponse struct {
	Data catalogcore.DecorationDetailResponse `json:"data"`
}

func getDecorationsService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
