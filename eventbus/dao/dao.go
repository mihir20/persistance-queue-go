package dao

import "persistent-queue/api/event"

type EventsDao interface {
	CreateEvent(event *event.Event) error
	GetEvent(eventName string) (*event.Event, error)
	UpdateEvent(event *event.Event) error
	DeleteEvent(event *event.Event) error
}
