package agenda

import (
	"encoding/json"
	"time"

	"alinea.com/internal/core"
)

type Slot struct {
	Weekday time.Weekday `json:"weekday"`
	StartAt string       `json:"start_at"`
	EndsAt  string       `json:"ends_at"`
}

type AgendaJSON struct {
	Name  string `json:"name"`
	Slots []Slot `json:"slots"`
}

type Parser struct {
	agenda core.Agenda
}

func FromJSON(b []byte) (Parser, error) {
	var a AgendaJSON

	err := json.Unmarshal(b, &a)
	if err != nil {
		return Parser{}, err
	}

	return FromAgenda(a)
}

func FromAgenda(a AgendaJSON) (Parser, error) {
	var parser Parser
	var slots []core.Slot

	for _, slot := range a.Slots {
		s, err := core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))
		if err != nil {
			return parser, err
		}

		slots = append(slots, s)
	}

	parser.agenda = core.Agenda{
		Name:  a.Name,
		Slots: slots,
	}

	return parser, nil
}

func (p *Parser) ToJSON() ([]byte, error) {
	a, err := p.ToJSONStruct()
	if err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(a, "", "  ")
}

func (p *Parser) ToAgenda() (core.Agenda, error) {
	return p.agenda, nil
}

func (p *Parser) ToJSONStruct() (AgendaJSON, error) {
	var slots []Slot

	for _, slot := range p.agenda.Slots {
		slots = append(slots, Slot{
			Weekday: slot.Weekday,
			StartAt: string(slot.StartAt),
			EndsAt:  string(slot.EndsAt),
		})
	}

	return AgendaJSON{
		Name:  p.agenda.Name,
		Slots: slots,
	}, nil
}
