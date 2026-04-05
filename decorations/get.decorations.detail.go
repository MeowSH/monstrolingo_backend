package decorations

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetDecorationsDetail returns detailed decoration data in target language.
//
//encore:api public method=GET path=/decorations/detail/:external_key
func GetDecorationsDetail(ctx context.Context, external_key string, params *decorationsTargetLanguageRequest) (*decorationsDetailResponse, error) {
	svc, err := getDecorationsService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetDecorationDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &decorationsDetailResponse{Data: *out}, nil
}
