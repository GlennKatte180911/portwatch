// Package portaudit provides a structured audit trail for port state changes,
// recording who triggered a scan, what changed, and when.
package portaudit

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit record for a port change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Trigger   string    `json:"trigger"`
	Added     []int     `json:"added"`
	Removed   []int     `json:"removed"`
}

// Auditor records and persists audit entries to a file.
type Auditor struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// New creates a new Auditor that persists entries to the given file path.
// Existing entries are loaded from disk if the file is present.
func New(path string) (*Auditor, error) {
	a := &Auditor{path: path}
	if err := a.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return a, nil
}

// Record appends an audit entry for the given trigger and port diff.
// Entries with no added or removed ports are silently ignored.
func (a *Auditor) Record(trigger string, added, removed []int) error {
	if len(added) == 0 && len(removed) == 0 {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, Entry{
		Timestamp: time.Now().UTC(),
		Trigger:   trigger,
		Added:     added,
		Removed:   removed,
	})
	return a.save()
}

// Entries returns a copy of all recorded audit entries.
func (a *Auditor) Entries() []Entry {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]Entry, len(a.entries))
	copy(out, a.entries)
	return out
}

func (a *Auditor) load() error {
	data, err := os.ReadFile(a.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &a.entries)
}

func (a *Auditor) save() error {
	data, err := json.MarshalIndent(a.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.path, data, 0o644)
}
