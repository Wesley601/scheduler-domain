package main

import (
	"fmt"
	"time"

	"alinea.com/internal/schedule"
	"alinea.com/pkg/utils"
)

func main() {
	service := schedule.Service{
		Name:     "My Service",
		Duration: utils.Must(time.ParseDuration("5m")),
	}

	s := schedule.Schedule{
		Name: "My Schedule",
		Slots: []schedule.Slot{
			utils.Must(schedule.NewSlot(time.Monday, "10:00:00", "10:45:00")),
			utils.Must(schedule.NewSlot(time.Tuesday, "11:00:00", "11:45:00")),
		},
	}

	w, err := s.ListAvailableSlots(schedule.Window{
		From: utils.Must(time.Parse(time.RFC3339, "2023-10-01T00:00:00Z")),
		To:   utils.Must(time.Parse(time.RFC3339, "2023-10-05T00:00:00Z")),
	}, service)

	if err != nil {
		panic(err)
	}

	for _, w2 := range w {
		fmt.Printf("%s - %s %s\n", w2.From.Format(time.RFC3339), w2.To.Format(time.RFC3339), w2.From.Weekday())
	}

	result, err := s.Fits(schedule.Booking{
		Window: schedule.Window{
			From: utils.Must(time.Parse(time.RFC3339, "2023-10-03T15:30:00Z")),
			To:   utils.Must(time.Parse(time.RFC3339, "2023-10-03T15:35:00Z")),
		},
	}, service)
	if err != nil {
		panic(err)
	}

	fmt.Printf("result: %v\n", result)
}
