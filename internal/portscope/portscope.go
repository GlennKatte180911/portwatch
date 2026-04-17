// Package portscope defines named scan scopes that restrict which port
// ranges are active during a given monitoring session.
package portscope

import (
	"fmt"
	"sort"
)

// Scope represents a named collection of port ranges.
type Scope struct {
	Name   string
	Ranges [][2]int // inclusive [low, high] pairs
}

// Registry holds a set of named scopes.
type Registry struct {
	scopes map[string]Scope
}

// New returns a Registry pre-loaded with sensible default scopes.
func New() *Registry {
	r := &Registry{scopes: make(map[string]Scope)}
	r.Register(Scope{Name: "system", Ranges: [][2]int{{1, 1023}}})
	r.Register(Scope{Name: "registered", Ranges: [][2]int{{1024, 49151}}})
	r.Register(Scope{Name: "dynamic", Ranges: [][2]int{{49152, 65535}}})
	r.Register(Scope{Name: "all", Ranges: [][2]int{{1, 65535}}})
	return r
}

// Register adds or replaces a scope in the registry.
func (r *Registry) Register(s Scope) error {
	if s.Name == "" {
		return fmt.Errorf("portscope: scope name must not be empty")
	}
	r.scopes[s.Name] = s
	return nil
}

// Get returns the scope with the given name, or an error if not found.
func (r *Registry) Get(name string) (Scope, error) {
	s, ok := r.scopes[name]
	if !ok {
		return Scope{}, fmt.Errorf("portscope: unknown scope %q", name)
	}
	return s, nil
}

// Contains reports whether port p falls within any range of the scope.
func (s Scope) Contains(p int) bool {
	for _, rng := range s.Ranges {
		if p >= rng[0] && p <= rng[1] {
			return true
		}
	}
	return false
}

// Names returns a sorted list of registered scope names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.scopes))
	for n := range r.scopes {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
