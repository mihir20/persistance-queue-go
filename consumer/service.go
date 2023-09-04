package consumer

import (
	"log"
	"math"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/pkg/errors"
	taskqueueNs "persistent-queue/pkg/taskqueue"
	"time"
)

type Service struct {
	eventBusService eventbus.IService
}

func NewService(eventBusService eventbus.IService) *Service {
	return &Service{
		eventBusService: eventBusService,
	}
}

func (s *Service) PollEventsQueue() error {
	event, err := s.eventBusService.GetEventToProcess(s.GetTaskQueueName())
	if err != nil {
		log.Printf("error while polling for event, err:%s\n", err.Error())
		return err
	}
	if event != nil {
		err = s.ConsumeEvent(event)
		s.processEventConsumption(err, event)
	} else {
		log.Printf("no event to process\n")
	}
	return nil
}

func (s *Service) ConsumeEvent(event *eventbus.PassengerEvent) error {
	log.Printf("consuming event %s\n", event.Event.Name)
	return errors.ErrTransientFailure
}

func (s *Service) GetTaskQueueName() taskqueue.TaskQueue {
	return taskqueueNs.ConsumerTaskQueue
}

func (s *Service) processEventConsumption(err error, event *eventbus.PassengerEvent) {
	if err == nil || errors.IsPermanentError(err) {
		deleteErr := s.eventBusService.DequeueEventFromTaskQueue(s.GetTaskQueueName(), event)
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

	log.Fatalf("unknown error from consumer, err: %s", err.Error())
}

func (s *Service) processTransientFailure(passengerEvent *eventbus.PassengerEvent) error {
	log.Printf("transient failure while processing event %s\n", passengerEvent.Event.Name)
	maxAttempts := 3
	baseInterval := 6 * time.Second
	oldPassenger := &eventbus.PassengerEvent{
		Event:         passengerEvent.Event,
		RetryAttempts: passengerEvent.RetryAttempts,
		EventTime:     passengerEvent.EventTime,
	}
	if passengerEvent.RetryAttempts+1 == maxAttempts {
		err := s.eventBusService.DequeueEventFromTaskQueue(s.GetTaskQueueName(), passengerEvent)
		if err != nil {
			return err
		}
		return nil
	}
	passengerEvent.RetryAttempts++
	nextExecutionTime := passengerEvent.EventTime.Add(time.Duration(math.Pow(2, float64(passengerEvent.RetryAttempts))) * baseInterval)
	err := s.eventBusService.UpdatePassengerEvent(s.GetTaskQueueName(), oldPassenger, passengerEvent, nextExecutionTime)
	if err != nil {
		return err
	}
	return nil
}
