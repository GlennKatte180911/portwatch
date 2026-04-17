// Package portwatch wires together the core scan-diff-alert pipeline.
package portwatch

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/portdiff"
	"github.com/user/portwatch/internal/portstate"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/notifier"
)

// Watcher runs periodic port scans and emits diffs via a Notifier.
type Watcher struct {
	scanner  *scanner.Scanner
	state    *portstate.State
	notifier notifier.Notifier
	interval time.Duration
}

// Config holds Watcher configuration.
type Config struct {
	StartPort int
	EndPort   int
	Interval  time.Duration
	Notifier  notifier.Notifier
}

// New creates a Watcher from the given Config.
func New(cfg Config) (*Watcher, error) {
	s, err := scanner.New(cfg.StartPort, cfg.EndPort)
	if err != nil {
		return nil, err
	}
	return &Watcher{
		scanner:  s,
		state:    portstate.New(),
		notifier: cfg.Notifier,
		interval: cfg.Interval,
	}, nil
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *Watcher) tick(ctx context.Context) error {
	ports, err := w.scanner.Scan(ctx)
	if err != nil {
		return err
	}
	prev := w.state.Ports()
	w.state.Update(ports)
	diff := portdiff.Compute(prev, ports)
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return nil
	}
	return w.notifier.Notify(ctx, diff)
}
