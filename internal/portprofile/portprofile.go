// Package portprofile builds a composite profile for a port by combining
// label, rank, classification, and scope information into a single summary.
package portprofile

import (
	"fmt"
	"sync"
)

// Profile holds aggregated metadata for a single port.
type Profile struct {
	Port     int
	Label    string
	Rank     string
	Class    string
	Scope    string
	Notes    []string
}

// Profiler builds Port profiles from registered providers.
type Profiler struct {
	mu       sync.RWMutex
	labeler  func(int) string
	ranker   func(int) string
	classer  func(int) string
	scoper   func(int) string
}

// New returns a Profiler with no-op providers. Use the With* methods to
// attach real implementations.
func New() *Profiler {
	noop := func(int) string { return "" }
	return &Profiler{
		labeler: noop,
		ranker:  noop,
		classer: noop,
		scoper:  noop,
	}
}

// WithLabeler sets the label provider.
func (p *Profiler) WithLabeler(fn func(int) string) *Profiler {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.labeler = fn
	return p
}

// WithRanker sets the rank provider.
func (p *Profiler) WithRanker(fn func(int) string) *Profiler {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ranker = fn
	return p
}

// WithClasser sets the classification provider.
func (p *Profiler) WithClasser(fn func(int) string) *Profiler {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.classer = fn
	return p
}

// WithScoper sets the scope provider.
func (p *Profiler) WithScoper(fn func(int) string) *Profiler {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.scoper = fn
	return p
}

// Build constructs a Profile for the given port.
func (p *Profiler) Build(port int) Profile {
	p.mu.RLock()
	defer p.mu.RUnlock()

	pr := Profile{
		Port:  port,
		Label: p.labeler(port),
		Rank:  p.ranker(port),
		Class: p.classer(port),
		Scope: p.scoper(port),
	}
	pr.Notes = buildNotes(pr)
	return pr
}

// BuildAll constructs profiles for each port in the slice.
func (p *Profiler) BuildAll(ports []int) []Profile {
	out := make([]Profile, len(ports))
	for i, port := range ports {
		out[i] = p.Build(port)
	}
	return out
}

func buildNotes(pr Profile) []string {
	var notes []string
	if pr.Rank == "critical" {
		notes = append(notes, fmt.Sprintf("port %d is ranked critical", pr.Port))
	}
	if pr.Class == "system" && pr.Port > 1023 {
		notes = append(notes, "system-class port outside well-known range")
	}
	return notes
}
