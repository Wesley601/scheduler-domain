package main

import (
	"fmt"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/utils"
)

func main() {
	service := core.Service{
		Name:     "My Service",
		Duration: utils.Must(time.ParseDuration("5m")),
	}

	s := core.Schedule{
		Name: "My Schedule",
		Slots: []core.Slot{
			utils.Must(core.NewSlot(time.Monday, "10:00:00", "10:45:00")),
			utils.Must(core.NewSlot(time.Tuesday, "11:00:00", "11:45:00")),
		},
	}

	w, err := s.ListAvailableSlots(core.Window{
		From: utils.Must(time.Parse(time.RFC3339, "2023-10-01T00:00:00Z")),
		To:   utils.Must(time.Parse(time.RFC3339, "2023-10-05T00:00:00Z")),
	}, service)

	if err != nil {
		panic(err)
	}

	for _, w2 := range w {
		fmt.Printf("%s - %s %s\n", w2.From.Format(time.RFC3339), w2.To.Format(time.RFC3339), w2.From.Weekday())
	}

	result, err := s.Fits(core.Booking{
		Window: core.Window{
			From: utils.Must(time.Parse(time.RFC3339, "2023-10-03T15:30:00Z")),
			To:   utils.Must(time.Parse(time.RFC3339, "2023-10-03T15:35:00Z")),
		},
	}, service)
	if err != nil {
		panic(err)
	}

	fmt.Printf("result: %v\n", result)
}
