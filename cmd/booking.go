/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"alinea.com/internal/booking"
	"alinea.com/internal/core"
	"alinea.com/pkg/utils"
	"github.com/spf13/cobra"
)

var b *string
var s *string
var e *string

// bookingCmd represents the booking command
var bookingCmd = &cobra.Command{
	Use:   "booking",
	Short: "handle the bookings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *b != "" {
			if *s == "" || *e == "" {
				panic(fmt.Errorf("start and end are required"))
			}

			parser, err := booking.FromJSON([]byte(*b))
			if err != nil {
				panic(err)
			}

			b, err := parser.ToJSONStruct()
			if err != nil {
				panic(err)
			}

			w, err := core.NewWindow(b.From, b.To)
			if err != nil {
				panic(err)
			}

			bookingService.Book(context.Background(), booking.CreateBookingDTO{
				AgendaID: "1",
				Window:   w,
				Service: core.Service{
					Name:     "My Service",
					Duration: utils.Must(time.ParseDuration("1h")),
				},
			})

			fmt.Println("booking created")
		}
	},
}

var bookingService booking.BookingService

func init() {
	rootCmd.AddCommand(bookingCmd)

	b = bookingCmd.Flags().StringP("book", "b", "", "book an agenda")
	s = bookingCmd.Flags().StringP("start", "s", "", "book start example: 2023-01-01 10:00:00")
	e = bookingCmd.Flags().StringP("end", "e", "", "book end example: 2023-01-01 11:00:00")
}
