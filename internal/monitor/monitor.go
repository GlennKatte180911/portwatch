package monitor

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Monitor orchestrates periodic port scanning and change detection.
type Monitor struct {
	cfg      *config.Config
	scanner  *scanner.Scanner
	notifier *alert.Notifier
}

// New creates a new Monitor with the given configuration.
func New(cfg *config.Config, notifier *alert.Notifier) *Monitor {
	return &Monitor{
		cfg:      cfg,
		scanner:  scanner.New(cfg.Timeout),
		notifier: notifier,
	}
}

// Run starts the monitoring loop. It blocks until the done channel is closed.
func (m *Monitor) Run(done <-chan struct{}) error {
	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	if err := m.tick(); err != nil {
		return fmt.Errorf("initial scan failed: %w", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := m.tick(); err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}
		case <-done:
			return nil
		}
	}
}

// tick performs a single scan cycle: scan ports, diff against previous snapshot, alert on changes.
func (m *Monitor) tick() error {
	ports, err := m.scanner.Scan(m.cfg.StartPort, m.cfg.EndPort)
	if err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	current := snapshot.New(ports)

	previous, err := snapshot.Load(m.cfg.SnapshotPath)
	if err == nil {
		added, removed := snapshot.Diff(previous, current)
		if len(added) > 0 || len(removed) > 0 {
			m.notifier.Notify(added, removed)
		}
	}

	if err := current.Save(m.cfg.SnapshotPath); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	return nil
}
