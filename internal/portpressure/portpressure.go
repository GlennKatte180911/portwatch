// Package portpressure tracks how frequently ports appear across successive
// scans and computes a "pressure" score that reflects how consistently a port
// is observed. A high score indicates a stable, repeatedly-seen port; a low
// score indicates a transient or noisy one.
package portpressure

import (
	"sync"
	"time"
)

// Entry holds the pressure data for a single port.
type Entry struct {
	Port      int
	Score     float64   // rolling average presence rate [0.0, 1.0]
	LastSeen  time.Time
	ScanCount int
}

// Tracker maintains pressure scores for observed ports.
type Tracker struct {
	mu      sync.Mutex
	entries map[int]*Entry
	alpha   float64 // EMA smoothing factor
}

// New returns a Tracker with the given EMA smoothing factor alpha.
// alpha must be in (0, 1]; a typical value is 0.3.
func New(alpha float64) *Tracker {
	if alpha <= 0 || alpha > 1 {
		alpha = 0.3
	}
	return &Tracker{
		entries: make(map[int]*Entry),
		alpha:   alpha,
	}
}

// Observe updates pressure scores given the set of ports seen in the latest
// scan. Ports present in seen receive a positive signal; ports previously
// tracked but absent receive a negative signal.
func (t *Tracker) Observe(seen []int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	seenSet := make(map[int]struct{}, len(seen))
	for _, p := range seen {
		seenSet[p] = struct{}{}
	}

	// Positive signal for observed ports.
	for p := range seenSet {
		e, ok := t.entries[p]
		if !ok {
			e = &Entry{Port: p}
			t.entries[p] = e
		}
		e.Score = t.alpha*1.0 + (1-t.alpha)*e.Score
		e.LastSeen = now
		e.ScanCount++
	}

	// Negative signal for tracked but absent ports.
	for p, e := range t.entries {
		if _, present := seenSet[p]; !present {
			e.Score = t.alpha*0.0 + (1-t.alpha)*e.Score
		}
	}
}

// Get returns the Entry for a port and whether it exists.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Above returns all ports whose pressure score is strictly above threshold.
func (t *Tracker) Above(threshold float64) []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	var result []Entry
	for _, e := range t.entries {
		if e.Score > threshold {
			result = append(result, *e)
		}
	}
	return result
}

// Reset clears all tracked pressure data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int]*Entry)
}
