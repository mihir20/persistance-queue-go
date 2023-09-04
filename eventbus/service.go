package eventbus

import (
	"fmt"
	"log"
	"persistent-queue/api/event"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"persistent-queue/eventbus/dao"
	taskqueueNs "persistent-queue/pkg/taskqueue"
	"time"
)

type Service struct {
	eventsDao            dao.EventsDao
	registeredTaskQueues []taskqueue.TaskQueue
}

func NewService(eventsDao dao.EventsDao) *Service {
	return &Service{
		eventsDao: eventsDao,
		registeredTaskQueues: []taskqueue.TaskQueue{
			taskqueueNs.ConsumerTaskQueue,
		},
	}
}

func (s *Service) EnqueueEvent(event *event.Event) error {
	fmt.Printf("enqueuing new event in the bus, name: %s\n", event.Name)
	err := s.eventsDao.CreateEvent(&eventbus.PassengerEvent{
		Event:         event,
		RetryAttempts: 0,
		EventTime:     time.Now(),
	}, s.registeredTaskQueues)
	if err != nil {
		log.Printf("error saving event, err: %s\n", err.Error())
	}
	return err
}

func (s *Service) GetEventToProcess(taskQueue taskqueue.TaskQueue) (*eventbus.PassengerEvent, error) {
	savedEvent, err := s.eventsDao.GetEvent(taskQueue)
	if err != nil {
		log.Printf("error getting event, err: %s\n", err.Error())
		return nil, err
	}
	return savedEvent, nil
}

func (s *Service) DequeueEventFromTaskQueue(taskQueue taskqueue.TaskQueue, passengerEvent *eventbus.PassengerEvent) error {
	err := s.eventsDao.DeleteEvent(taskQueue, passengerEvent)
	if err != nil {
		log.Printf("error deleting event, err: %s\n", err.Error())
		return err
	}
	return nil
}
