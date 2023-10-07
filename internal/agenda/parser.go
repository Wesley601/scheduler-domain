package agenda

import (
	"encoding/json"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/utils"
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

type ListAgendaJSON struct {
	Data []AgendaJSON `json:"data"`
}

type ListParser struct {
	agendas []core.Agenda
}

func (p ListParser) FromAgenda(a ListAgendaJSON) (ListParser, error) {
	var parser ListParser

	for _, a := range a.Data {
		parser.agendas = append(parser.agendas, utils.Must(assembleAgenda(a.Name, a.Slots)))
	}

	return parser, nil
}

func (p ListParser) FromJSON(b []byte) (ListParser, error) {
	var a ListAgendaJSON

	err := json.Unmarshal(b, &a)
	if err != nil {
		return ListParser{}, err
	}

	return p.FromAgenda(a)
}

func (p *ListParser) ToJSON() ([]byte, error) {
	a, err := p.ToJSONStruct()
	if err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(a, "", "  ")
}

func (p *ListParser) ToAgenda() ([]core.Agenda, error) {
	return p.agendas, nil
}

func (p *ListParser) ToJSONStruct() ([]AgendaJSON, error) {

	var listAgenda []AgendaJSON

	for _, agenda := range p.agendas {
		var slots []Slot
		for _, slot := range agenda.Slots {
			slots = append(slots, Slot{
				Weekday: slot.Weekday,
				StartAt: string(slot.StartAt),
				EndsAt:  string(slot.EndsAt),
			})
		}

		listAgenda = append(listAgenda, AgendaJSON{
			Name:  agenda.Name,
			Slots: slots,
		})
	}

	return listAgenda, nil
}

type Parser struct {
	agenda core.Agenda
}

func (p Parser) FromJSON(b []byte) (Parser, error) {
	var a AgendaJSON

	err := json.Unmarshal(b, &a)
	if err != nil {
		return Parser{}, err
	}

	return p.FromAgenda(a)
}

func (p Parser) FromAgenda(a AgendaJSON) (Parser, error) {
	var parser Parser

	parser.agenda = utils.Must(assembleAgenda(a.Name, a.Slots))

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

func assembleAgenda(name string, slots []Slot) (core.Agenda, error) {
	var a core.Agenda

	a.Name = name

	for _, slot := range slots {
		s, err := core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))
		if err != nil {
			return core.Agenda{}, err
		}

		a.Slots = append(a.Slots, s)
	}

	return a, nil
}
