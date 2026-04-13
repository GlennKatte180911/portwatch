// Package throttle implements a sliding-window token-bucket throttle
// for controlling the rate of outbound notifications in portwatch.
//
// Unlike the ratelimit package, which enforces a minimum cooldown
// between successive events for the same port, throttle allows a
// configurable burst of up to N events within a rolling time window
// before silencing further events until the window advances.
//
// Typical usage:
//
//	th := throttle.New(time.Minute, 5)
//
//	if th.Allow("port:8080") {
//	    notifier.Notify(event)
//	}
//
// Each key is tracked independently, so bursts on one port do not
// affect the allowance of another.
package throttle
