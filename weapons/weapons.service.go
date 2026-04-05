package weapons

import "monstrolingo_backend/internal/catalogcore"

type weaponsTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type weaponsTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type weaponsTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type weaponsDetailResponse struct {
	Data catalogcore.WeaponDetailResponse `json:"data"`
}

func getWeaponsService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
