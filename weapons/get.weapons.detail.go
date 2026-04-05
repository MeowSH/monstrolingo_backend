package weapons

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetWeaponsDetail returns detailed weapon data in target language.
//
//encore:api public method=GET path=/weapons/detail/:external_key
func GetWeaponsDetail(ctx context.Context, external_key string, params *weaponsTargetLanguageRequest) (*weaponsDetailResponse, error) {
	svc, err := getWeaponsService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetWeaponDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &weaponsDetailResponse{Data: *out}, nil
}
