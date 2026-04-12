// Package history provides persistent storage and retrieval of port scan history.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a single historical record of a port scan diff event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Added     []int     `json:"added"`
	Removed   []int     `json:"removed"`
}

// History holds a collection of scan diff entries.
type History struct {
	Entries []Entry `json:"entries"`
	path    string
}

// New creates a new History backed by the given file path.
func New(path string) *History {
	return &History{path: path}
}

// Record appends a new entry to the history and persists it to disk.
func (h *History) Record(added, removed []int) error {
	if len(added) == 0 && len(removed) == 0 {
		return nil
	}
	h.Entries = append(h.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Added:     added,
		Removed:   removed,
	})
	return h.save()
}

// Load reads history from disk into h. If the file does not exist, h.Entries
// is left empty and no error is returned.
func (h *History) Load() error {
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("history: read %s: %w", h.path, err)
	}
	return json.Unmarshal(data, h)
}

// Last returns up to n most recent entries, newest first.
func (h *History) Last(n int) []Entry {
	entries := make([]Entry, len(h.Entries))
	copy(entries, h.Entries)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})
	if n > len(entries) {
		n = len(entries)
	}
	return entries[:n]
}

func (h *History) save() error {
	if err := os.MkdirAll(filepath.Dir(h.path), 0o755); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	if err := os.WriteFile(h.path, data, 0o644); err != nil {
		return fmt.Errorf("history: write %s: %w", h.path, err)
	}
	return nil
}
