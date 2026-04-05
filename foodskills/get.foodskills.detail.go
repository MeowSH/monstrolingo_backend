package foodskills

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetFoodskillsDetail returns detailed food skill data in target language.
//
//encore:api public method=GET path=/food-skills/detail/:external_key
func GetFoodskillsDetail(ctx context.Context, external_key string, params *foodskillsTargetLanguageRequest) (*foodskillsDetailResponse, error) {
	svc, err := getFoodskillsService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetFoodSkillDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &foodskillsDetailResponse{Data: *out}, nil
}
