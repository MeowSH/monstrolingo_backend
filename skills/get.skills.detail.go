package skills

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetSkillsDetail returns detailed skill data in target language.
//
//encore:api public method=GET path=/skills/detail/:external_key
func GetSkillsDetail(ctx context.Context, external_key string, params *skillsTargetLanguageRequest) (*skillsDetailResponse, error) {
	svc, err := getSkillsService()
	if err != nil {
		return nil, err
	}
	req := &catalogcore.CategoryDetailRequest{ExternalKey: external_key}
	if params != nil {
		req.TargetLang = params.TargetLang
	}
	out, err := svc.GetSkillDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return &skillsDetailResponse{Data: *out}, nil
}
