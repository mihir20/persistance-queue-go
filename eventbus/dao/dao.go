package dao

import "persistent-queue/api/event"

type EventsDao interface {
	CreateEvent(event event.Event) (event.Event, error)
	UpdateEvent(event event.Event) error
	DeleteEvent(event event.Event) error
}
