// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently alert notifications are emitted for a given port.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate at which events are allowed through per port.
type Limiter struct {
	mu       sync.Mutex
	last     map[int]time.Time
	cooldown time.Duration
}

// New returns a Limiter that enforces the given cooldown period between
// successive allowed events for the same port.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		last:     make(map[int]time.Time),
		cooldown: cooldown,
	}
}

// Allow reports whether an event for the given port should be allowed through.
// It returns true the first time a port is seen, and subsequently only after
// the configured cooldown has elapsed since the last allowed event.
func (l *Limiter) Allow(port int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[port]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[port] = now
	return true
}

// Reset clears the recorded timestamp for a port, allowing the next event
// for that port to pass immediately regardless of the cooldown.
func (l *Limiter) Reset(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, port)
}

// ResetAll clears all recorded timestamps.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[int]time.Time)
}
