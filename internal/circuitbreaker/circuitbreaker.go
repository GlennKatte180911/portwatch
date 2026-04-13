// Package circuitbreaker implements a simple circuit breaker pattern
// for protecting downstream services (e.g. webhook endpoints) from
// repeated failures. It transitions between Closed, Open, and Half-Open
// states based on consecutive failure counts and a recovery timeout.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the Open state.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // blocking calls
	StateHalfOpen              // testing recovery
)

// Breaker is a circuit breaker that tracks failures and controls execution.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	maxFailures  int
	resetTimeout time.Duration
	nextAttempt  time.Time
}

// New creates a Breaker that opens after maxFailures consecutive failures
// and attempts recovery after resetTimeout.
func New(maxFailures int, resetTimeout time.Duration) *Breaker {
	return &Breaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Allow returns nil if the call is permitted, or ErrOpen if the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Now().After(b.nextAttempt) {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess resets the breaker to Closed on a successful call.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure count and opens the circuit if the
// threshold is exceeded.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.nextAttempt = time.Now().Add(b.resetTimeout)
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
