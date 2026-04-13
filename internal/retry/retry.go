// Package retry provides a simple exponential-backoff retry mechanism
// for transient failures such as webhook delivery or file I/O errors.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config controls the retry behaviour.
type Config struct {
	// MaxAttempts is the total number of tries (including the first).
	MaxAttempts int
	// BaseDelay is the wait time before the second attempt.
	BaseDelay time.Duration
	// MaxDelay caps the exponential growth.
	MaxDelay time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 4,
		BaseDelay:   250 * time.Millisecond,
		MaxDelay:    10 * time.Second,
	}
}

// Retryer executes operations with exponential backoff.
type Retryer struct {
	cfg Config
}

// New creates a Retryer with the given Config.
func New(cfg Config) *Retryer {
	return &Retryer{cfg: cfg}
}

// Do calls fn repeatedly until it returns nil, the context is cancelled,
// or MaxAttempts is exhausted. It returns ErrMaxAttempts when all tries fail.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if lastErr = fn(); lastErr == nil {
			return nil
		}
		if attempt == r.cfg.MaxAttempts-1 {
			break
		}
		delay := r.delay(attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return errors.Join(ErrMaxAttempts, lastErr)
}

// delay computes the backoff duration for a given attempt index (0-based).
func (r *Retryer) delay(attempt int) time.Duration {
	d := float64(r.cfg.BaseDelay) * math.Pow(2, float64(attempt))
	if d > float64(r.cfg.MaxDelay) {
		d = float64(r.cfg.MaxDelay)
	}
	return time.Duration(d)
}
