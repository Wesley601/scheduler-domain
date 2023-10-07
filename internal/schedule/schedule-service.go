package schedule

import (
	"alinea.com/internal/core"
)

type ScheduleRepository interface {
	FindByID(id string) (core.Schedule, error)
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
