// Package portexpiry tracks how long a port has been continuously open
// and emits an alert when it exceeds a configured maximum age.
package portexpiry

import (
	"sync"
	"time"
)

// Entry records when a port was first observed open.
type Entry struct {
	Port      int
	FirstSeen time.Time
}

// Expiry holds per-port open timestamps and a maximum age threshold.
type Expiry struct {
	mu      sync.Mutex
	entries map[int]Entry
	maxAge  time.Duration
}

// New creates an Expiry tracker with the given maximum open duration.
func New(maxAge time.Duration) *Expiry {
	return &Expiry{
		entries: make(map[int]Entry),
		maxAge:  maxAge,
	}
}

// Observe records a port as currently open. If the port is seen for the
// first time, its FirstSeen timestamp is set to now.
func (e *Expiry) Observe(port int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.entries[port]; !ok {
		e.entries[port] = Entry{Port: port, FirstSeen: time.Now()}
	}
}

// Remove forgets a port (e.g. when it closes).
func (e *Expiry) Remove(port int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.entries, port)
}

// Expired returns all ports that have been open longer than maxAge.
func (e *Expiry) Expired() []Entry {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	var out []Entry
	for _, en := range e.entries {
		if now.Sub(en.FirstSeen) > e.maxAge {
			out = append(out, en)
		}
	}
	return out
}

// Age returns how long the port has been open. ok is false if unknown.
func (e *Expiry) Age(port int) (time.Duration, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	en, ok := e.entries[port]
	if !ok {
		return 0, false
	}
	return time.Since(en.FirstSeen), true
}
