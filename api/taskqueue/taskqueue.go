package taskqueue

import "persistent-queue/api/event"

type TaskQueue string

type ITaskQueue interface {
	GetTaskQueueName() TaskQueue
	PollEventsQueue() error
	ConsumeEvent(event *event.Event) error
}
