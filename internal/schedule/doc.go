// Package schedule provides a Scheduler type for portwatch that emits
// periodic ticks to drive port-scan cycles.
//
// Basic usage:
//
//	s := schedule.New(30*time.Second, 2*time.Second)
//	defer s.Stop()
//	s.RunWithContext(ctx, func() {
//		// perform a port scan
//	})
//
// Interval is the base period between scans. Jitter adds a random
// duration in [0, jitter) to each interval, which is useful when
// multiple portwatch instances run in the same environment and you
// want to spread their load.
//
// The scheduler is safe for concurrent use. Stop may be called from
// any goroutine and is idempotent.
package schedule
