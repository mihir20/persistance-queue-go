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
	eventbus.IService
	eventsDao            dao.EventsDao
	registeredTaskQueues []taskqueue.TaskQueue
}

func NewService(eventsDao dao.EventsDao) *Service {
	return &Service{
		eventsDao: eventsDao,
		registeredTaskQueues: []taskqueue.TaskQueue{
			taskqueueNs.SnowflakeConsumerTaskQueue,
			taskqueueNs.VendorApiConsumerTaskQueue,
			taskqueueNs.FileConsumerTaskQueue,
		},
	}
}

func (s *Service) EnqueueEvent(event *event.Event) error {
	fmt.Printf("enqueuing new event in the bus, userid: %s\n", event.UserID)
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
	savedEvent, executionTimestamp, err := s.eventsDao.GetEvent(taskQueue)
	if err != nil {
		log.Printf("error getting event, err: %s\n", err.Error())
		return nil, err
	}
	if time.Now().Before(time.Unix(executionTimestamp, 0).UTC()) {
		return nil, nil
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

func (s *Service) UpdatePassengerEvent(taskQueue taskqueue.TaskQueue, oldPassenger, newPassengerEvent *eventbus.PassengerEvent, nextExecutionTime time.Time) error {
	err := s.eventsDao.UpdateEvent(taskQueue, oldPassenger, newPassengerEvent, nextExecutionTime)
	if err != nil {
		log.Printf("error deleting event, err: %s\n", err.Error())
		return err
	}
	return nil
}

func (s *Service) CountEventsInQueue() (map[taskqueue.TaskQueue]int64, error) {
	mp, err := s.eventsDao.CountEventsInQueue(s.registeredTaskQueues)
	if err != nil {
		log.Printf("error deleting event, err: %s\n", err.Error())
		return nil, err
	}
	return mp, nil
}
