package app

import (
	"context"

	"alinea.com/internal/agenda"
	"alinea.com/internal/booking"
	"alinea.com/internal/service"
	"alinea.com/pkg/mongo"
	"alinea.com/pkg/utils"
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
