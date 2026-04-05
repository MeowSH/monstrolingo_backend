package catalogcore

func ToTargetLanguageRequest(targetLang string) *TargetLanguageRequest {
	return &TargetLanguageRequest{TargetLang: targetLang}
}
