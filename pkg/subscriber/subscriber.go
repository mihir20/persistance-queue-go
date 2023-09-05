package subscriber

import (
	"github.com/go-co-op/gocron"
	"log"
	"math"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/pkg/errors"
	"persistent-queue/pkg/retrystrategy"
	"time"
)

const (
	minPollingFrequency = 1
)

type Subscriber struct {
	// time to poll the queue in seconds
	pollingFrequency int
	taskQueue        taskqueue.TaskQueue
	eventBusService  eventbus.IService
	retryStrategy    retrystrategy.IRetryStrategy
	consumeFunction  func(*eventbus.PassengerEvent) error
}

func NewSubscriber(pollingFrequency int, taskQueue taskqueue.TaskQueue, eventBusService eventbus.IService,
	retryStrategy retrystrategy.IRetryStrategy, consumeFunction func(*eventbus.PassengerEvent) error) *Subscriber {
	return &Subscriber{
		pollingFrequency: int(math.Max(float64(pollingFrequency), minPollingFrequency)),
		taskQueue:        taskQueue,
		eventBusService:  eventBusService,
		retryStrategy:    retryStrategy,
		consumeFunction:  consumeFunction,
	}
}

func (s *Subscriber) StartWorker() {
	scheduler := gocron.NewScheduler(time.UTC)

	// Every starts the job immediately and then runs at the
	// specified interval
	_, err := scheduler.Every(2).Seconds().Do(func() {
		err := s.pollAndConsumeEvents()
		if err != nil {
			log.Printf("error performing polling on snowflakeconsumer, err: %s\n", err.Error())
			return
		}
	})
	if err != nil {
		log.Printf("error setting up job, err %s\n", err.Error())
	}
	scheduler.StartBlocking()
}

func (s *Subscriber) pollAndConsumeEvents() error {
	event, err := s.eventBusService.GetEventToProcess(s.taskQueue)
	if err != nil {
		log.Printf("error while polling for event, err:%s\n", err.Error())
		return err
	}
	if event != nil {
		err = s.consumeEvent(event)
		s.processEventConsumption(err, event)
	} else {
		log.Printf("no event to process\n")
	}
	return nil
}

func (s *Subscriber) consumeEvent(event *eventbus.PassengerEvent) error {
	return s.consumeFunction(event)
}

func (s *Subscriber) processEventConsumption(err error, event *eventbus.PassengerEvent) {
	if err == nil || errors.IsPermanentError(err) {
		deleteErr := s.eventBusService.DequeueEventFromTaskQueue(s.taskQueue, event)
		if deleteErr != nil {
			log.Printf("failed to delete event from queue, err:%s", deleteErr.Error())
		}
		return
	}
	if errors.IsTransientError(err) {
		processErr := s.processTransientFailure(event)
		if processErr != nil {
			log.Printf("failed to process transient failure, err:%s", processErr.Error())
		}
		return
	}

	log.Fatalf("unknown error from snowflakeconsumer, err: %s", err.Error())
}

func (s *Subscriber) processTransientFailure(passengerEvent *eventbus.PassengerEvent) error {
	log.Printf("transient failure while processing event %s\n", passengerEvent.Event.UserID)
	oldPassenger := &eventbus.PassengerEvent{
		Event:         passengerEvent.Event,
		RetryAttempts: passengerEvent.RetryAttempts,
		EventTime:     passengerEvent.EventTime,
	}
	if s.retryStrategy.IsMaxRetryMet(passengerEvent.RetryAttempts + 1) {
		err := s.eventBusService.DequeueEventFromTaskQueue(s.taskQueue, passengerEvent)
		if err != nil {
			return err
		}
		return nil
	}
	passengerEvent.RetryAttempts++
	nextExecutionTime := s.retryStrategy.GetNextRetryTime(passengerEvent.RetryAttempts, passengerEvent.EventTime)
	err := s.eventBusService.UpdatePassengerEvent(s.taskQueue, oldPassenger, passengerEvent, nextExecutionTime)
	if err != nil {
		return err
	}
	return nil
}
