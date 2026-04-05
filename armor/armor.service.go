package armor

import "monstrolingo_backend/internal/catalogcore"

type armorTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type armorTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type armorTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type armorDetailResponse struct {
	Data catalogcore.ArmorDetailResponse `json:"data"`
}

func getArmorService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
