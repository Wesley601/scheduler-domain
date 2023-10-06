package core

import (
	"fmt"
	"time"

	"alinea.com/pkg/utils"
)

type Schedule struct {
	Name  string
	Slots []Slot
}

func (s *Schedule) ListAvailableSlots(w Window, sv Service) (ws []Window, err error) {
	startAt := w.From
	endsAt := w.To

	for !startAt.After(endsAt) {
		for _, s2 := range s.Slots {
			if s2.Weekday != startAt.Weekday() {
				continue
			}

			from := utils.Must(s2.StartAt.Time(startAt))
			to := utils.Must(s2.EndsAt.Time(startAt))

			for from.Before(to) {
				ws = append(ws, Window{
					From: from,
					To:   from.Add(sv.Duration),
				})

				from = from.Add(sv.Duration)
			}
		}

		startAt = startAt.Add(24 * time.Hour)
	}

	return
}

func (s *Schedule) Fits(b Booking, sv Service) (bool, error) {
	if !sv.WindowFits(b.Window) {
		return false, fmt.Errorf("invalid booking duration")
	}

	slots, err := s.ListAvailableSlots(b.Window, sv)
	if err != nil {
		return false, err
	}

	for _, slot := range slots {
		if b.Window.Equal(slot) {
			return true, nil
		}
	}

	return false, nil
}

type ScheduleService struct {
	BookingRepository  interface{}
	ScheduleRepository interface{}
	BlockRepository    interface{}
}
