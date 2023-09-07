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
		{
			name: "1st retry",
			fields: fields{
				baseInterval: 2 * time.Second,
				maxAttempts:  3,
			},
			args: args{
				nextAttemptNumber:  1,
				eventPublishedTime: time.Unix(1000000, 0),
			},
			want: time.Unix(1000004, 0),
		},
		{
			name: "2nd retry",
			fields: fields{
				baseInterval: 2 * time.Second,
				maxAttempts:  3,
			},
			args: args{
				nextAttemptNumber:  2,
				eventPublishedTime: time.Unix(1000000, 0),
			},
			want: time.Unix(1000008, 0),
		},
		{
			name: "3rd retry",
			fields: fields{
				baseInterval: 2 * time.Second,
				maxAttempts:  3,
			},
			args: args{
				nextAttemptNumber:  3,
				eventPublishedTime: time.Unix(1000000, 0),
			},
			want: time.Unix(1000016, 0),
		},
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
		{
			name: "first retry out of 3",
			fields: fields{
				baseInterval: 2 * time.Second,
				maxAttempts:  3,
			},
			args: args{
				retryCount: 1,
			},
			want: false,
		},
		{
			name: "second retry out of 3",
			fields: fields{
				baseInterval: 2 * time.Second,
				maxAttempts:  3,
			},
			args: args{
				retryCount: 2,
			},
			want: false,
		},
		{
			name: "second retry out of 3",
			fields: fields{
				baseInterval: 2 * time.Second,
				maxAttempts:  3,
			},
			args: args{
				retryCount: 3,
			},
			want: true,
		},
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
