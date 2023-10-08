package booking

import (
	"context"
	"fmt"

	"alinea.com/internal/core"
	"alinea.com/pkg/event"
)

type BookingRepository interface {
	IsAvailable(c context.Context, w core.Window) (bool, error)
	Save(c context.Context, b core.Booking) error
}

type AgendaRepository interface {
	FindByID(c context.Context, id string) (core.Agenda, error)
}

type BlockRepository interface {
	IsAvailable(c context.Context, w core.Window) (bool, error)
}

type EventPublisher interface {
	Notify(e event.Event)
}

type BookingService struct {
	bookingRepository BookingRepository
	agendaRepository  AgendaRepository
	blockRepository   BlockRepository
	publisher         EventPublisher
}

func NewBookingService(bookingRepository BookingRepository, agendaRepository AgendaRepository, blockRepository BlockRepository, publisher EventPublisher) *BookingService {
	return &BookingService{
		bookingRepository: bookingRepository,
		agendaRepository:  agendaRepository,
		blockRepository:   blockRepository,
		publisher:         publisher,
	}
}

type CreateBookingDTO struct {
	AgendaID string
	Window   core.Window
	Service  core.Service
}

func (useCase *BookingService) Book(c context.Context, dto CreateBookingDTO) (Parser, error) {
	var parser Parser
	schedule, err := useCase.agendaRepository.FindByID(c, dto.AgendaID)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	b := core.Booking{
		Window: dto.Window,
	}

	fits, err := schedule.Fits(b, dto.Service)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	if !fits {
		return parser, fmt.Errorf("booking does not fit in schedule")
	}

	available, err := useCase.bookingRepository.IsAvailable(c, b.Window)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	if !available {
		return parser, fmt.Errorf("already booked")
	}

	available, err = useCase.blockRepository.IsAvailable(c, b.Window)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	if !available {
		return parser, fmt.Errorf("blocked")
	}

	err = useCase.bookingRepository.Save(c, b)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	useCase.publisher.Notify(core.BookedEvent{})

	return parser, err
}
