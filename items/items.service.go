package items

import "monstrolingo_backend/internal/catalogcore"

type itemsTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type itemsTargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type itemsTableResponse struct {
	Items      []catalogcore.CategoryTableRow `json:"items"`
	Pagination catalogcore.Pagination         `json:"pagination"`
}

type itemsDetailResponse struct {
	Data catalogcore.ItemDetailResponse `json:"data"`
}

func getItemsService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
