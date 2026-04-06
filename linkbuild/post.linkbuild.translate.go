package linkbuild

import (
	"context"

	"monstrolingo_backend/internal/simbuildcore"
)

// PostLinkbuildTranslate translates entries from a supported sim link.
//
//encore:api public method=POST path=/linkbuild/translate
func PostLinkbuildTranslate(ctx context.Context, req *linkbuildTranslateRequest) (*linkbuildTranslateResponse, error) {
	svc, err := getLinkbuildService()
	if err != nil {
		return nil, err
	}

	coreReq := &simbuildcore.TranslateRequest{}
	if req != nil {
		coreReq.URL = req.URL
		coreReq.TargetLang = req.TargetLang
	}

	out, err := svc.TranslateSimBuild(ctx, coreReq)
	if err != nil {
		return nil, err
	}

	return &linkbuildTranslateResponse{
		SourceLangDetected: out.SourceLangDetected,
		TargetLang:         out.TargetLang,
		TranslationMode:    out.TranslationMode,
		SkillsOriginal:     out.SkillsOriginal,
		SkillsTranslated:   out.SkillsTranslated,
		UnmatchedElements:  out.UnmatchedElements,
	}, nil
}
