package event

type Event interface {
	Name() string
	IsAsynchronous() bool
}

type GeneralError string

func (e GeneralError) Name() string {
	return "event.general.error"
}

type EventHandler interface {
	Notify(event Event)
}

type EventPublisher struct {
	handlers map[string][]EventHandler
}

func NewEventPublisher() *EventPublisher {
	p := EventPublisher{
		handlers: make(map[string][]EventHandler),
	}

	return &p
}

func (e *EventPublisher) Subscribe(handler EventHandler, events ...Event) {
	for _, event := range events {
		handlers := e.handlers[event.Name()]
		handlers = append(handlers, handler)
		e.handlers[event.Name()] = handlers
	}
}

func (e *EventPublisher) Notify(event Event) {
	if event.IsAsynchronous() {
		go e.notify(event)
	}

	e.notify(event)
}

func (e *EventPublisher) notify(event Event) {
	for _, handler := range e.handlers[event.Name()] {
		handler.Notify(event)
	}
}
