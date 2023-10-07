package schedule

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

type Agenda struct {
	Name  string `json:"name"`
	Slots []Slot `json:"slots"`
}

type Parser struct {
	schedule core.Schedule
}

func (p *Parser) ToAgenda() (Agenda, error) {
	var slots []Slot

	for _, slot := range p.schedule.Slots {
		slots = append(slots, Slot{
			Weekday: slot.Weekday,
			StartAt: string(slot.StartAt),
			EndsAt:  string(slot.EndsAt),
		})
	}

	return Agenda{
		Name:  p.schedule.Name,
		Slots: slots,
	}, nil
}

func (p *Parser) FromAgenda(a Agenda) error {
	var slots []core.Slot

	for _, slot := range a.Slots {
		s, err := core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))
		if err != nil {
			return err
		}

		slots = append(slots, s)
	}

	p.schedule = core.Schedule{
		Name:  a.Name,
		Slots: slots,
	}

	return nil
}

func (p *Parser) ToJSON() ([]byte, error) {
	a, err := p.ToAgenda()
	if err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(a, "", "  ")
}

func (p *Parser) FromJSON(b []byte) error {
	var a Agenda

	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	return p.FromAgenda(a)
}
