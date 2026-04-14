// Package porttrend tracks how frequently individual ports appear across
// successive scans, providing a simple open-frequency metric per port.
package porttrend

import (
	"sync"
	"time"
)

// Entry holds trend data for a single port.
type Entry struct {
	Port      int
	SeenCount int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Tracker accumulates per-port observation counts.
type Tracker struct {
	mu      sync.RWMutex
	entries map[int]*Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[int]*Entry),
	}
}

// Record marks each port in ports as observed at the given timestamp.
func (t *Tracker) Record(ports []int, at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, p := range ports {
		e, ok := t.entries[p]
		if !ok {
			e = &Entry{Port: p, FirstSeen: at}
			t.entries[p] = e
		}
		e.SeenCount++
		e.LastSeen = at
	}
}

// Get returns the Entry for the given port and whether it was found.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	e, ok := t.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all tracked entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

// Reset clears all tracked data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int]*Entry)
}
