package eventbus

import (
	"persistent-queue/api/event"
	"persistent-queue/api/taskqueue"
	"time"
)

type IService interface {
	EnqueueEvent(event *event.Event) error
	GetEventToProcess(queue taskqueue.TaskQueue) (*PassengerEvent, error)
	DequeueEventFromTaskQueue(queue taskqueue.TaskQueue, passengerEvent *PassengerEvent) error
	UpdatePassengerEvent(queue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *PassengerEvent, nextExecutionTime time.Time) error
}
