//go:generate mockgen -source=retry_strategy.go -destination=mocks/mock.go -package=mocks
package retrystrategy

import "time"

type IRetryStrategy interface {
	// IsMaxRetryMet checks if maximum retries are met nor not
	IsMaxRetryMet(retryCount int) bool
	// GetNextRetryTime returns next retry time for an event
	// nextAttemptNumber is next attempt number of a retry, i.e. if an event will be processed 2nd time on next retry then nextAttemptNumber = 2
	GetNextRetryTime(nextAttemptNumber int, eventPublishedTime time.Time) time.Time
}
