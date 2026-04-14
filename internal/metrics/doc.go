// Package metrics provides a lightweight, concurrency-safe collector for
// portwatch runtime statistics.
//
// # Overview
//
// A Collector accumulates scan and alert events during the lifetime of the
// portwatch process. Call RecordScan after every port sweep completes and
// RecordAlert whenever an alert is dispatched to a notifier.
//
// # Usage
//
//	c := metrics.New()
//
//	// after each scan:
//	c.RecordScan(elapsed)
//
//	// after each alert:
//	c.RecordAlert()
//
//	// read a consistent snapshot at any time:
//	snap := c.Snapshot()
//
// # HTTP endpoint
//
// The Collector exposes an http.HandlerFunc via Handler() that serialises
// the current Snapshot as JSON. Mount it alongside the healthcheck endpoint:
//
//	http.Handle("/metrics", collector.Handler())
package metrics
