// Package portevict provides a grace-period eviction tracker for ports.
//
// When a port disappears from a scan it may be a transient blip rather than
// a genuine closure. Evictor holds such ports in a pending state for a
// configurable grace period; only once the period has elapsed without the
// port reappearing is it considered truly gone and returned via Confirmed.
//
// Typical usage:
//
//	evict := portevict.New(30 * time.Second)
//
//	// on each scan, mark newly-absent ports:
//	for _, p := range removedPorts {
//		evict.Mark(p)
//	}
//
//	// lift ports that came back:
//	for _, p := range reappearedPorts {
//		evict.Lift(p)
//	}
//
//	// act only on ports confirmed gone:
//	for _, p := range evict.Confirmed() {
//		alert(p)
//	}
package portevict
