package catalog

func toDetailTranslation(row translationRecord) DetailTranslation {
	return DetailTranslation{
		LanguageCode:  row.LanguageCode,
		Name:          row.Name,
		Description:   row.Description,
		FlavorText:    row.FlavorText,
		EffectSummary: row.EffectSummary,
		Slug:          row.Slug,
	}
}

func toTableRow(row categoryTableDBRow, sourceLang, targetLang string) CategoryTableRow {
	return CategoryTableRow{
		ExternalKey: row.ExternalKey,
		Source: TableTranslation{
			Language:    sourceLang,
			Name:        row.SourceName,
			Description: row.SourceDescription,
		},
		Target: TableTranslation{
			Language:    targetLang,
			Name:        row.TargetName,
			Description: row.TargetDescription,
		},
	}
}

func buildPagination(page int, limit int, total int64) Pagination {
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
	}
}
