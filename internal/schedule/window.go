package schedule

import "time"

type Window struct {
	From time.Time
	To   time.Time
}

func (w *Window) Colapse(window Window) bool {
	return w.From.Before(window.To) && w.To.After(window.From)
}

func (w *Window) Equal(window Window) bool {
	return w.From.Equal(window.From) && w.To.Equal(window.To)
}
