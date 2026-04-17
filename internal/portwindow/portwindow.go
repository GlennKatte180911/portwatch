// Package portwindow provides a sliding-window view of port activity,
// tracking how many times each port has been seen open within a recent
// time window.
package portwindow

import (
	"sync"
	"time"
)

// Entry records a single observation timestamp for a port.
type Entry struct {
	Port      int
	ObservedAt time.Time
}

// Window maintains a sliding time window of port observations.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	log      []Entry
}

// New creates a Window that retains observations within the given duration.
func New(duration time.Duration) *Window {
	return &Window{duration: duration}
}

// Observe records that port was seen open at the current time.
func (w *Window) Observe(port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.log = append(w.log, Entry{Port: port, ObservedAt: time.Now()})
	w.evict()
}

// Count returns the number of times port has been observed within the window.
func (w *Window) Count(port int) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	count := 0
	for _, e := range w.log {
		if e.Port == port {
			count++
		}
	}
	return count
}

// Active returns all ports that have at least one observation in the window.
func (w *Window) Active() []int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	seen := make(map[int]struct{})
	for _, e := range w.log {
		seen[e.Port] = struct{}{}
	}
	ports := make([]int, 0, len(seen))
	for p := range seen {
		ports = append(ports, p)
	}
	return ports
}

// Reset clears all observations.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.log = nil
}

// evict removes entries older than the window duration. Must be called with mu held.
func (w *Window) evict() {
	cutoff := time.Now().Add(-w.duration)
	i := 0
	for i < len(w.log) && w.log[i].ObservedAt.Before(cutoff) {
		i++
	}
	w.log = w.log[i:]
}
