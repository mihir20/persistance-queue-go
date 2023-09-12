//go:generate mockgen -source=service.go -destination=./mocks/mock.go
package eventbus

import (
	"persistent-queue/api/event"
	"persistent-queue/api/taskqueue"
	"time"
)

// IService (EventBusService) is responsible for storing and processing of all the events
type IService interface {
	// EnqueueEvent will add new event in the event queue
	EnqueueEvent(event *event.Event) error
	// GetEventsToProcess get events ready to be consumed by the consumer
	// taskQueue will be the consumer queue that is requesting the events
	// countOfElements is the number of events that are available for processing
	GetEventsToProcess(taskQueue taskqueue.TaskQueue, countOfElements int64) ([]*PassengerEvent, error)
	// DequeueEventFromTaskQueue will remove the event from the queue when event consumption comes to a terminal state(successfully consumed or permanently failed)
	DequeueEventFromTaskQueue(queue taskqueue.TaskQueue, passengerEvent *PassengerEvent) error
	// UpdatePassengerEvent updates event in the queue with updatedPassengerEvent and next execution time
	UpdatePassengerEvent(queue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *PassengerEvent, nextExecutionTime time.Time) error
	// CountEventsInQueue counts remaining events to be processed in a queue
	CountEventsInQueue() (map[taskqueue.TaskQueue]int64, error)
}
