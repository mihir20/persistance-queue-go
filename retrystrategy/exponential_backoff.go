package retrystrategy

import (
	"math"
	"time"
)

type ExponentialBackOffRetryStrategy struct {
	baseInterval time.Duration
	maxAttempts  int
}

func NewExponentialBackOffRetryStrategy(baseInterval time.Duration, maxAttempts int) *ExponentialBackOffRetryStrategy {
	return &ExponentialBackOffRetryStrategy{
		baseInterval: baseInterval,
		maxAttempts:  maxAttempts,
	}
}

// IsMaxRetryMet returns true if retry count is greater or equal to maxNumber of attempts
func (s *ExponentialBackOffRetryStrategy) IsMaxRetryMet(retryCount int) bool {
	return s.maxAttempts <= retryCount
}

// GetNextRetryTime returns eventPublishedTime + baseInterval*(2^nextAttemptNumber)
// e.g. event published time is 02:35:20 PM and nextAttemptNumber is 3 and base interval is 2 Sec
// Next retry time will be 02:35:20PM + 2sec*(2^3) = 02:35:36PM
func (s *ExponentialBackOffRetryStrategy) GetNextRetryTime(nextAttemptNumber int, eventPublishedTime time.Time) time.Time {
	return eventPublishedTime.Add(time.Duration(math.Pow(2, float64(nextAttemptNumber))) * s.baseInterval)
}
