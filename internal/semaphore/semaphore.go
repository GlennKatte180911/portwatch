// Package semaphore provides a counting semaphore to bound concurrent
// port-scan goroutines, preventing resource exhaustion during wide-range scans.
package semaphore

import (
	"context"
	"fmt"
)

// Semaphore is a counting semaphore backed by a buffered channel.
type Semaphore struct {
	slots chan struct{}
}

// New returns a Semaphore that allows at most n concurrent acquisitions.
// It returns an error if n is less than 1.
func New(n int) (*Semaphore, error) {
	if n < 1 {
		return nil, fmt.Errorf("semaphore: capacity must be at least 1, got %d", n)
	}
	return &Semaphore{slots: make(chan struct{}, n)}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
// It returns ctx.Err() if the context is done before a slot is obtained.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees one slot. It panics if called more times than Acquire.
func (s *Semaphore) Release() {
	select {
	case <-s.slots:
	default:
		panic("semaphore: Release called without a matching Acquire")
	}
}

// Cap returns the maximum concurrency the semaphore was created with.
func (s *Semaphore) Cap() int {
	return cap(s.slots)
}

// Available returns the number of additional acquisitions currently possible.
func (s *Semaphore) Available() int {
	return cap(s.slots) - len(s.slots)
}
