package eventbus

import (
	"persistent-queue/api/event"
	"time"
)

type PassengerEvent struct {
	Event         *event.Event
	RetryAttempts int
	EventTime     time.Time
}
