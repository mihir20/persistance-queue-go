package dao

import (
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
)

type EventsDao interface {
	CreateEvent(event *eventbus.PassengerEvent, taskQueues []taskqueue.TaskQueue) error
	GetEvent(taskQueue taskqueue.TaskQueue) (*eventbus.PassengerEvent, error)
	UpdateEvent(taskQueue taskqueue.TaskQueue, updatedPassengerEvent *eventbus.PassengerEvent) error
	DeleteEvent(event *eventbus.PassengerEvent) error
}
