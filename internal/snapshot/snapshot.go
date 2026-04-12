package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// New creates a new Snapshot with the current timestamp.
func New(ports []int) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
}

// Save writes the snapshot to the given file path as JSON.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Diff compares two snapshots and returns added and removed ports.
func Diff(previous, current *Snapshot) (added, removed []int) {
	prev := toSet(previous.Ports)
	curr := toSet(current.Ports)

	for p := range curr {
		if !prev[p] {
			added = append(added, p)
		}
	}
	for p := range prev {
		if !curr[p] {
			removed = append(removed, p)
		}
	}
	return added, removed
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
