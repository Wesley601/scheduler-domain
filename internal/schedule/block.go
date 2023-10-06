package schedule

import "time"

type Block struct {
	Weekday time.Weekday
	Window  Window
}
