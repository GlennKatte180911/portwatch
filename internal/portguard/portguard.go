// Package portguard enforces per-port access guards, blocking or allowing
// scanner activity based on a configurable allowlist and denylist. It is
// intended as a lightweight policy gate sitting between the scanner and the
// alert pipeline.
package portguard

import (
	"errors"
	"fmt"
	"sync"
)

// Action describes the outcome of a guard evaluation.
type Action int

const (
	// Allow means the port passes the guard.
	Allow Action = iota
	// Deny means the port is blocked by the guard.
	Deny
)

// Guard holds allowlist and denylist sets and evaluates ports against them.
// Denylist takes precedence over allowlist. If both lists are empty every port
// is allowed.
type Guard struct {
	mu       sync.RWMutex
	allowSet map[int]struct{}
	denySet  map[int]struct{}
}

// New returns a new Guard with empty allowlist and denylist.
func New() *Guard {
	return &Guard{
		allowSet: make(map[int]struct{}),
		denySet:  make(map[int]struct{}),
	}
}

// Allow adds port to the allowlist.
func (g *Guard) Permit(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portguard: invalid port %d", port)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.allowSet[port] = struct{}{}
	return nil
}

// Block adds port to the denylist.
func (g *Guard) Block(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portguard: invalid port %d", port)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.denySet[port] = struct{}{}
	return nil
}

// Evaluate returns Allow or Deny for the given port.
func (g *Guard) Evaluate(port int) Action {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, denied := g.denySet[port]; denied {
		return Deny
	}
	if len(g.allowSet) == 0 {
		return Allow
	}
	if _, allowed := g.allowSet[port]; allowed {
		return Allow
	}
	return Deny
}

// Apply filters ports, returning only those that pass the guard.
func (g *Guard) Apply(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if g.Evaluate(p) == Allow {
			out = append(out, p)
		}
	}
	return out
}

// ErrInvalidPort is returned when a port number is out of the valid range.
var ErrInvalidPort = errors.New("portguard: port out of range")
