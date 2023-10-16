package utils

import (
	"context"
	"fmt"
	"time"
)

// Must panics if e is not nil.
// Do not use in production!
// This is just a helper for the examples or test environment.
func Must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}

	return t
}

func LoadingMessage(c context.Context, message string) {
	symbols := []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"}
	for {
		select {
		case <-c.Done():
			fmt.Print("\r\033[K")
			return
		default:
			for _, s := range symbols {
				fmt.Printf("\r\033[K%s %s", s, message)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
