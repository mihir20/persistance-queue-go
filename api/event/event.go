package event

import "time"

// Event is data structure for the event consumed by the queue
type Event struct {
	// identifier for the event
	EventId string
	// user identifier for event
	UserID string
	// payload for event
	Payload string
	// time at which event was published
	PublishedAt time.Time
}
