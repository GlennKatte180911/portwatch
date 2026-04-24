// Package portguard provides a thread-safe allowlist/denylist gate for port
// numbers.
//
// # Overview
//
// A Guard holds two sets of ports: an allowlist and a denylist. When
// Evaluate is called for a given port the following precedence rules apply:
//
//  1. If the port appears in the denylist it is always denied.
//  2. If the allowlist is non-empty the port must appear in it to be allowed.
//  3. If both lists are empty every port is allowed (open policy).
//
// # Usage
//
//	g := portguard.New()
//	_ = g.Permit(443)
//	_ = g.Permit(22)
//	_ = g.Block(23) // telnet — always deny
//
//	filtered := g.Apply(scannedPorts)
//
// Guard is safe for concurrent use.
package portguard
