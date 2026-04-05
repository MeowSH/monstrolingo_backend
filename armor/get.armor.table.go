package armor

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetArmorTable returns bilingual paginated rows for armor pieces.
//
//encore:api public method=GET path=/armor/table
func GetArmorTable(ctx context.Context, req *armorTableRequest) (*armorTableResponse, error) {
	svc, err := getArmorService()
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
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategoryArmor, coreReq)
	if err != nil {
		return nil, err
	}
	return &armorTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
