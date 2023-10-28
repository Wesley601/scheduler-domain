package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	b, err := core.CreateNewBooking(dto.From, dto.To)
	if err != nil {
		return err
	}

	agenda, err := useCase.agendaRepository.FindByID(c, dto.AgendaID)
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

	err = useCase.HasSomeBlock(c, b.Window)
	if err != nil {
		return err
	}

	return useCase.bookingRepository.Save(c, b)
}

type RebookingDTO struct {
	BookingID string `json:"booking_id"`
	AgendaID  string `json:"agenda_id"`
	ServiceID string `json:"service_id"`
	From      string `json:"from"`
	To        string `json:"to"`
}

func (useCase *BookingService) Rebook(c context.Context, dto RebookingDTO) error {
	oldBooking, err := useCase.bookingRepository.FindByID(c, dto.BookingID)
	if err != nil {
		return err
	}

	if oldBooking.Window.From.Before(time.Now()) {
		return errors.New("cannot rebook passed bookings")
	}

	err = useCase.Book(c, CreateBookingDTO{
		AgendaID:  dto.AgendaID,
		ServiceID: dto.ServiceID,
		From:      dto.From,
		To:        dto.To,
	})

	if err != nil {
		return err
	}

	return useCase.bookingRepository.Remove(c, oldBooking)
}

func (useCase *BookingService) HasSomeBlock(c context.Context, w core.Window) error {
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
