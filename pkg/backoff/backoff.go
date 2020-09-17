package backoff

import "time"

type Backoff interface {
	Duration() time.Duration
}

type staticBackoff struct {
	factor time.Duration
}

func NewStaticBackoff(factor time.Duration) Backoff {
	return &staticBackoff{factor: factor}
}

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
