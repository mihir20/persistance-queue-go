package eventbus

import (
	"persistent-queue/api/event"
)

type IService interface {
	EnqueueEvent(event *event.Event) error
}
