package core

import "time"

type Service struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}

func (s *Service) Fits(w Window) bool {
	return w.To.Sub(w.From) == s.Duration
}
