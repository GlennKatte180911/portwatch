// Package portcache provides a thread-safe, TTL-based in-memory cache for
// port scan results.
//
// # Overview
//
// Scanning the full port range on every tick can be expensive. portcache sits
// between the scheduler and the scanner: if a valid (non-expired) result
// exists for a given host key it is returned immediately, avoiding a
// redundant network round-trip.
//
// # Usage
//
//	c := portcache.New(30 * time.Second)
//
//	// After a scan completes, store the result.
//	c.Set("127.0.0.1", openPorts)
//
//	// On the next tick, check the cache first.
//	if entry, ok := c.Get("127.0.0.1"); ok {
//		// reuse entry.Ports
//	}
//
//	// Periodically evict stale entries.
//	c.Purge()
//
// A TTL of zero disables caching: every Get call returns a miss.
package portcache
