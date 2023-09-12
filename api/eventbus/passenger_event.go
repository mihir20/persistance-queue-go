package eventbus

import (
	"persistent-queue/api/event"
	"time"
)

// PassengerEvent is data structure for the event received by EventBusService
type PassengerEvent struct {
	//
	Event         *event.Event
	RetryAttempts int
	EventTime     time.Time
}
