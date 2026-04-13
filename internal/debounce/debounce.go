// Package debounce provides a mechanism to suppress rapid repeated port-change
// events, ensuring that alerts are only emitted once a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Event represents a debounced port-change event keyed by an arbitrary string
// identifier (e.g. a port number or event type).
type Event struct {
	Key       string
	FiredAt   time.Time
}

// Debouncer delays forwarding of events until no new event with the same key
// has arrived within the configured wait duration.
type Debouncer struct {
	wait    time.Duration
	mu      sync.Mutex
	timers  map[string]*time.Timer
	output  chan Event
}

// New creates a Debouncer that waits for the given duration of silence before
// emitting an event on the returned channel.
func New(wait time.Duration, bufSize int) (*Debouncer, <-chan Event) {
	ch := make(chan Event, bufSize)
	d := &Debouncer{
		wait:   wait,
		timers: make(map[string]*time.Timer),
		output: ch,
	}
	return d, ch
}

// Push schedules an event for the given key. If an event for the same key is
// already pending, its timer is reset, effectively extending the quiet period.
func (d *Debouncer) Push(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Reset(d.wait)
		return
	}

	d.timers[key] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()

		d.output <- Event{Key: key, FiredAt: time.Now()}
	})
}

// Stop cancels all pending timers and closes the output channel. The Debouncer
// must not be used after Stop is called.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for k, t := range d.timers {
		t.Stop()
		delete(d.timers, k)
	}
	close(d.output)
}
