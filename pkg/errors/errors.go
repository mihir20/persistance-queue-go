package errors

import "errors"

var (
	// ErrTransientFailure error for transient failure
	ErrTransientFailure = errors.New("transient error")
	// ErrPermanentFailure error for permanent failure
	ErrPermanentFailure = errors.New("permanent failure")
)

// IsTransientError checks if err is for transient failure or not
func IsTransientError(err error) bool {
	return errors.Is(err, ErrTransientFailure)
}

// IsPermanentError checks if err is for permanent failure or not
func IsPermanentError(err error) bool {
	return errors.Is(err, ErrPermanentFailure)
}
