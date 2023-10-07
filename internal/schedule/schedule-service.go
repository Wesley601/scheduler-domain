package schedule

import (
	"time"

	"alinea.com/internal/core"
)

type ScheduleRepository interface {
	FindByID(id string) (core.Schedule, error)
	Save(s core.Schedule) error
}

type ScheduleService struct {
	scheduleRepository ScheduleRepository
}

func NewScheduleService(scheduleRepository ScheduleRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepository: scheduleRepository,
	}
}

func (useCase *ScheduleService) ListSlots(id string, w core.Window, sv core.Service) ([]core.Window, error) {
	s, err := useCase.scheduleRepository.FindByID(id)
	if err != nil {
		return []core.Window{}, err
	}

	slots, err := s.ListAvailableSlots(w, sv)
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

type CreateScheduleDTO struct {
	Name  string
	Slots []CreateSlotDTO
}

func (useCase *ScheduleService) Create(dto CreateScheduleDTO) (Parser, error) {
	var slots []core.Slot

	for _, slot := range dto.Slots {
		s, err := core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))
		if err != nil {
			return Parser{}, err
		}

		slots = append(slots, s)
	}

	schedule := core.NewSchedule(dto.Name, slots)

	if err := useCase.scheduleRepository.Save(*schedule); err != nil {
		return Parser{}, err
	}

	parser := Parser{
		schedule: *schedule,
	}

	return parser, nil
}

func (useCase *ScheduleService) FindById(id string) (Parser, error) {
	schedule, err := useCase.scheduleRepository.FindByID(id)
	if err != nil {
		return Parser{}, err
	}

	parser := Parser{
		schedule: schedule,
	}

	return parser, nil
}
