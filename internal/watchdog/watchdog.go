// Package watchdog provides a self-monitoring component that detects when
// the scan loop has stalled and emits an alert if no scan has completed
// within the expected interval plus a configurable grace period.
package watchdog

import (
	"context"
	"sync"
	"time"
)

// StalledFunc is called when the watchdog determines the scanner has stalled.
type StalledFunc func(lastSeen time.Time)

// Watchdog monitors scan heartbeats and fires a callback when they stop.
type Watchdog struct {
	mu       sync.Mutex
	lastBeat time.Time
	interval time.Duration
	grace    time.Duration
	onStall  StalledFunc
}

// New creates a Watchdog that expects a heartbeat at least every interval+grace.
// onStall is invoked (once per stall event) when the deadline is exceeded.
func New(interval, grace time.Duration, onStall StalledFunc) *Watchdog {
	return &Watchdog{
		lastBeat: time.Now(),
		interval: interval,
		grace:    grace,
		onStall:  onStall,
	}
}

// Beat records that a scan has just completed. Call this after every scan.
func (w *Watchdog) Beat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastBeat = time.Now()
}

// LastBeat returns the time of the most recent heartbeat.
func (w *Watchdog) LastBeat() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastBeat
}

// Run starts the watchdog loop. It blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval + w.grace)
	defer ticker.Stop()

	stalled := false

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			w.mu.Lock()
			last := w.lastBeat
			w.mu.Unlock()

			if now.Sub(last) > w.interval+w.grace {
				if !stalled {
					stalled = true
					if w.onStall != nil {
						w.onStall(last)
					}
				}
			} else {
				stalled = false
			}
		}
	}
}
