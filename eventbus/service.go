package eventbus

import (
	"fmt"
	"persistent-queue/api/event"
	"persistent-queue/eventbus/dao"
)

type Service struct {
	eventsDao dao.EventsDao
}

func NewService(eventsDao dao.EventsDao) *Service {
	return &Service{
		eventsDao: eventsDao,
	}
}

func (s *Service) EnqueueEvent(event *event.Event) error {
	fmt.Printf("enqueuing new event in the bus, name: %s\n", event.Name)
	return nil
}
