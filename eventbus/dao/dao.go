//go:generate mockgen -source=dao.go -destination=mocks/mock_dao.go -package=mocks
package dao

import (
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"time"
)

type EventsDao interface {
	CreateEvent(event *eventbus.PassengerEvent, taskQueues []taskqueue.TaskQueue) error
	GetEvents(taskQueue taskqueue.TaskQueue, cutOffTime time.Time, countOfEvents int64) ([]*eventbus.PassengerEvent, error)
	UpdateEvent(taskQueue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *eventbus.PassengerEvent, nextExecutionTime time.Time) error
	DeleteEvent(taskQueue taskqueue.TaskQueue, event *eventbus.PassengerEvent) error
	CountEventsInQueue(taskQueues []taskqueue.TaskQueue) (map[taskqueue.TaskQueue]int64, error)
}
