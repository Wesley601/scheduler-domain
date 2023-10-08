package core

type Booking struct {
	Window Window
}

type BookedEvent struct{}

func (b BookedEvent) Name() string {
	return "booking.booked"
}

func (b BookedEvent) IsAsynchronous() bool {
	return false
}

type BookedErrorEvent struct{}

func (b BookedErrorEvent) Name() string {
	return "booking.booked.error"
}

func (b BookedErrorEvent) IsAsynchronous() bool {
	return false
}
