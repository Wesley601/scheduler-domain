package core

import "go.mongodb.org/mongo-driver/bson/primitive"

type Status int

const (
	Booked Status = iota
	Canceled
)

type Booking struct {
	ID     string `json:"id"`
	Window Window `json:"window"`
}

func CreateNewBooking(from, to string) (Booking, error) {
	w, err := NewWindowFromString(from, to)
	if err != nil {
		return Booking{}, err
	}

	b := Booking{
		ID:     primitive.NewObjectID().Hex(),
		Window: w,
	}

	return b, nil
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
