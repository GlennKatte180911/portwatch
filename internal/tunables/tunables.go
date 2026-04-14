// Package tunables provides runtime-adjustable parameters for portwatch.
// Values can be updated without restarting the process and are safe for
// concurrent use.
package tunables

import (
	"sync"
	"time"
)

// Tunables holds adjustable runtime parameters.
type Tunables struct {
	mu           sync.RWMutex
	scanInterval time.Duration
	scanTimeout  time.Duration
	alertCooldown time.Duration
	maxAlerts    int
}

// Defaults returns a Tunables instance populated with sensible defaults.
func Defaults() *Tunables {
	return &Tunables{
		scanInterval:  30 * time.Second,
		scanTimeout:   2 * time.Second,
		alertCooldown: 5 * time.Minute,
		maxAlerts:     10,
	}
}

// ScanInterval returns the current scan interval.
func (t *Tunables) ScanInterval() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.scanInterval
}

// SetScanInterval updates the scan interval. Returns an error if d <= 0.
func (t *Tunables) SetScanInterval(d time.Duration) error {
	if d <= 0 {
		return ErrInvalidValue{Field: "ScanInterval", Value: d.String()}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.scanInterval = d
	return nil
}

// ScanTimeout returns the current per-port scan timeout.
func (t *Tunables) ScanTimeout() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.scanTimeout
}

// SetScanTimeout updates the per-port scan timeout. Returns an error if d <= 0.
func (t *Tunables) SetScanTimeout(d time.Duration) error {
	if d <= 0 {
		return ErrInvalidValue{Field: "ScanTimeout", Value: d.String()}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.scanTimeout = d
	return nil
}

// AlertCooldown returns the minimum time between repeated alerts for the same port.
func (t *Tunables) AlertCooldown() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.alertCooldown
}

// SetAlertCooldown updates the alert cooldown. Returns an error if d < 0.
func (t *Tunables) SetAlertCooldown(d time.Duration) error {
	if d < 0 {
		return ErrInvalidValue{Field: "AlertCooldown", Value: d.String()}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.alertCooldown = d
	return nil
}

// MaxAlerts returns the maximum number of alerts emitted per scan cycle.
func (t *Tunables) MaxAlerts() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.maxAlerts
}

// SetMaxAlerts updates the max alerts cap. Returns an error if n < 1.
func (t *Tunables) SetMaxAlerts(n int) error {
	if n < 1 {
		return ErrInvalidValue{Field: "MaxAlerts", Value: fmt.Sprintf("%d", n)}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.maxAlerts = n
	return nil
}
