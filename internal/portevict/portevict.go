// Package portevict tracks ports that have been evicted (closed) and
// provides a grace-period mechanism before treating them as truly gone.
package portevict

import (
	"sync"
	"time"
)

// Evictor holds ports that have recently disappeared and releases them
// only once their grace period has elapsed.
type Evictor struct {
	mu      sync.Mutex
	grace   time.Duration
	evicted map[int]time.Time // port -> time it was first evicted
}

// New returns an Evictor with the given grace period.
func New(grace time.Duration) *Evictor {
	return &Evictor{
		grace:   grace,
		evicted: make(map[int]time.Time),
	}
}

// Mark records port as evicted at the current time if not already tracked.
func (e *Evictor) Mark(port int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.evicted[port]; !ok {
		e.evicted[port] = time.Now()
	}
}

// Lift removes a port from the eviction list (e.g. it came back up).
func (e *Evictor) Lift(port int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.evicted, port)
}

// Confirmed returns ports whose grace period has fully elapsed.
// Those ports are removed from internal tracking.
func (e *Evictor) Confirmed() []int {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	var out []int
	for port, evictedAt := range e.evicted {
		if now.Sub(evictedAt) >= e.grace {
			out = append(out, port)
			delete(e.evicted, port)
		}
	}
	return out
}

// Pending returns the number of ports currently in the grace window.
func (e *Evictor) Pending() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.evicted)
}
