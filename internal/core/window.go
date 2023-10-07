package core

import "time"

type Window struct {
	From time.Time
	To   time.Time
}

func NewWindow(from, to string) (Window, error) {
	f, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return Window{}, err
	}

	t, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return Window{}, err
	}

	return Window{
		From: f,
		To:   t,
	}, nil
}

func (w *Window) Colapse(window Window) bool {
	return w.From.Before(window.To) && w.To.After(window.From)
}

func (w *Window) Equal(window Window) bool {
	return w.From.Equal(window.From) && w.To.Equal(window.To)
}
