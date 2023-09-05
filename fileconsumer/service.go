package fileconsumer

import (
	"log"
	"persistent-queue/api/eventbus"
	"persistent-queue/pkg/errors"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ConsumeEvent(event *eventbus.PassengerEvent) error {
	log.Printf("consuming event %s\n", event.Event.UserID)
	return errors.ErrPermanentFailure
}
