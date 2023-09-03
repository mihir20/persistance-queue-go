package consumer

import (
	"log"
	eventModel "persistent-queue/api/event"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
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
		if err != nil {
			log.Printf("failed to consume event, err:%s\n", err.Error())
			return err
		}
	}
	return nil
}

func (s *Service) ConsumeEvent(event *eventModel.Event) error {
	log.Printf("consuming event %s\n", event.Name)
	return nil
}

func (s *Service) GetTaskQueueName() taskqueue.TaskQueue {
	return taskqueueNs.ConsumerTaskQueue
}
