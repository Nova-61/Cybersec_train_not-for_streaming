package retry

import (
	"fmt"
	"time"
)

type Retrier struct {
	maxAttempts int
	delay       time.Duration
}

func NewRetrier(maxAttempts int, delay time.Duration) *Retrier {
	return &Retrier{maxAttempts: maxAttempts, delay: delay}
}

// Do вызывает fn, повторяя при ошибке до maxAttempts раз с паузой delay между попытками
func (r *Retrier) Do(fn func() error) error {
	var lastErr error
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			fmt.Printf("Attempt %d/%d failed: %v\n", attempt, r.maxAttempts, err)
		}
		if attempt < r.maxAttempts {
			time.Sleep(r.delay)
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", r.maxAttempts, lastErr)
}
