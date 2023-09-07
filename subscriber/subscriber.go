package subscriber

import (
	"github.com/go-co-op/gocron"
	"log"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/retrystrategy"
	eventprocessor2 "persistent-queue/subscriber/eventprocessor"
	"time"
)

type Subscriber struct {
	// time to poll the queue in seconds
	pollingFrequency int
	consumeFunction  func(*eventbus.PassengerEvent) error
	eventProcessor   eventprocessor2.IEventProcessor
}

func NewSubscriber(numberOfWorkers, pollingFrequency int, taskQueue taskqueue.TaskQueue, eventBusService eventbus.IService,
	retryStrategy retrystrategy.IRetryStrategy, consumeFunction func(*eventbus.PassengerEvent) error) *Subscriber {
	return &Subscriber{
		pollingFrequency: pollingFrequency,
		consumeFunction:  consumeFunction,
		eventProcessor:   eventprocessor2.NewBatchEventProcessor(numberOfWorkers, taskQueue, eventBusService, retryStrategy),
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
	return s.eventProcessor.PollAndProcessEvents(s.consumeFunction)
}
