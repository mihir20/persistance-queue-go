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

func (s *ExponentialBackOffRetryStrategy) IsMaxRetryMet(retryCount int) bool {
	return s.maxAttempts <= retryCount
}

func (s *ExponentialBackOffRetryStrategy) GetNextRetryTime(nextAttemptNumber int, eventPublishedTime time.Time) time.Time {
	return eventPublishedTime.Add(time.Duration(math.Pow(2, float64(nextAttemptNumber))) * s.baseInterval)
}
