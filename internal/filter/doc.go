// Package filter provides port filtering functionality for portwatch.
//
// A Filter allows callers to define exclusion rules based on port ranges,
// so that well-known or intentionally open ports can be suppressed from
// change alerts during monitoring.
//
// # Creating a Filter
//
// Filters are constructed from a slice of range strings. Each string may be
// either a single port number (e.g. "22") or a hyphen-separated inclusive
// range (e.g. "8000-9000").
//
//	f, err := filter.New([]string{"22", "80", "8000-9000"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Checking Individual Ports
//
// Use Allowed to test whether a single port passes through the filter:
//
//	if f.Allowed(8080) {
//	    fmt.Println("port 8080 is not excluded")
//	}
//
// # Filtering a Slice
//
// Apply returns a new slice containing only the ports not matched by any
// exclusion rule:
//
//	visible := f.Apply([]int{22, 443, 8080, 8443})
//	// visible == []int{443, 8443}  (assuming 22 and 8080-range are excluded)
package filter
