package agenda

import (
	"context"
	"time"

	"alinea.com/internal/core"
)

type AgendaRepository interface {
	FindByID(c context.Context, id string) (core.Agenda, error)
	Save(c context.Context, s core.Agenda) error
	List(c context.Context) ([]core.Agenda, error)
}

type AgendaService struct {
	agendaRepository AgendaRepository
}

func NewAgendaService(agendaRepository AgendaRepository) *AgendaService {
	return &AgendaService{
		agendaRepository: agendaRepository,
	}
}

func (useCase *AgendaService) ListSlots(c context.Context, id string, w core.Window, s core.Service) ([]core.Window, error) {
	a, err := useCase.agendaRepository.FindByID(c, id)
	if err != nil {
		return []core.Window{}, err
	}

	slots, err := a.ListAvailableSlots(w, s)
	if err != nil {
		return []core.Window{}, err
	}

	return slots, nil
}

type CreateSlotDTO struct {
	Weekday time.Weekday
	StartAt string
	EndsAt  string
}

type CreateAgendaDTO struct {
	Name  string
	Slots []CreateSlotDTO
}

func (useCase *AgendaService) Create(c context.Context, dto CreateAgendaDTO) (Parser, error) {
	var slots []core.Slot

	for _, slot := range dto.Slots {
		s, err := core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))
		if err != nil {
			return Parser{}, err
		}

		slots = append(slots, s)
	}

	agenda := core.NewAgenda(dto.Name, slots)

	if err := useCase.agendaRepository.Save(c, *agenda); err != nil {
		return Parser{}, err
	}

	parser := Parser{
		agenda: *agenda,
	}

	return parser, nil
}

func (useCase *AgendaService) FindByID(c context.Context, id string) (Parser, error) {
	agenda, err := useCase.agendaRepository.FindByID(c, id)
	if err != nil {
		return Parser{}, err
	}

	parser := Parser{
		agenda: agenda,
	}

	return parser, nil
}

func (useCase *AgendaService) List(c context.Context) (ListParser, error) {
	agenda, err := useCase.agendaRepository.List(c)
	if err != nil {
		return ListParser{}, err
	}

	return ListParser{
		agendas: agenda,
	}, nil
}
