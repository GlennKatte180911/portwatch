// Package labelmap provides a registry for attaching human-readable labels
// (e.g. service names) to well-known port numbers, enriching alert output.
package labelmap

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// LabelMap maps port numbers to descriptive service labels.
type LabelMap struct {
	mu     sync.RWMutex
	entries map[int]string
}

// New returns a LabelMap pre-seeded with common well-known ports.
func New() *LabelMap {
	return &LabelMap{
		entries: map[int]string{
			22:   "ssh",
			25:   "smtp",
			53:   "dns",
			80:   "http",
			443:  "https",
			3306: "mysql",
			5432: "postgres",
			6379: "redis",
			8080: "http-alt",
			27017: "mongodb",
		},
	}
}

// Set registers a label for the given port, overwriting any existing entry.
func (l *LabelMap) Set(port int, label string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries[port] = label
}

// Get returns the label for port and whether one was found.
func (l *LabelMap) Get(port int) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	v, ok := l.entries[port]
	return v, ok
}

// Label returns the label for port, or a formatted fallback "port/<n>".
func (l *LabelMap) Label(port int) string {
	if name, ok := l.Get(port); ok {
		return name
	}
	return fmt.Sprintf("port/%d", port)
}

// LoadFile merges labels from a JSON file (map of "port": "label" pairs).
// Existing entries are overwritten where keys overlap.
func (l *LabelMap) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("labelmap: open %s: %w", path, err)
	}
	defer f.Close()

	var raw map[int]string
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return fmt.Errorf("labelmap: decode %s: %w", path, err)
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	for port, label := range raw {
		l.entries[port] = label
	}
	return nil
}

// Snapshot returns a copy of all current entries.
func (l *LabelMap) Snapshot() map[int]string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[int]string, len(l.entries))
	for k, v := range l.entries {
		out[k] = v
	}
	return out
}
