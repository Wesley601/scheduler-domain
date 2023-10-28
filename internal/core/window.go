package core

import (
	"fmt"
	"time"
)

type Window struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Create a new Window based on RFC3339 format
func NewWindowFromString(from, to string) (Window, error) {
	f, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return Window{}, err
	}

	t, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return Window{}, err
	}

	return NewWindow(f, t)
}

func NewWindow(from, to time.Time) (Window, error) {
	if to.Before(from) {
		return Window{}, fmt.Errorf("to must be after from")
	}

	return Window{
		From: from,
		To:   to,
	}, nil
}

func (w *Window) Colapse(window Window) bool {
	return w.From.Before(window.To) && w.To.After(window.From)
}

func (w *Window) Equal(window Window) bool {
	return w.From.Equal(window.From) && w.To.Equal(window.To)
}
