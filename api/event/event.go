package event

import "time"

type Event struct {
	EventId     string
	UserID      string
	Payload     string
	PublishedAt time.Time
}
