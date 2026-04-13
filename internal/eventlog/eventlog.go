// Package eventlog provides a structured, append-only log of port change
// events that can be queried and streamed for audit or diagnostic purposes.
package eventlog

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded port-change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Added     []int     `json:"added"`
	Removed   []int     `json:"removed"`
}

// EventLog is a thread-safe, file-backed log of port change events.
type EventLog struct {
	mu   sync.Mutex
	path string
}

// New returns an EventLog that persists entries to path.
func New(path string) *EventLog {
	return &EventLog{path: path}
}

// Append writes a new entry to the log file.
func (l *EventLog) Append(added, removed []int) error {
	if len(added) == 0 && len(removed) == 0 {
		return nil
	}
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Added:     added,
		Removed:   removed,
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("eventlog: open %s: %w", l.path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("eventlog: encode entry: %w", err)
	}
	return nil
}

// Load reads all entries from the log file.
// If the file does not exist an empty slice is returned.
func (l *EventLog) Load() ([]Entry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("eventlog: open %s: %w", l.path, err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("eventlog: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
