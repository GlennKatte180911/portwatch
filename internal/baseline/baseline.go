// Package baseline manages the trusted port baseline used to detect
// unexpected changes during monitoring sessions.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ErrNoBaseline is returned when no baseline file exists at the given path.
var ErrNoBaseline = errors.New("baseline: no baseline file found")

// Baseline represents a trusted snapshot of open ports at a point in time.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	Ports     []int     `json:"ports"`
	Path      string    `json:"-"`
}

// New creates a new Baseline from the provided ports and persists it to path.
func New(path string, ports []int) (*Baseline, error) {
	b := &Baseline{
		CreatedAt: time.Now().UTC(),
		Ports:     ports,
		Path:      path,
	}
	if err := b.Save(); err != nil {
		return nil, err
	}
	return b, nil
}

// Load reads a baseline from the given file path.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoBaseline
		}
		return nil, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	b.Path = path
	return &b, nil
}

// Save persists the baseline to its configured path.
func (b *Baseline) Save() error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.Path, data, 0o644)
}

// Contains reports whether port p is part of the baseline.
func (b *Baseline) Contains(p int) bool {
	for _, port := range b.Ports {
		if port == p {
			return true
		}
	}
	return false
}

// Unexpected returns ports from current that are not in the baseline.
func (b *Baseline) Unexpected(current []int) []int {
	var out []int
	for _, p := range current {
		if !b.Contains(p) {
			out = append(out, p)
		}
	}
	return out
}
