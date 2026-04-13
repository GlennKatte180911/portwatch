// Package filter provides port filtering utilities for portwatch.
// It allows users to exclude specific ports or port ranges from scan results.
package filter

import "fmt"

// Rule represents a single port filter rule.
type Rule struct {
	Low  int
	High int
}

// Filter holds a set of exclusion rules applied to scan results.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of range strings like "22", "8000-9000".
func New(ranges []string) (*Filter, error) {
	f := &Filter{}
	for _, r := range ranges {
		rule, err := parseRange(r)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid range %q: %w", r, err)
		}
		f.rules = append(f.rules, rule)
	}
	return f, nil
}

// Allowed returns true if the port is not excluded by any rule.
func (f *Filter) Allowed(port int) bool {
	for _, r := range f.rules {
		if port >= r.Low && port <= r.High {
			return false
		}
	}
	return true
}

// Apply filters a slice of ports, returning only those not excluded.
func (f *Filter) Apply(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if f.Allowed(p) {
			out = append(out, p)
		}
	}
	return out
}

func parseRange(s string) (Rule, error) {
	var low, high int
	n, err := fmt.Sscanf(s, "%d-%d", &low, &high)
	if err != nil || n != 2 {
		// Try single port
		n, err = fmt.Sscanf(s, "%d", &low)
		if err != nil || n != 1 {
			return Rule{}, fmt.Errorf("expected PORT or PORT-PORT")
		}
		high = low
	}
	if low < 1 || high > 65535 || low > high {
		return Rule{}, fmt.Errorf("port values must be 1-65535 and low <= high")
	}
	return Rule{Low: low, High: high}, nil
}
