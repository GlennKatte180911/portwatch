// Package portsampler provides statistical sampling of port scan results,
// allowing portwatch to reduce noise by only reporting ports that appear
// consistently across multiple scans within a rolling window.
package portsampler

import (
	"sync"
	"time"
)

// Sample holds the observation count and timing metadata for a single port.
type Sample struct {
	Port      int
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Sampler tracks how frequently each port appears across scans.
type Sampler struct {
	mu        sync.Mutex
	samples   map[int]*Sample
	window    time.Duration
	threshold int
}

// New creates a Sampler that retains observations within the given window
// and considers a port "stable" once it has been seen at least threshold times.
func New(window time.Duration, threshold int) *Sampler {
	if threshold < 1 {
		threshold = 1
	}
	return &Sampler{
		samples:   make(map[int]*Sample),
		window:    window,
		threshold: threshold,
	}
}

// Observe records that the given ports were seen at the current time.
func (s *Sampler) Observe(ports []int) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict(now)
	for _, p := range ports {
		if sm, ok := s.samples[p]; ok {
			sm.Count++
			sm.LastSeen = now
		} else {
			s.samples[p] = &Sample{Port: p, Count: 1, FirstSeen: now, LastSeen: now}
		}
	}
}

// Stable returns all ports whose observation count meets or exceeds the threshold.
func (s *Sampler) Stable() []int {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict(now)
	out := make([]int, 0, len(s.samples))
	for p, sm := range s.samples {
		if sm.Count >= s.threshold {
			out = append(out, p)
		}
	}
	return out
}

// Get returns the Sample for a port and whether it exists.
func (s *Sampler) Get(port int) (Sample, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sm, ok := s.samples[port]
	if !ok {
		return Sample{}, false
	}
	return *sm, true
}

// Reset clears all recorded samples.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples = make(map[int]*Sample)
}

// evict removes samples whose last observation has fallen outside the window.
// Caller must hold s.mu.
func (s *Sampler) evict(now time.Time) {
	cutoff := now.Add(-s.window)
	for p, sm := range s.samples {
		if sm.LastSeen.Before(cutoff) {
			delete(s.samples, p)
		}
	}
}
