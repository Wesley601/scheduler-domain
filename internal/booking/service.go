package booking

import (
	"context"
	"fmt"

	"alinea.com/internal/core"
	"alinea.com/pkg/mongo"
)

type BookingService struct {
	bookingRepository *mongo.BookingRepository
	agendaRepository  *mongo.AgendaRepository
	blockRepository   *mongo.BlockRepository
	serviceRepository *mongo.ServiceRepository
}

func NewBookingService(
	bookingRepository *mongo.BookingRepository,
	agendaRepository *mongo.AgendaRepository,
	blockRepository *mongo.BlockRepository,
	serviceRepository *mongo.ServiceRepository,
) *BookingService {
	return &BookingService{
		bookingRepository: bookingRepository,
		agendaRepository:  agendaRepository,
		blockRepository:   blockRepository,
		serviceRepository: serviceRepository,
	}
}

type CreateBookingDTO struct {
	AgendaID  string `json:"agenda_id"`
	ServiceID string `json:"service_id"`
	From      string `json:"from"`
	To        string `json:"to"`
}

func (useCase *BookingService) Book(c context.Context, dto CreateBookingDTO) error {
	agenda, err := useCase.agendaRepository.FindByID(c, dto.AgendaID)
	if err != nil {
		return err
	}

	b, err := core.CreateNewBooking(dto.From, dto.To)
	if err != nil {
		return err
	}

	service, err := useCase.serviceRepository.FindByID(c, dto.ServiceID)
	if err != nil {
		return err
	}

	fits, err := agenda.Fits(b, service)
	if err != nil {
		return err
	}

	if !fits {
		return fmt.Errorf("booking does not fit in schedule")
	}

	err = useCase.IsAvailable(c, b.Window)
	if err != nil {
		return err
	}

	return useCase.bookingRepository.Save(c, b)
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
