// Package dedup provides event deduplication for portwatch.
// It suppresses repeated identical diffs within a configurable window,
// ensuring downstream consumers are not flooded with duplicate alerts.
package dedup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Deduplicator filters out duplicate diffs seen within a time window.
type Deduplicator struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[string]time.Time
}

// New creates a Deduplicator with the given deduplication window.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		window: window,
		seen:   make(map[string]time.Time),
	}
}

// IsDuplicate returns true if an identical diff was already seen within the window.
func (d *Deduplicator) IsDuplicate(diff snapshot.Diff) bool {
	key := diffKey(diff)
	if key == "" {
		return false
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	if _, ok := d.seen[key]; ok {
		return true
	}
	d.seen[key] = time.Now()
	return false
}

// Reset clears all recorded diffs, allowing them to pass through again.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

// evict removes entries older than the deduplication window. Must be called with mu held.
func (d *Deduplicator) evict() {
	now := time.Now()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

// diffKey produces a stable string key for a Diff.
func diffKey(diff snapshot.Diff) string {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return ""
	}
	return fmt.Sprintf("added=%v removed=%v", diff.Added, diff.Removed)
}
