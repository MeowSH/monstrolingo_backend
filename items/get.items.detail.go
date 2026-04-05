package items

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetItemsDetail returns detailed item data in target language.
//
//encore:api public method=GET path=/items/detail/:external_key
func GetItemsDetail(ctx context.Context, external_key string, params *itemsTargetLanguageRequest) (*itemsDetailResponse, error) {
	svc, err := getItemsService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetItemDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &itemsDetailResponse{Data: *out}, nil
}
