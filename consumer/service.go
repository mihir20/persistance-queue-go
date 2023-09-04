package consumer

import (
	"log"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/pkg/errors"
	taskqueueNs "persistent-queue/pkg/taskqueue"
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
	}
	return nil
}

func (s *Service) ConsumeEvent(event *eventbus.PassengerEvent) error {
	log.Printf("consuming event %s\n", event.Event.Name)
	return nil
}

func (s *Service) GetTaskQueueName() taskqueue.TaskQueue {
	return taskqueueNs.ConsumerTaskQueue
}

func (s *Service) processEventConsumption(err error, event *eventbus.PassengerEvent) {
	if err == nil || errors.IsPermanentError(err) {
		// TODO: delete event
		return
	}
	if errors.IsTransientError(err) {
		// TODO: retry event
		return
	}

	log.Fatalf("unknown error from consumer, err: %s", err.Error())
}
