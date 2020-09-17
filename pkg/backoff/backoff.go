package backoff

import "time"

// Backoff provides the next wait duration.
type Backoff interface {
	Duration() time.Duration
}

type staticBackoff struct {
	factor time.Duration
}

// NewStaticBackoff creates a new backoff of provided factor time duration.
func NewStaticBackoff(factor time.Duration) Backoff {
	return &staticBackoff{factor: factor}
}

// Duration returns the time duration to wait.
func (s *staticBackoff) Duration() time.Duration {
	return s.factor
}

// TODO: Create ExponentialBackoff
// type exponentialBackoff struct {
// 	factor time.Duration
// 	count  int
// }

// func NewExponentialBackoff(factor time.Duration) Backoff {
// 	return &exponentialBackoff{factor: factor}
// }

// func (e *exponentialBackoff) Duration() time.Duration {
// 	e.count++
// 	return e.factor * time.Duration(e.count
// }
