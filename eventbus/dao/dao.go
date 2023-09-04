package dao

import (
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"time"
)

type EventsDao interface {
	CreateEvent(event *eventbus.PassengerEvent, taskQueues []taskqueue.TaskQueue) error
	GetEvent(taskQueue taskqueue.TaskQueue) (*eventbus.PassengerEvent, int64, error)
	UpdateEvent(taskQueue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *eventbus.PassengerEvent, nextExecutionTime time.Time) error
	DeleteEvent(taskQueue taskqueue.TaskQueue, event *eventbus.PassengerEvent) error
}
