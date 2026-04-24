// Package portfence provides a boundary enforcement layer that blocks
// or flags ports that fall outside a declared acceptable set.
package portfence

import (
	"fmt"
	"sort"
	"sync"
)

// Verdict describes the outcome of a fence check.
type Verdict int

const (
	Allow Verdict = iota
	Deny
)

func (v Verdict) String() string {
	if v == Allow {
		return "allow"
	}
	return "deny"
}

// Violation records a port that violated the fence.
type Violation struct {
	Port    int
	Verdict Verdict
}

// Fence holds the set of explicitly permitted ports.
type Fence struct {
	mu      sync.RWMutex
	allowed map[int]struct{}
}

// New creates a Fence from the provided list of permitted ports.
func New(permitted []int) (*Fence, error) {
	allowed := make(map[int]struct{}, len(permitted))
	for _, p := range permitted {
		if p < 1 || p > 65535 {
			return nil, fmt.Errorf("portfence: invalid port %d", p)
		}
		allowed[p] = struct{}{}
	}
	return &Fence{allowed: allowed}, nil
}

// Permit adds a port to the allowed set.
func (f *Fence) Permit(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portfence: invalid port %d", port)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowed[port] = struct{}{}
	return nil
}

// Revoke removes a port from the allowed set.
func (f *Fence) Revoke(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.allowed, port)
}

// Check returns Allow if the port is permitted, Deny otherwise.
func (f *Fence) Check(port int) Verdict {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if _, ok := f.allowed[port]; ok {
		return Allow
	}
	return Deny
}

// Evaluate checks each port and returns a Violation for every denied port.
func (f *Fence) Evaluate(ports []int) []Violation {
	var violations []Violation
	for _, p := range ports {
		if f.Check(p) == Deny {
			violations = append(violations, Violation{Port: p, Verdict: Deny})
		}
	}
	return violations
}

// Permitted returns a sorted slice of all currently allowed ports.
func (f *Fence) Permitted() []int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]int, 0, len(f.allowed))
	for p := range f.allowed {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}
