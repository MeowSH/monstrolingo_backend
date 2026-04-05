package catalogcore

import "encore.dev/beta/errs"

func IsInvalidArgument(err error) bool {
	return errs.Code(err) == errs.InvalidArgument
}

func IsNotFound(err error) bool {
	return errs.Code(err) == errs.NotFound
}

func IsInternal(err error) bool {
	return errs.Code(err) == errs.Internal
}
