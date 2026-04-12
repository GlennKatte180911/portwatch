// Package history records and retrieves a persistent log of port change events
// detected by portwatch.
//
// Each time the monitor detects that ports have been added or removed, a new
// [Entry] is appended to the on-disk history file. Entries can later be
// retrieved in reverse-chronological order via [History.Last], making it easy
// to review recent activity without replaying the full log.
//
// Typical usage:
//
//	h := history.New("/var/lib/portwatch/history.json")
//	if err := h.Load(); err != nil {
//		log.Fatal(err)
//	}
//	// After a diff is computed:
//	if err := h.Record(diff.Added, diff.Removed); err != nil {
//		log.Println("history write error:", err)
//	}
package history
