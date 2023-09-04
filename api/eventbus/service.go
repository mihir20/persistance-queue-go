package eventbus

import (
	"persistent-queue/api/event"
	"persistent-queue/api/taskqueue"
)

type IService interface {
	EnqueueEvent(event *event.Event) error
	GetEventToProcess(queue taskqueue.TaskQueue) (*PassengerEvent, error)
}
