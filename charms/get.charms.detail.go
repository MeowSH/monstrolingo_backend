package charms

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetCharmsDetail returns detailed charm data in target language.
//
//encore:api public method=GET path=/charms/detail/:external_key
func GetCharmsDetail(ctx context.Context, external_key string, params *charmsTargetLanguageRequest) (*charmsDetailResponse, error) {
	svc, err := getCharmsService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetCharmDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &charmsDetailResponse{Data: *out}, nil
}
