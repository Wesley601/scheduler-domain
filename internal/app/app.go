package app

import (
	"context"
	"fmt"

	"alinea.com/internal/agenda"
	"alinea.com/internal/booking"
	"alinea.com/internal/core"
	"alinea.com/internal/service"
	"alinea.com/pkg/event"
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

	BookingService = booking.NewBookingService(bookingRepository, agendaRepository, blockRepository, serviceRepository, createEventPublisher())
	AgendaService = agenda.NewAgendaService(agendaRepository, serviceRepository)
	ServiceService = service.NewServiceService(serviceRepository)
}

type ApplicationErrors struct{}

func (a ApplicationErrors) Notify(event event.Event) {
	fmt.Printf("Error at %s", event.Name())
}

type ApplicationEvents struct{}

func (a ApplicationEvents) Notify(event event.Event) {
	fmt.Printf("Event at %s", event.Name())
}

func createEventPublisher() *event.EventPublisher {
	p := event.NewEventPublisher()

	p.Subscribe(ApplicationErrors{}, core.BookedErrorEvent{})
	p.Subscribe(ApplicationEvents{}, core.BookedEvent{})

	return p
}
