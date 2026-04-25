// Package portage provides age tracking for open ports.
//
// A Tracker records the first and last time each port was observed as open,
// allowing operators to identify long-lived or unexpectedly persistent
// connections.
//
// Basic usage:
//
//	tr := portage.New(time.Now)
//
//	// After each scan, pass the list of open ports:
//	tr.Observe(openPorts)
//
//	// When a port closes, remove it so age resets if it reopens:
//	tr.Remove(closedPort)
//
//	// Find ports open longer than 1 hour:
//	stale := tr.OlderThan(time.Hour)
//
The zero value is not usable; always construct with New.
package portage
