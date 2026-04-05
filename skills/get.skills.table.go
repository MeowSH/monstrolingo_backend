package skills

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetSkillsTable returns bilingual paginated rows for skills.
//
//encore:api public method=GET path=/skills/table
func GetSkillsTable(ctx context.Context, req *skillsTableRequest) (*skillsTableResponse, error) {
	svc, err := getSkillsService()
	if err != nil {
		return nil, err
	}
	var coreReq *catalogcore.CategoryTableRequest
	if req != nil {
		coreReq = &catalogcore.CategoryTableRequest{
			SourceLang: req.SourceLang,
			TargetLang: req.TargetLang,
			Page:       req.Page,
			Limit:      req.Limit,
		}
	}
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategorySkills, coreReq)
	if err != nil {
		return nil, err
	}
	return &skillsTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
