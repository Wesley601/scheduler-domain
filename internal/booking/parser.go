package booking

import (
	"encoding/json"

	"alinea.com/internal/core"
)

type BookingJSON struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
}

type Parser struct {
	booking core.Booking
}

func FromJSON(b []byte) (Parser, error) {
	var a BookingJSON

	err := json.Unmarshal(b, &a)
	if err != nil {
		return Parser{}, err
	}

	return FromBooking(a)
}

func FromBooking(a BookingJSON) (Parser, error) {
	var parser Parser

	window, err := core.NewWindow(a.From, a.To)
	if err != nil {
		return parser, err
	}

	parser.booking = core.Booking{
		ID:     a.ID,
		Window: window,
	}

	return parser, nil
}

func (p *Parser) ToBooking() (core.Booking, error) {
	return p.booking, nil
}

func (p *Parser) ToJSON() ([]byte, error) {
	b, err := p.ToBooking()
	if err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(b, "", "  ")
}

func (p *Parser) ToJSONStruct() (BookingJSON, error) {
	return BookingJSON{
		ID:   p.booking.ID,
		From: p.booking.Window.From.Format("2006-01-02 15:04:05"),
		To:   p.booking.Window.To.Format("2006-01-02 15:04:05"),
	}, nil
}
