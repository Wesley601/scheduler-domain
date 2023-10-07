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
	FindByID(id string) (core.Agenda, error)
}

type BlockRepository interface {
	IsAvailable(w core.Window) (bool, error)
}

type BookingService struct {
	bookingRepository  BookingRepository
	scheduleRepository ScheduleRepository
	blockRepository    BlockRepository
}

func NewBookingService(bookingRepository BookingRepository, scheduleRepository ScheduleRepository, blockRepository BlockRepository) *BookingService {
	return &BookingService{
		bookingRepository:  bookingRepository,
		scheduleRepository: scheduleRepository,
		blockRepository:    blockRepository,
	}
}

func (useCase *BookingService) Book(sID string, b core.Booking, s core.Service) error {
	schedule, err := useCase.scheduleRepository.FindByID(sID)
	if err != nil {
		return err
	}

	fits, err := schedule.Fits(b, s)
	if err != nil {
		return err
	}

	if !fits {
		return fmt.Errorf("booking does not fit in schedule")
	}

	available, err := useCase.bookingRepository.IsAvailable(b.Window)
	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("already booked")
	}

	available, err = useCase.blockRepository.IsAvailable(b.Window)
	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("blocked")
	}

	return useCase.bookingRepository.Save(b)
}
