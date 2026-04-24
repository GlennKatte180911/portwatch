// Package portquota enforces per-label or per-scope port count limits,
// alerting when the number of open ports in a group exceeds a configured ceiling.
package portquota

import (
	"fmt"
	"sync"
)

// Violation describes a quota breach for a named group.
type Violation struct {
	Group string
	Limit int
	Actual int
}

// String returns a human-readable description of the violation.
func (v Violation) String() string {
	return fmt.Sprintf("quota exceeded for group %q: limit %d, actual %d", v.Group, v.Limit, v.Actual)
}

// Quota holds per-group port count limits.
type Quota struct {
	mu     sync.RWMutex
	limits map[string]int
}

// New returns a Quota with no limits set.
func New() *Quota {
	return &Quota{limits: make(map[string]int)}
}

// Set registers a maximum port count for the named group.
// A limit of zero or below is rejected.
func (q *Quota) Set(group string, limit int) error {
	if limit <= 0 {
		return fmt.Errorf("portquota: limit for %q must be > 0, got %d", group, limit)
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.limits[group] = limit
	return nil
}

// Remove deletes the limit for the named group.
func (q *Quota) Remove(group string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.limits, group)
}

// Check evaluates counts (group → port count) against registered limits.
// It returns a Violation for every group whose count exceeds its limit.
func (q *Quota) Check(counts map[string]int) []Violation {
	q.mu.RLock()
	defer q.mu.RUnlock()

	var violations []Violation
	for group, limit := range q.limits {
		actual, ok := counts[group]
		if !ok {
			continue
		}
		if actual > limit {
			violations = append(violations, Violation{
				Group:  group,
				Limit:  limit,
				Actual: actual,
			})
		}
	}
	return violations
}

// Limits returns a snapshot of all configured group limits.
func (q *Quota) Limits() map[string]int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	out := make(map[string]int, len(q.limits))
	for k, v := range q.limits {
		out[k] = v
	}
	return out
}
