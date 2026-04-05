package languages

import "monstrolingo_backend/internal/catalogcore"

type languagesListResponse struct {
	Languages []catalogcore.LanguageOption `json:"languages"`
}

func getLanguagesService() (*catalogcore.Service, error) {
	return catalogcore.GetService()
}
