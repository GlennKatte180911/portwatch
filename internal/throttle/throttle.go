// Package throttle provides a token-bucket style throttle that limits
// how many notifications can be emitted within a sliding time window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks per-key event counts within a rolling window and
// reports whether a new event should be allowed through.
type Throttle struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	buckets map[string][]time.Time
}

// New returns a Throttle that allows at most max events per key within
// the given window duration.
func New(window time.Duration, max int) *Throttle {
	return &Throttle{
		window:  window,
		max:     max,
		buckets: make(map[string][]time.Time),
	}
}

// Allow returns true if the event identified by key is within the
// allowed rate, and records the event. It returns false when the key
// has already reached the maximum count inside the current window.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)

	times := t.buckets[key]
	// evict timestamps outside the window
	valid := times[:0]
	for _, ts := range times {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= t.max {
		t.buckets[key] = valid
		return false
	}

	t.buckets[key] = append(valid, now)
	return true
}

// Reset clears all recorded events for the given key, allowing it to
// pass through immediately on the next call to Allow.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.buckets, key)
}

// Remaining returns the number of additional events that key may emit
// before being throttled within the current window.
func (t *Throttle) Remaining(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)

	count := 0
	for _, ts := range t.buckets[key] {
		if ts.After(cutoff) {
			count++
		}
	}

	rem := t.max - count
	if rem < 0 {
		return 0
	}
	return rem
}
