package errors

import "errors"

var (
	ErrTransientFailure = errors.New("transient error")
	ErrPermanentFailure = errors.New("permanent failure")
)

func IsTransientError(err error) bool {
	return errors.Is(err, ErrTransientFailure)
}

func IsPermanentError(err error) bool {
	return errors.Is(err, ErrPermanentFailure)
}
