// Package portdiff provides utilities for computing and summarising
// differences between two sets of ports.
package portdiff

import "sort"

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	Added   []int
	Removed []int
}

// IsEmpty reports whether the diff contains no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0
}

// Compute returns the ports added and removed when moving from prev to curr.
func Compute(prev, curr []int) Diff {
	prevSet := toSet(prev)
	currSet := toSet(curr)

	var added, removed []int

	for p := range currSet {
		if !prevSet[p] {
			added = append(added, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			removed = append(removed, p)
		}
	}

	sort.Ints(added)
	sort.Ints(removed)

	return Diff{Added: added, Removed: removed}
}

// Summary returns a human-readable one-line description of the diff.
func Summary(d Diff) string {
	if d.IsEmpty() {
		return "no changes"
	}
	msg := ""
	if len(d.Added) > 0 {
		msg += fmt.Sprintf("+%d added", len(d.Added))
	}
	if len(d.Removed) > 0 {
		if msg != "" {
			msg += ", "
		}
		msg += fmt.Sprintf("-%d removed", len(d.Removed))
	}
	return msg
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
