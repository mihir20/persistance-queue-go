package retrystrategy

import (
	"reflect"
	"testing"
	"time"
)

func TestExponentialBackOffRetryStrategy_GetNextRetryTime(t *testing.T) {
	type fields struct {
		baseInterval time.Duration
		maxAttempts  int
	}
	type args struct {
		nextAttemptNumber  int
		eventPublishedTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ExponentialBackOffRetryStrategy{
				baseInterval: tt.fields.baseInterval,
				maxAttempts:  tt.fields.maxAttempts,
			}
			if got := s.GetNextRetryTime(tt.args.nextAttemptNumber, tt.args.eventPublishedTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNextRetryTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExponentialBackOffRetryStrategy_IsMaxRetryMet(t *testing.T) {
	type fields struct {
		baseInterval time.Duration
		maxAttempts  int
	}
	type args struct {
		retryCount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ExponentialBackOffRetryStrategy{
				baseInterval: tt.fields.baseInterval,
				maxAttempts:  tt.fields.maxAttempts,
			}
			if got := s.IsMaxRetryMet(tt.args.retryCount); got != tt.want {
				t.Errorf("IsMaxRetryMet() = %v, want %v", got, tt.want)
			}
		})
	}
}
