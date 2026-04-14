// Package portstate tracks the current observed state of ports and
// provides a thread-safe snapshot of which ports are open.
package portstate

import (
	"sync"
	"time"
)

// State holds the most recently observed open ports along with metadata.
type State struct {
	mu        sync.RWMutex
	ports     []int
	updatedAt time.Time
	scanCount int64
}

// New returns an initialised, empty State.
func New() *State {
	return &State{}
}

// Update atomically replaces the current port list with ports.
func (s *State) Update(ports []int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := make([]int, len(ports))
	copy(copy, ports)
	s.ports = copy
	s.updatedAt = time.Now()
	s.scanCount++
}

// Ports returns a copy of the current open port list.
func (s *State) Ports() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.ports) == 0 {
		return nil
	}
	out := make([]int, len(s.ports))
	copy(out, s.ports)
	return out
}

// UpdatedAt returns the time of the last Update call.
func (s *State) UpdatedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updatedAt
}

// ScanCount returns the total number of updates applied.
func (s *State) ScanCount() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scanCount
}

// Contains reports whether port is in the current open port list.
func (s *State) Contains(port int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.ports {
		if p == port {
			return true
		}
	}
	return false
}
