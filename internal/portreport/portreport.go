// Package portreport provides a structured summary report of the current
// port landscape, combining profile, policy, and trend data into a single
// human-readable or machine-readable output.
package portreport

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

// Entry holds the report data for a single port.
type Entry struct {
	Port      int       `json:"port"`
	Label     string    `json:"label"`
	Class     string    `json:"class"`
	Rank      string    `json:"rank"`
	Policy    string    `json:"policy"`
	SeenCount int       `json:"seen_count"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Notes     []string  `json:"notes,omitempty"`
}

// Report is an ordered collection of port entries.
type Report struct {
	GeneratedAt time.Time `json:"generated_at"`
	Entries     []Entry   `json:"entries"`
}

// Builder assembles a Report from provided entries.
type Builder struct {
	entries []Entry
}

// New returns a new Builder.
func New() *Builder {
	return &Builder{}
}

// Add appends an entry to the report.
func (b *Builder) Add(e Entry) {
	b.entries = append(b.entries, e)
}

// Build sorts entries by port number and returns the final Report.
func (b *Builder) Build() Report {
	sorted := make([]Entry, len(b.entries))
	copy(sorted, b.entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Port < sorted[j].Port
	})
	return Report{
		GeneratedAt: time.Now().UTC(),
		Entries:     sorted,
	}
}

// WriteText writes a human-readable text summary to w.
func WriteText(w io.Writer, r Report) error {
	_, err := fmt.Fprintf(w, "Port Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	for _, e := range r.Entries {
		_, err = fmt.Fprintf(w, "  %-6d %-20s %-10s %-8s policy:%-6s seen:%d\n",
			e.Port, e.Label, e.Class, e.Rank, e.Policy, e.SeenCount)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes the report as JSON to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
