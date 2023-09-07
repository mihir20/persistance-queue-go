package eventprocessor

import (
	"log"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/pkg/errors"
	"persistent-queue/pkg/goroutine"
	"persistent-queue/retrystrategy"
	"persistent-queue/subscriber/consumer"
	"sync"
)

type BatchEventProcessor struct {
	batchSize       int
	taskQueue       taskqueue.TaskQueue
	eventBusService eventbus.IService
	retryStrategy   retrystrategy.IRetryStrategy
	consumer        consumer.IConsumer
}

func NewBatchEventProcessor(batchSize int, taskQueue taskqueue.TaskQueue,
	eventBusService eventbus.IService, retryStrategy retrystrategy.IRetryStrategy,
	consumer consumer.IConsumer) *BatchEventProcessor {
	return &BatchEventProcessor{
		batchSize:       batchSize,
		taskQueue:       taskQueue,
		eventBusService: eventBusService,
		retryStrategy:   retryStrategy,
		consumer:        consumer,
	}
}

func (p *BatchEventProcessor) PollAndProcessEvents() error {
	events, err := p.eventBusService.GetEventsToProcess(p.taskQueue, int64(p.batchSize))
	if err != nil {
		log.Printf("error while polling for event, err:%s\n", err.Error())
		return err
	}
	if events != nil && len(events) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(events))
		for _, event := range events {
			localEvent := event
			goroutine.Run(func() {
				p.processEvent(localEvent, p.consumer.Consume)
				wg.Done()
			})
		}
		wg.Wait()
	} else {
		log.Printf("no event to process\n")
	}
	return nil
}

func (p *BatchEventProcessor) processEvent(event *eventbus.PassengerEvent, consumerMethod func(event *eventbus.PassengerEvent) error) {
	err := consumerMethod(event)
	p.processEventConsumptionResponse(err, event)
}

func (p *BatchEventProcessor) processEventConsumptionResponse(err error, event *eventbus.PassengerEvent) {
	if err == nil || errors.IsPermanentError(err) {
		deleteErr := p.eventBusService.DequeueEventFromTaskQueue(p.taskQueue, event)
		if deleteErr != nil {
			log.Printf("failed to delete event from queue, err:%s", deleteErr.Error())
		}
		return
	}
	if errors.IsTransientError(err) {
		processErr := p.processTransientFailure(event)
		if processErr != nil {
			log.Printf("failed to process transient failure, err:%s", processErr.Error())
		}
		return
	}

	log.Fatalf("unknown error from snowflakeconsumer, err: %s", err.Error())
}

func (p *BatchEventProcessor) processTransientFailure(passengerEvent *eventbus.PassengerEvent) error {
	log.Printf("transient failure while processing event %s\n", passengerEvent.Event.UserID)
	oldPassenger := &eventbus.PassengerEvent{
		Event:         passengerEvent.Event,
		RetryAttempts: passengerEvent.RetryAttempts,
		EventTime:     passengerEvent.EventTime,
	}
	if p.retryStrategy.IsMaxRetryMet(passengerEvent.RetryAttempts + 1) {
		err := p.eventBusService.DequeueEventFromTaskQueue(p.taskQueue, passengerEvent)
		if err != nil {
			return err
		}
		return nil
	}
	passengerEvent.RetryAttempts++
	nextExecutionTime := p.retryStrategy.GetNextRetryTime(passengerEvent.RetryAttempts, passengerEvent.EventTime)
	err := p.eventBusService.UpdatePassengerEvent(p.taskQueue, oldPassenger, passengerEvent, nextExecutionTime)
	if err != nil {
		return err
	}
	return nil
}
