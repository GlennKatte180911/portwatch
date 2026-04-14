// Package portlabel resolves human-readable service names for port numbers,
// combining the built-in well-known port list with user-supplied overrides.
package portlabel

import (
	"fmt"
	"sync"
)

// Resolver maps port numbers to service label strings.
type Resolver struct {
	mu     sync.RWMutex
	labels map[int]string
}

// New returns a Resolver pre-loaded with common well-known port names.
func New() *Resolver {
	r := &Resolver{labels: make(map[int]string)}
	for port, name := range wellKnown {
		r.labels[port] = name
	}
	return r
}

// Set registers a custom label for the given port, overwriting any existing entry.
func (r *Resolver) Set(port int, label string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.labels[port] = label
}

// Label returns the service name for port. If no entry exists the fallback
// string "port/<n>" is returned so callers always receive a non-empty value.
func (r *Resolver) Label(port int) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if name, ok := r.labels[port]; ok {
		return name
	}
	return fmt.Sprintf("port/%d", port)
}

// Annotate returns a copy of the provided port slice as a slice of
// PortLabel pairs, each carrying the resolved service name.
func (r *Resolver) Annotate(ports []int) []PortLabel {
	out := make([]PortLabel, len(ports))
	for i, p := range ports {
		out[i] = PortLabel{Port: p, Label: r.Label(p)}
	}
	return out
}

// PortLabel pairs a port number with its resolved human-readable name.
type PortLabel struct {
	Port  int    `json:"port"`
	Label string `json:"label"`
}
