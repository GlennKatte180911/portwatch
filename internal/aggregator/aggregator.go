// Package aggregator batches port change events over a time window
// before forwarding them downstream, reducing notification noise during
// rapid port churn.
package aggregator

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Event holds an aggregated diff collected within one flush window.
type Event struct {
	Added   []int
	Removed []int
	At      time.Time
}

// Aggregator collects snapshot diffs and emits batched Events.
type Aggregator struct {
	mu       sync.Mutex
	added    map[int]struct{}
	removed  map[int]struct{}
	window   time.Duration
	output   chan Event
	stopOnce sync.Once
	done     chan struct{}
}

// New creates an Aggregator that flushes accumulated diffs every window.
func New(window time.Duration) *Aggregator {
	a := &Aggregator{
		added:   make(map[int]struct{}),
		removed: make(map[int]struct{}),
		window:  window,
		output:  make(chan Event, 8),
		done:    make(chan struct{}),
	}
	go a.loop()
	return a
}

// Push records a diff produced by snapshot.Diff.
func (a *Aggregator) Push(d snapshot.Diff) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, p := range d.Added {
		a.added[p] = struct{}{}
		delete(a.removed, p)
	}
	for _, p := range d.Removed {
		a.removed[p] = struct{}{}
		delete(a.added, p)
	}
}

// Events returns the channel on which batched Events are delivered.
func (a *Aggregator) Events() <-chan Event {
	return a.output
}

// Stop halts the flush loop and closes the output channel.
func (a *Aggregator) Stop() {
	a.stopOnce.Do(func() { close(a.done) })
}

func (a *Aggregator) loop() {
	ticker := time.NewTicker(a.window)
	defer ticker.Stop()
	defer close(a.output)
	for {
		select {
		case <-ticker.C:
			a.flush()
		case <-a.done:
			a.flush()
			return
		}
	}
}

func (a *Aggregator) flush() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.added) == 0 && len(a.removed) == 0 {
		return
	}
	ev := Event{At: time.Now()}
	for p := range a.added {
		ev.Added = append(ev.Added, p)
	}
	for p := range a.removed {
		ev.Removed = append(ev.Removed, p)
	}
	a.added = make(map[int]struct{})
	a.removed = make(map[int]struct{})
	a.output <- ev
}
