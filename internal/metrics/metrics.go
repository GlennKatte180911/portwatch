// Package metrics tracks runtime statistics for portwatch scans,
// including scan counts, alert counts, and timing information.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	TotalScans   int64         `json:"total_scans"`
	TotalAlerts  int64         `json:"total_alerts"`
	LastScanAt   time.Time     `json:"last_scan_at"`
	LastAlertAt  time.Time     `json:"last_alert_at"`
	AvgScanTime  time.Duration `json:"avg_scan_time_ns"`
	UptimeSeconds float64      `json:"uptime_seconds"`
}

// Collector accumulates metrics over the lifetime of the process.
type Collector struct {
	mu          sync.Mutex
	startedAt   time.Time
	totalScans  int64
	totalAlerts int64
	lastScanAt  time.Time
	lastAlertAt time.Time
	scanTimeSum time.Duration
}

// New returns a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{startedAt: time.Now()}
}

// RecordScan records a completed scan and its duration.
func (c *Collector) RecordScan(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalScans++
	c.lastScanAt = time.Now()
	c.scanTimeSum += d
}

// RecordAlert records that an alert was emitted.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalAlerts++
	c.lastAlertAt = time.Now()
}

// Snapshot returns a consistent point-in-time copy of the metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	var avg time.Duration
	if c.totalScans > 0 {
		avg = c.scanTimeSum / time.Duration(c.totalScans)
	}

	return Snapshot{
		TotalScans:    c.totalScans,
		TotalAlerts:   c.totalAlerts,
		LastScanAt:    c.lastScanAt,
		LastAlertAt:   c.lastAlertAt,
		AvgScanTime:   avg,
		UptimeSeconds: time.Since(c.startedAt).Seconds(),
	}
}

// Reset zeroes all counters while preserving the start time.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalScans = 0
	c.totalAlerts = 0
	c.lastScanAt = time.Time{}
	c.lastAlertAt = time.Time{}
	c.scanTimeSum = 0
}
