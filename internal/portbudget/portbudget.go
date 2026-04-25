// Package portbudget tracks the number of open ports against a configured
// budget and reports whether the current count is within acceptable bounds.
package portbudget

import (
	"errors"
	"fmt"
	"sync"
)

// ErrBudgetExceeded is returned when the observed port count exceeds the limit.
var ErrBudgetExceeded = errors.New("port budget exceeded")

// Violation describes a single budget breach.
type Violation struct {
	Limit   int
	Actual  int
	Excess  int
}

func (v Violation) Error() string {
	return fmt.Sprintf("port budget exceeded: limit %d, actual %d, excess %d",
		v.Limit, v.Actual, v.Excess)
}

// Budget enforces a maximum number of simultaneously open ports.
type Budget struct {
	mu    sync.RWMutex
	limit int
}

// New creates a Budget with the given limit. Returns an error if limit < 1.
func New(limit int) (*Budget, error) {
	if limit < 1 {
		return nil, errors.New("portbudget: limit must be at least 1")
	}
	return &Budget{limit: limit}, nil
}

// SetLimit updates the budget limit at runtime. Returns an error if limit < 1.
func (b *Budget) SetLimit(limit int) error {
	if limit < 1 {
		return errors.New("portbudget: limit must be at least 1")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.limit = limit
	return nil
}

// Limit returns the current budget limit.
func (b *Budget) Limit() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.limit
}

// Check evaluates whether the supplied port slice is within budget.
// It returns a Violation (which also satisfies the error interface) when the
// count exceeds the limit, and nil otherwise.
func (b *Budget) Check(ports []int) error {
	b.mu.RLock()
	limit := b.limit
	b.mu.RUnlock()

	actual := len(ports)
	if actual > limit {
		return Violation{
			Limit:  limit,
			Actual: actual,
			Excess: actual - limit,
		}
	}
	return nil
}

// Within returns true when the port count is at or below the budget limit.
func (b *Budget) Within(ports []int) bool {
	return b.Check(ports) == nil
}
