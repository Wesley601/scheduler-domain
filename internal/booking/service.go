package booking

import (
	"context"
	"fmt"

	"alinea.com/internal/core"
)

type BookingRepository interface {
	IsAvailable(c context.Context, w core.Window) (bool, error)
	Save(c context.Context, b core.Booking) error
}

type ScheduleRepository interface {
	FindByID(c context.Context, id string) (core.Agenda, error)
}

type BlockRepository interface {
	IsAvailable(c context.Context, w core.Window) (bool, error)
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

type CreateBookingDTO struct {
	AgendaID string
	Window   core.Window
	Service  core.Service
}

func (useCase *BookingService) Book(c context.Context, dto CreateBookingDTO) (Parser, error) {
	var parser Parser
	schedule, err := useCase.scheduleRepository.FindByID(c, dto.AgendaID)
	if err != nil {
		return parser, err
	}

	b := core.Booking{
		Window: dto.Window,
	}

	fits, err := schedule.Fits(b, dto.Service)
	if err != nil {
		return parser, err
	}

	if !fits {
		return parser, fmt.Errorf("booking does not fit in schedule")
	}

	available, err := useCase.bookingRepository.IsAvailable(c, b.Window)
	if err != nil {
		return parser, err
	}

	if !available {
		return parser, fmt.Errorf("already booked")
	}

	available, err = useCase.blockRepository.IsAvailable(c, b.Window)
	if err != nil {
		return parser, err
	}

	if !available {
		return parser, fmt.Errorf("blocked")
	}

	err = useCase.bookingRepository.Save(c, b)

	return parser, err
}
