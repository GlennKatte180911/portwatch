// Package suppress provides a mechanism to temporarily silence alerts
// for specific ports, preventing repeated notifications during known
// maintenance windows or expected port state changes.
package suppress

import (
	"sync"
	"time"
)

// Suppressor tracks suppressed ports and their expiry times.
type Suppressor struct {
	mu      sync.Mutex
	entries map[int]time.Time
	now     func() time.Time
}

// New returns a new Suppressor.
func New() *Suppressor {
	return &Suppressor{
		entries: make(map[int]time.Time),
		now:     time.Now,
	}
}

// Suppress silences alerts for the given port for the specified duration.
// Calling Suppress on an already-suppressed port resets its expiry.
func (s *Suppressor) Suppress(port int, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = s.now().Add(duration)
}

// IsSuppressed reports whether the given port is currently suppressed.
func (s *Suppressor) IsSuppressed(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	expiry, ok := s.entries[port]
	if !ok {
		return false
	}
	if s.now().After(expiry) {
		delete(s.entries, port)
		return false
	}
	return true
}

// Lift removes the suppression for the given port immediately.
func (s *Suppressor) Lift(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// Apply filters a slice of ports, returning only those that are not suppressed.
func (s *Suppressor) Apply(ports []int) []int {
	result := make([]int, 0, len(ports))
	for _, p := range ports {
		if !s.IsSuppressed(p) {
			result = append(result, p)
		}
	}
	return result
}
