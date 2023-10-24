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

type ServiceRepository interface {
	FindByID(c context.Context, id string) (core.Service, error)
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
	serviceRepository ServiceRepository
	publisher         EventPublisher
}

func NewBookingService(
	bookingRepository BookingRepository,
	agendaRepository AgendaRepository,
	blockRepository BlockRepository,
	serviceRepository ServiceRepository,
	publisher EventPublisher,
) *BookingService {
	return &BookingService{
		bookingRepository: bookingRepository,
		agendaRepository:  agendaRepository,
		blockRepository:   blockRepository,
		serviceRepository: serviceRepository,
		publisher:         publisher,
	}
}

type CreateBookingDTO struct {
	AgendaID  string `json:"agenda_id"`
	ServiceID string `json:"service_id"`
	From      string `json:"from"`
	To        string `json:"to"`
}

func (useCase *BookingService) Book(c context.Context, dto CreateBookingDTO) (Parser, error) {
	var parser Parser
	agenda, err := useCase.agendaRepository.FindByID(c, dto.AgendaID)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	b, err := core.CreateNewBooking(dto.From, dto.To)
	if err != nil {
		return parser, err
	}

	service, err := useCase.serviceRepository.FindByID(c, dto.ServiceID)
	if err != nil {
		return parser, err
	}

	fits, err := agenda.Fits(b, service)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	if !fits {
		return parser, fmt.Errorf("booking does not fit in schedule")
	}

	err = useCase.IsAvailable(c, b.Window)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	err = useCase.bookingRepository.Save(c, b)
	if err != nil {
		useCase.publisher.Notify(core.BookedErrorEvent{})
		return parser, err
	}

	useCase.publisher.Notify(core.BookedEvent{})
	return parser, err
}

func (useCase *BookingService) IsAvailable(c context.Context, w core.Window) error {
	available, err := useCase.bookingRepository.IsAvailable(c, w)
	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("already booked")
	}

	available, err = useCase.blockRepository.IsAvailable(c, w)
	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("blocked")
	}

	return nil
}
