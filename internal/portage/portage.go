// Package portage tracks the age of open ports — how long each port has
// been continuously observed as open — and exposes helpers for querying
// ports that have been open longer than a given threshold.
package portage

import (
	"sync"
	"time"
)

// Entry holds age information for a single port.
type Entry struct {
	Port      int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Age returns how long the port has been continuously observed.
func (e Entry) Age() time.Duration {
	return e.LastSeen.Sub(e.FirstSeen)
}

// Tracker records the first and last observation time for each port.
type Tracker struct {
	mu      sync.Mutex
	entries map[int]Entry
	now     func() time.Time
}

// New returns a new Tracker. The now function is used for time injection
// (pass time.Now for production use).
func New(now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	return &Tracker{
		entries: make(map[int]Entry),
		now:     now,
	}
}

// Observe records an observation for the given ports. Ports seen for the
// first time have their FirstSeen set; all ports have LastSeen updated.
func (t *Tracker) Observe(ports []int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	for _, p := range ports {
		e, ok := t.entries[p]
		if !ok {
			e = Entry{Port: p, FirstSeen: now}
		}
		e.LastSeen = now
		t.entries[p] = e
	}
}

// Remove deletes the tracking entry for a port (e.g. when it closes).
func (t *Tracker) Remove(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Get returns the Entry for a port and whether it was found.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	return e, ok
}

// OlderThan returns all entries whose age exceeds the given threshold.
func (t *Tracker) OlderThan(threshold time.Duration) []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	var result []Entry
	for _, e := range t.entries {
		if e.Age() > threshold {
			result = append(result, e)
		}
	}
	return result
}

// All returns a snapshot of every tracked entry.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}

// Reset clears all tracking state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int]Entry)
}
