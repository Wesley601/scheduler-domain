package booking

import (
	"fmt"

	"alinea.com/internal/core"
)

type BookingRepository interface {
	IsAvailable(w core.Window) (bool, error)
	Save(b core.Booking) error
}
type ScheduleRepository interface {
	FindByID(id string) (core.Schedule, error)
}
type BlockRepository interface {
	IsAvailable(w core.Window) (bool, error)
}

type BookingUseCase struct {
	bookingRepository  BookingRepository
	scheduleRepository ScheduleRepository
	blockRepository    BlockRepository
}

func (useCase *BookingUseCase) Book(sID string, b core.Booking, s core.Service) error {
	schedule, err := useCase.scheduleRepository.FindByID(sID)
	if err != nil {
		return err
	}

	fits, err := schedule.Fits(b, s)
	if err != nil {
		return err
	}

	if !fits {
		return fmt.Errorf("already booked or not available")
	}

	return nil
}
