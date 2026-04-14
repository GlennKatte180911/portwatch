// Package ratelimit provides per-port event rate limiting for portwatch.
//
// A Limiter tracks the last time an event was emitted for a given port and
// suppresses subsequent events that arrive within a configurable cooldown
// window. This prevents alert storms when a port flaps rapidly.
//
// Basic usage:
//
//	rl := ratelimit.New(5 * time.Second)
//	if rl.Allow(8080) {
//		// emit alert
//	}
//
// The Reset method clears the timestamp for a specific port, allowing the
// next event to pass immediately regardless of the cooldown window.
package ratelimit
