// Package eventlog provides a structured, append-only event log for recording
// port change events detected by portwatch.
//
// Each entry captures the timestamp and the sets of ports that were added or
// removed during a single scan cycle. Entries are stored as newline-delimited
// JSON (JSONL) so the file remains human-readable and can be processed by
// standard Unix tools.
//
// Basic usage:
//
//	log := eventlog.New("/var/lib/portwatch/events.jsonl")
//
//	// Record a change.
//	if err := log.Append(added, removed); err != nil {
//		log.Printf("eventlog: %v", err)
//	}
//
//	// Retrieve the ten most recent events.
//	recent, err := log.Latest(10)
//
//	// Query events since a given time.
//	events, err := log.Query(eventlog.QueryOptions{
//		Since: time.Now().Add(-24 * time.Hour),
//		Limit: 50,
//	})
package eventlog
