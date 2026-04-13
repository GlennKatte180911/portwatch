// Package baseline provides functionality for recording and comparing a
// trusted set of open ports (the "baseline") against the current state of
// the host.
//
// A baseline is established once — typically at first run or on explicit user
// request — and then persisted to disk as a JSON file. Subsequent monitoring
// cycles load the baseline and compare it against freshly scanned ports to
// surface only the ports that are genuinely new or unexpected.
//
// Usage:
//
//	// Create and persist a new baseline
//	b, err := baseline.New("/var/lib/portwatch/baseline.json", currentPorts)
//
//	// Load an existing baseline
//	b, err := baseline.Load("/var/lib/portwatch/baseline.json")
//
//	// Find ports not present in the baseline
//	unexpected := b.Unexpected(currentPorts)
package baseline
