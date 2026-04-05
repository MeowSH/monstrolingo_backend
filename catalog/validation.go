package catalog

import (
	"regexp"
	"strings"
)

const (
	defaultPage  = 1
	defaultLimit = 25
	maxLimit     = 100
)

var languageCodePattern = regexp.MustCompile(`^[a-z]{2}(?:-[a-z]{2})?$`)

func normalizeTableQuery(req *CategoryTableRequest) (normalizedTableQuery, error) {
	if req == nil {
		return normalizedTableQuery{}, invalidArgument("missing request payload")
	}

	sourceLang, err := normalizeLanguageCode(req.SourceLang, "source_lang")
	if err != nil {
		return normalizedTableQuery{}, err
	}
	targetLang, err := normalizeLanguageCode(req.TargetLang, "target_lang")
	if err != nil {
		return normalizedTableQuery{}, err
	}

	page := req.Page
	if page <= 0 {
		page = defaultPage
	}

	limit := req.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		return normalizedTableQuery{}, invalidArgumentf("limit must be <= %d", maxLimit)
	}

	return normalizedTableQuery{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		Page:       page,
		Limit:      limit,
		Offset:     (page - 1) * limit,
	}, nil
}

func normalizeDetailQuery(req *CategoryDetailRequest) (normalizedDetailQuery, error) {
	if req == nil {
		return normalizedDetailQuery{}, invalidArgument("missing request payload")
	}
	externalKey := strings.TrimSpace(req.ExternalKey)
	if externalKey == "" {
		return normalizedDetailQuery{}, invalidArgument("external_key is required")
	}
	targetLang, err := normalizeLanguageCode(req.TargetLang, "target_lang")
	if err != nil {
		return normalizedDetailQuery{}, err
	}
	return normalizedDetailQuery{
		ExternalKey: externalKey,
		TargetLang:  targetLang,
	}, nil
}

func normalizeLanguageCode(raw string, field string) (string, error) {
	code := strings.ToLower(strings.TrimSpace(raw))
	if code == "" {
		return "", invalidArgumentf("%s is required", field)
	}
	if !languageCodePattern.MatchString(code) {
		return "", invalidArgumentf("%s has invalid format", field)
	}
	return code, nil
}
