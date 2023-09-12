//go:generate mockgen -source=consumer.go -destination=mocks/mock.go -package=mocks
package consumer

import "persistent-queue/api/eventbus"

// IConsumer knows how to consume an event
type IConsumer interface {
	Consume(event *eventbus.PassengerEvent) error
}

type Consumer struct {
	consume func(*eventbus.PassengerEvent) error
}

func NewConsumer(consume func(*eventbus.PassengerEvent) error) *Consumer {
	return &Consumer{consume: consume}
}

func (c *Consumer) Consume(event *eventbus.PassengerEvent) error {
	return c.consume(event)
}
