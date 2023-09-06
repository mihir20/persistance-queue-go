package eventprocessor

import (
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/pkg/retrystrategy"
)

type BatchEventProcessor struct {
	batchSize       int
	taskQueue       taskqueue.TaskQueue
	eventBusService eventbus.IService
	retryStrategy   retrystrategy.IRetryStrategy
}

func NewBatchEventProcessor(batchSize int, taskQueue taskqueue.TaskQueue,
	eventBusService eventbus.IService, retryStrategy retrystrategy.IRetryStrategy) *BatchEventProcessor {
	return &BatchEventProcessor{
		batchSize:       batchSize,
		taskQueue:       taskQueue,
		eventBusService: eventBusService,
		retryStrategy:   retryStrategy,
	}
}

func (p *BatchEventProcessor) PollAndProcessEvents(consumerMethod func(event *eventbus.PassengerEvent) error) error {
	return nil
}
