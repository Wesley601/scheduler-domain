package schedule

import "time"

type Service struct {
	Name     string
	Duration time.Duration
}

func (s *Service) WindowFits(w Window) bool {
	return w.To.Sub(w.From) == s.Duration
}
