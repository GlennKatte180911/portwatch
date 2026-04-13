// Package aggregator provides a time-windowed batching layer for port-change
// events.
//
// During a busy period — for example when a service restarts and briefly
// closes then re-opens several ports — the aggregator collects all
// snapshot.Diff values pushed within a configurable window and emits a single
// consolidated Event when the window expires.  Ports that are both added and
// removed within the same window cancel each other out and are not reported.
//
// Typical usage:
//
//	agg := aggregator.New(5 * time.Second)
//	defer agg.Stop()
//
//	// in the monitor loop:
//	agg.Push(diff)
//
//	// in a consumer goroutine:
//	for ev := range agg.Events() {
//	    fmt.Println("added:", ev.Added, "removed:", ev.Removed)
//	}
package aggregator
