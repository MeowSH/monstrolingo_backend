package charms

import "monstrolingo_backend/internal/catalogcore"

type charmsTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type charmsTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type charmsTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type charmsDetailResponse struct {
	Data catalogcore.CharmDetailResponse `json:"data"`
}

func getCharmsService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
