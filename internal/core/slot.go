package core

import (
	"fmt"
	"time"
)

type SlotTime string

func (s SlotTime) Time(baseTime time.Time) (time.Time, error) {
	return time.Parse(time.DateTime, baseTime.Format("2006-01-02")+" "+string(s))
}

type Slot struct {
	Weekday  time.Weekday
	Duration time.Duration
	StartAt  SlotTime
	EndsAt   SlotTime
}

func NewSlot(weekday time.Weekday, startAt, endsAt SlotTime) (Slot, error) {
	startAtTime, err := startAt.Time(time.Now())
	if err != nil {
		return Slot{}, err
	}
	endsAtTime, err := endsAt.Time(time.Now())
	if err != nil {
		return Slot{}, err
	}

	if endsAtTime.Before(startAtTime) {
		return Slot{}, fmt.Errorf("endsAt must be after startAt")
	}

	if err != nil {
		return Slot{}, err
	}

	return Slot{
		Weekday:  weekday,
		Duration: endsAtTime.Sub(startAtTime),
		StartAt:  startAt,
		EndsAt:   endsAt,
	}, nil
}
