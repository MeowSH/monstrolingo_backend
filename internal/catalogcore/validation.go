package catalogcore

const (
	DefaultPage  = 1
	DefaultLimit = 25
	MaxLimit     = 100
)

func NewTableRequest(sourceLang, targetLang string, page, limit int) *CategoryTableRequest {
	return &CategoryTableRequest{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		Page:       page,
		Limit:      limit,
	}
}

func NewDetailRequest(externalKey, targetLang string) *CategoryDetailRequest {
	return &CategoryDetailRequest{
		ExternalKey: externalKey,
		TargetLang:  targetLang,
	}
}
