// Package labelmap maintains a thread-safe registry that maps port numbers to
// human-readable service labels (e.g. 443 → "https", 5432 → "postgres").
//
// # Usage
//
// Create a registry pre-seeded with common well-known ports:
//
//	 lm := labelmap.New()
//
// Look up a label, falling back to "port/<n>" for unknown ports:
//
//	 name := lm.Label(port) // e.g. "ssh", "http", "port/9000"
//
// Extend the registry at runtime or from a JSON file:
//
//	 lm.Set(9000, "myapp")
//	 lm.LoadFile("/etc/portwatch/labels.json")
//
// The JSON file format is a simple object mapping port numbers to labels:
//
//	 {"9000": "myapp", "9001": "myapp-metrics"}
//
// All methods are safe for concurrent use.
package labelmap
