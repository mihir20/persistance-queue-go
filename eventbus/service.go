package eventbus

import (
	"fmt"
	"log"
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
	err := s.eventsDao.CreateEvent(event)
	if err != nil {
		log.Printf("error saving event, err: %s\n", err.Error())
	}
	savedEvent, err := s.eventsDao.GetEvent(event.Name)
	if err != nil {
		log.Printf("error saving event, err: %s\n", err.Error())
	}
	fmt.Printf("event saved in the bus, name: %s\n", savedEvent.Name)
	return err
}
