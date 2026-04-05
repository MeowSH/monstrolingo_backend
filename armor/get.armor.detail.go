package armor

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetArmorDetail returns detailed armor data in target language.
//
//encore:api public method=GET path=/armor/detail/:external_key
func GetArmorDetail(ctx context.Context, external_key string, params *armorTargetLanguageRequest) (*armorDetailResponse, error) {
	svc, err := getArmorService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetArmorDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &armorDetailResponse{Data: *out}, nil
}
