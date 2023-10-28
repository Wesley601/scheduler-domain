package core

import "time"

type Block struct {
	ID      string       `json:"id"`
	Weekday time.Weekday `json:"weekday"`
	Window  Window       `json:"window"`
}
