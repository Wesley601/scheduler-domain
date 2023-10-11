/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"alinea.com/internal/booking"
	"alinea.com/internal/core"
	"alinea.com/pkg/event"
	"github.com/spf13/cobra"
)

// bookingCmd represents the booking command
var bookingCmd = &cobra.Command{
	Use:   "booking",
	Short: "handle the bookings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *s == "" || *e == "" {
			panic(fmt.Errorf("start and end are required"))
		}

		w, err := core.NewWindow(*s, *e)
		if err != nil {
			panic(err)
		}

		book, err := bookingService.Book(context.Background(), booking.CreateBookingDTO{
			AgendaID:  *aID,
			Window:    w,
			ServiceID: *seID,
		})
		if err != nil {
			panic(err)
		}

		j, err := book.ToJSON()
		if err != nil {
			panic(err)
		}

		fmt.Println(string(j))
	},
}

var aID *string
var seID *string
var s *string
var e *string

func init() {
	rootCmd.AddCommand(bookingCmd)

	seID = bookingCmd.Flags().StringP("serviceID", "s", "", "service id")
	aID = bookingCmd.Flags().StringP("agendaID", "a", "", "agenda id")
	s = bookingCmd.Flags().StringP("start", "t", "", "book start example: 2023-01-01 10:00:00")
	e = bookingCmd.Flags().StringP("end", "e", "", "book end example: 2023-01-01 11:00:00")
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
