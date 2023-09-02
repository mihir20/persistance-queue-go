package event

import "time"

type Event struct {
	Name        string `json:"name"`
	PublishedAt time.Time
}
