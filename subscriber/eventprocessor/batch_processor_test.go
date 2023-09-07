package eventprocessor

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"persistent-queue/api/event"
	"persistent-queue/api/eventbus"
	mock_eventbus "persistent-queue/api/eventbus/mocks"
	"persistent-queue/api/taskqueue"
	"persistent-queue/pkg/errors"
	taskqueueNs "persistent-queue/pkg/taskqueue"
	"persistent-queue/retrystrategy/mocks"
	mocks2 "persistent-queue/subscriber/consumer/mocks"
	"testing"
	"time"
)

var (
	sampleTaskQueue = taskqueueNs.SnowflakeConsumerTaskQueue

	sampleEventPassenger1 = &eventbus.PassengerEvent{
		Event: &event.Event{
			EventId:     uuid.NewString(),
			UserID:      uuid.NewString(),
			Payload:     uuid.NewString(),
			PublishedAt: time.Unix(16000000, 0),
		},
		RetryAttempts: 0,
		EventTime:     time.Unix(16000000, 0),
	}
	sampleEventPassenger2 = &eventbus.PassengerEvent{
		Event: &event.Event{
			EventId:     uuid.NewString(),
			UserID:      uuid.NewString(),
			Payload:     uuid.NewString(),
			PublishedAt: time.Unix(16000000, 0),
		},
		RetryAttempts: 1,
		EventTime:     time.Unix(16000000, 0),
	}
)

func TestBatchEventProcessor_PollAndProcessEvents(t *testing.T) {
	type fields struct {
		batchSize int
		taskQueue taskqueue.TaskQueue
	}
	tests := []struct {
		name      string
		fields    fields
		mockSetup func(retryStrategy *mocks.MockIRetryStrategy, eventBusService *mock_eventbus.MockIService,
			consumer *mocks2.MockIConsumer)
		wantErr bool
	}{
		{
			name: "successful event process batch size 1",
			fields: fields{
				batchSize: 1,
				taskQueue: sampleTaskQueue,
			},
			mockSetup: func(retryStrategy *mocks.MockIRetryStrategy, eventBusService *mock_eventbus.MockIService, consumer *mocks2.MockIConsumer) {
				eventBusService.EXPECT().
					GetEventsToProcess(sampleTaskQueue, int64(1)).
					Return([]*eventbus.PassengerEvent{sampleEventPassenger1}, nil)
				consumer.EXPECT().Consume(sampleEventPassenger1).Return(nil)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger1).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful event process batch size 2",
			fields: fields{
				batchSize: 2,
				taskQueue: sampleTaskQueue,
			},
			mockSetup: func(retryStrategy *mocks.MockIRetryStrategy, eventBusService *mock_eventbus.MockIService, consumer *mocks2.MockIConsumer) {
				eventBusService.EXPECT().
					GetEventsToProcess(sampleTaskQueue, int64(2)).
					Return([]*eventbus.PassengerEvent{sampleEventPassenger1, sampleEventPassenger2}, nil)
				consumer.EXPECT().Consume(sampleEventPassenger1).Return(nil)
				consumer.EXPECT().Consume(sampleEventPassenger2).Return(nil)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger1).
					Return(nil)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger2).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "permanent fail event process batch size 1",
			fields: fields{
				batchSize: 1,
				taskQueue: sampleTaskQueue,
			},
			mockSetup: func(retryStrategy *mocks.MockIRetryStrategy, eventBusService *mock_eventbus.MockIService, consumer *mocks2.MockIConsumer) {
				eventBusService.EXPECT().
					GetEventsToProcess(sampleTaskQueue, int64(1)).
					Return([]*eventbus.PassengerEvent{sampleEventPassenger1}, nil)
				consumer.EXPECT().Consume(sampleEventPassenger1).Return(errors.ErrPermanentFailure)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger1).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "1 permanent 1 success event process batch size 2",
			fields: fields{
				batchSize: 2,
				taskQueue: sampleTaskQueue,
			},
			mockSetup: func(retryStrategy *mocks.MockIRetryStrategy, eventBusService *mock_eventbus.MockIService, consumer *mocks2.MockIConsumer) {
				eventBusService.EXPECT().
					GetEventsToProcess(sampleTaskQueue, int64(2)).
					Return([]*eventbus.PassengerEvent{sampleEventPassenger1, sampleEventPassenger2}, nil)
				consumer.EXPECT().Consume(sampleEventPassenger1).Return(nil)
				consumer.EXPECT().Consume(sampleEventPassenger2).Return(errors.ErrPermanentFailure)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger1).
					Return(nil)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger2).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "1 permanent 1 transient(retry exhaust) event process batch size 2",
			fields: fields{
				batchSize: 2,
				taskQueue: sampleTaskQueue,
			},
			mockSetup: func(retryStrategy *mocks.MockIRetryStrategy, eventBusService *mock_eventbus.MockIService, consumer *mocks2.MockIConsumer) {
				eventBusService.EXPECT().
					GetEventsToProcess(sampleTaskQueue, int64(2)).
					Return([]*eventbus.PassengerEvent{sampleEventPassenger1, sampleEventPassenger2}, nil)
				consumer.EXPECT().Consume(sampleEventPassenger1).Return(nil)
				consumer.EXPECT().Consume(sampleEventPassenger2).Return(errors.ErrTransientFailure)
				retryStrategy.EXPECT().IsMaxRetryMet(2).Return(true)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger1).
					Return(nil)
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger2).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "1 permanent 1 transient(retry not exhaust) event process batch size 2",
			fields: fields{
				batchSize: 2,
				taskQueue: sampleTaskQueue,
			},
			mockSetup: func(retryStrategy *mocks.MockIRetryStrategy, eventBusService *mock_eventbus.MockIService, consumer *mocks2.MockIConsumer) {
				eventBusService.EXPECT().
					GetEventsToProcess(sampleTaskQueue, int64(2)).
					Return([]*eventbus.PassengerEvent{sampleEventPassenger1, sampleEventPassenger2}, nil)
				consumer.EXPECT().Consume(sampleEventPassenger1).Return(nil)
				consumer.EXPECT().Consume(sampleEventPassenger2).Return(errors.ErrTransientFailure)
				retryStrategy.EXPECT().IsMaxRetryMet(2).Return(false)
				nextExecutionTime := time.Unix(18000000, 0)
				retryStrategy.EXPECT().GetNextRetryTime(2, sampleEventPassenger2.EventTime).Return(nextExecutionTime)
				oldPassenger := &eventbus.PassengerEvent{
					Event:         sampleEventPassenger2.Event,
					RetryAttempts: sampleEventPassenger2.RetryAttempts,
					EventTime:     sampleEventPassenger2.EventTime,
				}
				newPassenger := &eventbus.PassengerEvent{
					Event:         sampleEventPassenger2.Event,
					RetryAttempts: sampleEventPassenger2.RetryAttempts + 1,
					EventTime:     sampleEventPassenger2.EventTime,
				}
				eventBusService.EXPECT().DequeueEventFromTaskQueue(sampleTaskQueue, sampleEventPassenger1).
					Return(nil)
				eventBusService.EXPECT().UpdatePassengerEvent(sampleTaskQueue, oldPassenger, newPassenger, nextExecutionTime).
					Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockEventBus := mock_eventbus.NewMockIService(ctrl)
			mockRetryStrategy := mocks.NewMockIRetryStrategy(ctrl)
			mockConsumer := mocks2.NewMockIConsumer(ctrl)
			p := &BatchEventProcessor{
				batchSize:       tt.fields.batchSize,
				taskQueue:       tt.fields.taskQueue,
				eventBusService: mockEventBus,
				retryStrategy:   mockRetryStrategy,
				consumer:        mockConsumer,
			}
			tt.mockSetup(mockRetryStrategy, mockEventBus, mockConsumer)
			if err := p.PollAndProcessEvents(); (err != nil) != tt.wantErr {
				t.Errorf("PollAndProcessEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
