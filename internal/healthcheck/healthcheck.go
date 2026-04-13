// Package healthcheck provides a simple liveness probe for portwatch,
// exposing an HTTP endpoint that reports the monitor's operational status.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Status represents the current health of the monitor.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	Scans     int64     `json:"scans_completed"`
	LastScan  time.Time `json:"last_scan"`
	StartedAt time.Time `json:"started_at"`
}

// Checker tracks runtime metrics and serves a health endpoint.
type Checker struct {
	startedAt time.Time
	lastScan  atomic.Value // stores time.Time
	scans     atomic.Int64
}

// New creates a new Checker with the start time set to now.
func New() *Checker {
	c := &Checker{startedAt: time.Now()}
	c.lastScan.Store(time.Time{})
	return c
}

// RecordScan updates the last-scan timestamp and increments the scan counter.
func (c *Checker) RecordScan() {
	c.lastScan.Store(time.Now())
	c.scans.Add(1)
}

// Status returns a snapshot of the current health metrics.
func (c *Checker) Status() Status {
	return Status{
		OK:        true,
		Uptime:    time.Since(c.startedAt).Round(time.Second).String(),
		Scans:     c.scans.Load(),
		LastScan:  c.lastScan.Load().(time.Time),
		StartedAt: c.startedAt,
	}
}

// Handler returns an http.HandlerFunc that writes the health status as JSON.
func (c *Checker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(c.Status())
	}
}

// ListenAndServe starts an HTTP server on the given address serving /healthz.
// It blocks until the server exits.
func (c *Checker) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/healthz", c.Handler())
	server := &http.Server{
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}
	return fmt.Errorf("healthcheck server: %w", server.ListenAndServe())
}
