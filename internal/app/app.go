package app

import (
	"context"

	"wesley601.com/internal/agenda"
	"wesley601.com/internal/booking"
	"wesley601.com/internal/service"
	"wesley601.com/pkg/mongo"
	"wesley601.com/pkg/utils"
)

var AgendaService *agenda.AgendaService
var BookingService *booking.BookingService
var ServiceService *service.ServiceService

func init() {
	client := utils.Must(mongo.NewClient(context.Background()))

	bookingRepository := mongo.NewBookingRepository(client)
	agendaRepository := mongo.NewAgendaRepository(client)
	blockRepository := mongo.NewBlockRepository(client)
	serviceRepository := mongo.NewServiceRepository(client)

	BookingService = booking.NewBookingService(bookingRepository, agendaRepository, blockRepository, serviceRepository)
	AgendaService = agenda.NewAgendaService(agendaRepository, serviceRepository)
	ServiceService = service.NewServiceService(serviceRepository)
}
