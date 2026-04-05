package kinsects

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetKinsectsDetail returns detailed kinsect data in target language.
//
//encore:api public method=GET path=/kinsects/detail/:external_key
func GetKinsectsDetail(ctx context.Context, external_key string, params *kinsectsTargetLanguageRequest) (*kinsectsDetailResponse, error) {
	svc, err := getKinsectsService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetKinsectDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &kinsectsDetailResponse{Data: *out}, nil
}
