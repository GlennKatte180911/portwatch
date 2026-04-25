// Package portannounce broadcasts port change events to registered listeners.
package portannounce

import (
	"sync"

	"github.com/user/portwatch/internal/portdiff"
)

// Handler is a function that receives a port diff event.
type Handler func(diff portdiff.Diff)

// Announcer broadcasts diffs to all registered handlers.
type Announcer struct {
	mu       sync.RWMutex
	handlers []Handler
}

// New returns a new Announcer with no registered handlers.
func New() *Announcer {
	return &Announcer{}
}

// Subscribe registers a handler that will be called for each announced diff.
// Returns an unsubscribe function that removes the handler.
func (a *Announcer) Subscribe(h Handler) func() {
	a.mu.Lock()
	defer a.mu.Unlock()

	idx := len(a.handlers)
	a.handlers = append(a.handlers, h)

	return func() {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.handlers[idx] = nil
	}
}

// Announce sends the diff to all registered handlers.
// Nil handlers (unsubscribed) are skipped. Handlers are called sequentially.
func (a *Announcer) Announce(diff portdiff.Diff) {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return
	}

	a.mu.RLock()
	handlers := make([]Handler, len(a.handlers))
	copy(handlers, a.handlers)
	a.mu.RUnlock()

	for _, h := range handlers {
		if h != nil {
			h(diff)
		}
	}
}

// Count returns the number of active (non-nil) subscribers.
func (a *Announcer) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	count := 0
	for _, h := range a.handlers {
		if h != nil {
			count++
		}
	}
	return count
}
