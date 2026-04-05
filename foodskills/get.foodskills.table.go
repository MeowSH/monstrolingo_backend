package foodskills

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetFoodskillsTable returns bilingual paginated rows for food skills.
//
//encore:api public method=GET path=/food-skills/table
func GetFoodskillsTable(ctx context.Context, req *foodskillsTableRequest) (*foodskillsTableResponse, error) {
	svc, err := getFoodskillsService()
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
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategoryFoodSkills, coreReq)
	if err != nil {
		return nil, err
	}
	return &foodskillsTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
