package alert

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert describes a port change event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Ports     []int
}

// Notifier writes alerts to an output destination.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes an alert to the output.
func (n *Notifier) Notify(a Alert) {
	sort.Ints(a.Ports)
	fmt.Fprintf(n.out, "[%s] %s %s — ports: %v\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
		a.Ports,
	)
}

// NotifyDiff emits alerts for added and removed port sets.
func (n *Notifier) NotifyDiff(added, removed []int) {
	if len(added) > 0 {
		n.Notify(Alert{
			Timestamp: time.Now().UTC(),
			Level:     LevelAlert,
			Message:   "new open ports detected",
			Ports:     added,
		})
	}
	if len(removed) > 0 {
		n.Notify(Alert{
			Timestamp: time.Now().UTC(),
			Level:     LevelWarn,
			Message:   "ports no longer open",
			Ports:     removed,
		})
	}
	if len(added) == 0 && len(removed) == 0 {
		n.Notify(Alert{
			Timestamp: time.Now().UTC(),
			Level:     LevelInfo,
			Message:   "no port changes detected",
			Ports:     []int{},
		})
	}
}
