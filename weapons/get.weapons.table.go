package weapons

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetWeaponsTable returns bilingual paginated rows for weapons.
//
//encore:api public method=GET path=/weapons/table
func GetWeaponsTable(ctx context.Context, req *weaponsTableRequest) (*weaponsTableResponse, error) {
	svc, err := getWeaponsService()
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
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategoryWeapons, coreReq)
	if err != nil {
		return nil, err
	}
	return &weaponsTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
