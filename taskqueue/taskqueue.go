package taskqueue

import "persistent-queue/api/event"

type ITaskQueue interface {
	PollEventsQueue() error
	ConsumeEvent(event *event.Event) error
}
