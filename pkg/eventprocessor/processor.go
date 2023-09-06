package eventprocessor

import "persistent-queue/api/eventbus"

type IEventProcessor interface {
	PollAndProcessEvents(consumerMethod func(event *eventbus.PassengerEvent) error) error
}
