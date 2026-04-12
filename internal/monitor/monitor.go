// Package monitor orchestrates periodic port scanning and change detection.
package monitor

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Monitor periodically scans ports and emits alerts on changes.
type Monitor struct {
	scanner  *scanner.Scanner
	notifier *alert.Notifier
	history  *history.History
	interval time.Duration
	snapshotPath string
}

// New creates a Monitor with the provided components and scan interval.
func New(
	s *scanner.Scanner,
	n *alert.Notifier,
	h *history.History,
	interval time.Duration,
	snapshotPath string,
) *Monitor {
	return &Monitor{
		scanner:      s,
		notifier:     n,
		history:      h,
		interval:     interval,
		snapshotPath: snapshotPath,
	}
}

// Run starts the monitoring loop. It blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.tick()
		}unc (m *Monitor) tick {
	ports, err := m.scanner.("monitor err)
		return
	}

	current := snapshot.New(ports)

	prev, err := snapshot.Load(m.snapshotPath)
	if err != nil {
		log.Printf("monitor: load snapshot: %v", err)
		_ = current.Save(m.snapshotPath)
		return
	}

	diff := snapshot.Diff(prev, current)
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		_ = current.Save(m.snapshotPath)
		return
	}

	m.notifier.Notify(diff)

	if err := m.history.Record(diff.Added, diff.Removed); err != nil {
		log.Printf("monitor: history record error: %v", err)
	}

	_ = current.Save(m.snapshotPath)
}
