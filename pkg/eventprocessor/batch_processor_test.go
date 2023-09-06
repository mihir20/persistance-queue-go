package eventprocessor

import (
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/pkg/retrystrategy"
	"testing"
)

func TestBatchEventProcessor_PollAndProcessEvents(t *testing.T) {
	type fields struct {
		batchSize       int
		taskQueue       taskqueue.TaskQueue
		eventBusService eventbus.IService
		retryStrategy   retrystrategy.IRetryStrategy
	}
	type args struct {
		consumerMethod func(event *eventbus.PassengerEvent) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &BatchEventProcessor{
				batchSize:       tt.fields.batchSize,
				taskQueue:       tt.fields.taskQueue,
				eventBusService: tt.fields.eventBusService,
				retryStrategy:   tt.fields.retryStrategy,
			}
			if err := p.PollAndProcessEvents(tt.args.consumerMethod); (err != nil) != tt.wantErr {
				t.Errorf("PollAndProcessEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
