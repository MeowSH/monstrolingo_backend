package kinsects

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetKinsectsTable returns bilingual paginated rows for kinsects.
//
//encore:api public method=GET path=/kinsects/table
func GetKinsectsTable(ctx context.Context, req *kinsectsTableRequest) (*kinsectsTableResponse, error) {
	svc, err := getKinsectsService()
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
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategoryKinsects, coreReq)
	if err != nil {
		return nil, err
	}
	return &kinsectsTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
