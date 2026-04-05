package decorations

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetDecorationsTable returns bilingual paginated rows for decorations.
//
//encore:api public method=GET path=/decorations/table
func GetDecorationsTable(ctx context.Context, req *decorationsTableRequest) (*decorationsTableResponse, error) {
	svc, err := getDecorationsService()
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
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategoryDecorations, coreReq)
	if err != nil {
		return nil, err
	}
	return &decorationsTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
