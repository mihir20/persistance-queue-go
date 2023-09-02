package eventbus

import (
	"fmt"
	"persistent-queue/api/event"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) EnqueueEvent(event *event.Event) error {
	fmt.Printf("enqueuing new event in the bus, name: %s\n", event.Name)
	return nil
}
