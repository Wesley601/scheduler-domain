package core

import "time"

type Block struct {
	ID      string
	Weekday time.Weekday
	Window  Window
}
