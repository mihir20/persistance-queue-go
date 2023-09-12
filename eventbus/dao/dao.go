//go:generate mockgen -source=dao.go -destination=mocks/mock_dao.go -package=mocks
package dao

import (
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"time"
)

// EventsDao is Data access layer for events
type EventsDao interface {
	// CreateEvent creates a new db entry of a event in the taskQueue
	CreateEvent(event *eventbus.PassengerEvent, taskQueues []taskqueue.TaskQueue) error
	// GetEvents returns countOfEvents number of events to be processed till cutOffTime
	GetEvents(taskQueue taskqueue.TaskQueue, cutOffTime time.Time, countOfEvents int64) ([]*eventbus.PassengerEvent, error)
	// UpdateEvent updates the event in the db
	UpdateEvent(taskQueue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *eventbus.PassengerEvent, nextExecutionTime time.Time) error
	// DeleteEvent deletes the event from the db
	DeleteEvent(taskQueue taskqueue.TaskQueue, event *eventbus.PassengerEvent) error
	// CountEventsInQueue returns count of events in given task queues
	CountEventsInQueue(taskQueues []taskqueue.TaskQueue) (map[taskqueue.TaskQueue]int64, error)
}
