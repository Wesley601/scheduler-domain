package agenda

import (
	"context"
	"time"

	"alinea.com/internal/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgendaRepository interface {
	FindByID(c context.Context, id string) (core.Agenda, error)
	Save(c context.Context, s core.Agenda) error
	List(c context.Context) ([]core.Agenda, error)
}

type ServiceRepository interface {
	FindByID(c context.Context, id string) (core.Service, error)
}

type AgendaService struct {
	agendaRepository  AgendaRepository
	serviceRepository ServiceRepository
}

func NewAgendaService(agendaRepository AgendaRepository, serviceRepository ServiceRepository) *AgendaService {
	return &AgendaService{
		agendaRepository:  agendaRepository,
		serviceRepository: serviceRepository,
	}
}

func (useCase *AgendaService) ListSlots(c context.Context, agendaId, serviceId string, w core.Window) (ListWindowParser, error) {
	var parser ListWindowParser

	a, err := useCase.agendaRepository.FindByID(c, agendaId)
	if err != nil {
		return parser, err
	}

	s, err := useCase.serviceRepository.FindByID(c, serviceId)
	if err != nil {
		return parser, err
	}

	slots, err := a.ListAvailableSlots(w, s)
	if err != nil {
		return parser, err
	}

	parser = ListWindowParser{
		windows: slots,
	}

	return parser, nil
}

type CreateSlotDTO struct {
	Weekday time.Weekday `json:"weekday"`
	StartAt string       `json:"start_at"`
	EndsAt  string       `json:"ends_at"`
}

type CreateAgendaDTO struct {
	Name  string          `json:"name"`
	Slots []CreateSlotDTO `json:"slots"`
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

	agenda := core.NewAgenda(primitive.NewObjectID().Hex(), dto.Name, slots)

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
