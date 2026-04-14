// Package portstate provides a thread-safe container for the current
// observed open-port state of the local machine.
//
// # Overview
//
// A State is updated by the monitor after each scan cycle and can be
// queried concurrently by other components (reporters, health checks,
// notifiers) without additional locking.
//
// # Usage
//
//	st := portstate.New()
//
//	// After each scan:
//	st.Update(openPorts)
//
//	// Query at any time:
//	ports := st.Ports()
//	if st.Contains(443) {
//		// HTTPS is open
//	}
//
package portstate
