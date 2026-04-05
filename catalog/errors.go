package catalog

import (
	"fmt"

	"encore.dev/beta/errs"
)

func invalidArgument(msg string) error {
	return errs.B().
		Code(errs.InvalidArgument).
		Msg(msg).
		Err()
}

func invalidArgumentf(format string, args ...any) error {
	return invalidArgument(fmt.Sprintf(format, args...))
}

func unsupportedLanguage(code string) error {
	return invalidArgumentf("unsupported language: %s", code)
}

func notFoundf(entity string, key string) error {
	return errs.B().
		Code(errs.NotFound).
		Msgf("%s not found for key %q", entity, key).
		Err()
}

func internalError(msg string, cause error) error {
	b := errs.B().
		Code(errs.Internal).
		Msg(msg)
	if cause != nil {
		b = b.Cause(cause)
	}
	return b.Err()
}
