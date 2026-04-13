// Package circuitbreaker provides a thread-safe circuit breaker for portwatch.
//
// The circuit breaker wraps outbound calls (such as webhook notifications) and
// prevents cascading failures by short-circuiting requests when a downstream
// service becomes unhealthy.
//
// States:
//
//	- Closed:    Normal operation; all calls are allowed through.
//	- Open:      Failure threshold exceeded; calls are rejected with ErrOpen
//	             until the reset timeout elapses.
//	- Half-Open: Recovery probe; one call is allowed to test the downstream.
//	             A success transitions back to Closed; a failure re-opens.
//
// Usage:
//
//	b := circuitbreaker.New(5, 30*time.Second)
//	if err := b.Allow(); err != nil {
//	    // circuit is open — skip the call
//	    return err
//	}
//	if err := doCall(); err != nil {
//	    b.RecordFailure()
//	    return err
//	}
//	b.RecordSuccess()
package circuitbreaker
