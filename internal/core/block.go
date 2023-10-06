package core

import "time"

type Block struct {
	Weekday time.Weekday
	Window  Window
}
