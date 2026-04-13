// Package reporter formats and outputs port change reports
// to various destinations such as stdout or a log file.
package reporter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Format controls how reports are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes port change reports to a writer.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter that writes to out using the given format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Report writes a human- or machine-readable summary of a Diff.
func (r *Reporter) Report(d snapshot.Diff) error {
	if len(d.Added) == 0 && len(d.Removed) == 0 {
		return nil
	}
	switch r.format {
	case FormatJSON:
		return r.writeJSON(d)
	default:
		return r.writeText(d)
	}
}

func (r *Reporter) writeText(d snapshot.Diff) error {
	ts := time.Now().Format(time.RFC3339)
	for _, p := range d.Added {
		if _, err := fmt.Fprintf(r.out, "[%s] OPENED  port %d\n", ts, p); err != nil {
			return err
		}
	}
	for _, p := range d.Removed {
		if _, err := fmt.Fprintf(r.out, "[%s] CLOSED  port %d\n", ts, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeJSON(d snapshot.Diff) error {
	ts := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(r.out,
		`{"timestamp":%q,"added":%s,"removed":%s}\n`,
		ts, intSliceJSON(d.Added), intSliceJSON(d.Removed))
	return err
}

func intSliceJSON(ports []int) string {
	if len(ports) == 0 {
		return "[]"
	}
	out := "["
	for i, p := range ports {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%d", p)
	}
	return out + "]"
}
