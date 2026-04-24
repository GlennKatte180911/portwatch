// Package portquota enforces upper bounds on the number of open ports
// within named groups (e.g. labels or scopes).
//
// # Overview
//
// Create a Quota, register per-group limits with Set, then call Check
// with a map of current group → port-count values. Check returns a
// Violation for every group whose count exceeds its registered ceiling.
//
// # Example
//
//	q := portquota.New()
//	_ = q.Set("web", 5)
//	_ = q.Set("db", 2)
//
//	counts := map[string]int{"web": 8, "db": 1}
//	for _, v := range q.Check(counts) {
//		fmt.Println(v) // quota exceeded for group "web": limit 5, actual 8
//	}
//
// Quota is safe for concurrent use.
package portquota
