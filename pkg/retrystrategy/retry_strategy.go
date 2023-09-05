//go:generate mockgen -source=retry_strategy.go -destination=mocks/mock.go -package=mocks
package retrystrategy

import "time"

type IRetryStrategy interface {
	IsMaxRetryMet(retryCount int) bool
	GetNextRetryTime(nextAttemptNumber int, eventPublishedTime time.Time) time.Time
}
