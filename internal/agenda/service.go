package agenda

import (
	"context"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgendaService struct {
	agendaRepository  *mongo.AgendaRepository
	serviceRepository *mongo.ServiceRepository
}

func NewAgendaService(agendaRepository *mongo.AgendaRepository, serviceRepository *mongo.ServiceRepository) *AgendaService {
	return &AgendaService{
		agendaRepository:  agendaRepository,
		serviceRepository: serviceRepository,
	}
}

func (useCase *AgendaService) ListSlots(c context.Context, agendaId, serviceId string, w core.Window) ([]core.Window, error) {
	a, err := useCase.agendaRepository.FindByID(c, agendaId)
	if err != nil {
		return []core.Window{}, err
	}

	s, err := useCase.serviceRepository.FindByID(c, serviceId)
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
	Weekday time.Weekday `json:"weekday"`
	StartAt string       `json:"start_at"`
	EndsAt  string       `json:"ends_at"`
}

type CreateAgendaDTO struct {
	Name  string          `json:"name"`
	Slots []CreateSlotDTO `json:"slots"`
}

func (useCase *AgendaService) Create(c context.Context, dto CreateAgendaDTO) (*core.Agenda, error) {
	var slots []core.Slot

	for _, slot := range dto.Slots {
		s, err := core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))
		if err != nil {
			return nil, err
		}

		slots = append(slots, s)
	}

	agenda := core.NewAgenda(primitive.NewObjectID().Hex(), dto.Name, slots)

	if err := useCase.agendaRepository.Save(c, *agenda); err != nil {
		return nil, err
	}

	return agenda, nil
}

func (useCase *AgendaService) FindByID(c context.Context, id string) (*core.Agenda, error) {
	agenda, err := useCase.agendaRepository.FindByID(c, id)
	if err != nil {
		return nil, err
	}

	return &agenda, nil
}

func (useCase *AgendaService) List(c context.Context) ([]core.Agenda, error) {
	agendas, err := useCase.agendaRepository.List(c)
	if err != nil {
		return nil, err
	}

	return agendas, nil
}
