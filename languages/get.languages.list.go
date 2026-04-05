package languages

import (
	"context"
)

// GetLanguagesList returns active language codes and labels.
//
//encore:api public method=GET path=/languages
func GetLanguagesList(ctx context.Context) (*languagesListResponse, error) {
	svc, err := getLanguagesService()
	if err != nil {
		return nil, err
	}
	out, err := svc.ListLanguages(ctx)
	if err != nil {
		return nil, err
	}
	return &languagesListResponse{Languages: out.Languages}, nil
}
