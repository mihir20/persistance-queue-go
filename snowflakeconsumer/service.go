package snowflakeconsumer

import (
	"log"
	"persistent-queue/api/eventbus"
)

type Service struct {
	eventBusService eventbus.IService
}

func NewService(eventBusService eventbus.IService) *Service {
	return &Service{
		eventBusService: eventBusService,
	}
}

func (s *Service) ConsumeEvent(event *eventbus.PassengerEvent) error {
	log.Printf("consuming event %s\n", event.Event.UserID)
	return nil
}
