package dao

import (
	"persistent-queue/api/event"
	"persistent-queue/api/taskqueue"
)

type EventsDao interface {
	CreateEvent(event *event.Event, taskQueues []taskqueue.TaskQueue) error
	GetEvent(taskQueue taskqueue.TaskQueue) (*event.Event, error)
	UpdateEvent(taskQueue taskqueue.TaskQueue, event *event.Event) error
	DeleteEvent(event *event.Event) error
}
