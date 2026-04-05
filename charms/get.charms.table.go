package charms

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetCharmsTable returns bilingual paginated rows for charms.
//
//encore:api public method=GET path=/charms/table
func GetCharmsTable(ctx context.Context, req *charmsTableRequest) (*charmsTableResponse, error) {
	svc, err := getCharmsService()
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
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategoryCharms, coreReq)
	if err != nil {
		return nil, err
	}
	return &charmsTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
