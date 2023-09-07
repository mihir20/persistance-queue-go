package eventprocessor

type IEventProcessor interface {
	PollAndProcessEvents() error
}
