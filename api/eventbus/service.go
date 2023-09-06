//go:generate mockgen -source=service.go -destination=./mocks/mock.go
package eventbus

import (
	"persistent-queue/api/event"
	"persistent-queue/api/taskqueue"
	"time"
)

type IService interface {
	EnqueueEvent(event *event.Event) error
	GetEventsToProcess(taskQueue taskqueue.TaskQueue, countOfElements int64) ([]*PassengerEvent, error)
	DequeueEventFromTaskQueue(queue taskqueue.TaskQueue, passengerEvent *PassengerEvent) error
	UpdatePassengerEvent(queue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *PassengerEvent, nextExecutionTime time.Time) error
	CountEventsInQueue() (map[taskqueue.TaskQueue]int64, error)
}
