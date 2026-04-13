// Package schedule provides configurable scan interval management
// for portwatch, allowing scans to be triggered on a fixed cadence
// or with adaptive jitter to avoid thundering-herd patterns.
package schedule

import (
	"context"
	"math/rand"
	"time"
)

// Scheduler emits ticks on C at the configured interval.
type Scheduler struct {
	C        <-chan time.Time
	Interval time.Duration
	Jitter   time.Duration
	stop     chan struct{}
}

// New creates a Scheduler that fires every interval, optionally adding
// up to jitter of random delay to each period. Pass 0 for jitter to
// disable it.
func New(interval, jitter time.Duration) *Scheduler {
	ch := make(chan time.Time, 1)
	s := &Scheduler{
		C:        ch,
		Interval: interval,
		Jitter:   jitter,
		stop:     make(chan struct{}),
	}
	go s.run(ch)
	return s
}

func (s *Scheduler) run(ch chan<- time.Time) {
	for {
		wait := s.Interval
		if s.Jitter > 0 {
			//nolint:gosec // non-crypto jitter is fine
			wait += time.Duration(rand.Int63n(int64(s.Jitter)))
		}
		select {
		case <-time.After(wait):
			select {
			case ch <- time.Now():
			default:
			}
		case <-s.stop:
			return
		}
	}
}

// Stop shuts down the scheduler. Safe to call multiple times.
func (s *Scheduler) Stop() {
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
}

// RunWithContext blocks, calling fn each time the scheduler ticks,
// until ctx is cancelled or Stop is called.
func (s *Scheduler) RunWithContext(ctx context.Context, fn func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stop:
			return
		case <-s.C:
			fn()
		}
	}
}
