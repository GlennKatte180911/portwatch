// Package suppress provides temporary port alert suppression for portwatch.
//
// During planned maintenance or known configuration changes, operators may
// want to silence alerts for specific ports for a defined period. Suppressor
// tracks per-port expiry times and filters ports from alert pipelines.
//
// Example usage:
//
//	s := suppress.New()
//
//	// Silence alerts on port 8080 for 30 minutes.
//	s.Suppress(8080, 30*time.Minute)
//
//	// Check before sending an alert.
//	if !s.IsSuppressed(port) {
//		notifier.Notify(event)
//	}
//
//	// Filter a slice of changed ports before processing.
//	active := s.Apply(changedPorts)
//
// Suppressions expire automatically; no background goroutine is required.
package suppress
