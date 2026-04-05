package items

import (
	"context"

	"monstrolingo_backend/internal/catalogcore"
)

// GetItemsTable returns bilingual paginated rows for items.
//
//encore:api public method=GET path=/items/table
func GetItemsTable(ctx context.Context, req *itemsTableRequest) (*itemsTableResponse, error) {
	svc, err := getItemsService()
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
	out, err := svc.ListCategoryTable(ctx, catalogcore.CategoryItems, coreReq)
	if err != nil {
		return nil, err
	}
	return &itemsTableResponse{
		Items:      out.Items,
		Pagination: out.Pagination,
	}, nil
}
