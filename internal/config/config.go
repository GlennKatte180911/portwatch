package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// PortRange defines the range of ports to scan (e.g. "1-1024").
	PortRange string `json:"port_range"`

	// ScanInterval is how often to run a scan.
	ScanInterval time.Duration `json:"scan_interval"`

	// SnapshotPath is the file path used to persist port snapshots.
	SnapshotPath string `json:"snapshot_path"`

	// AlertOnNew triggers an alert when new open ports are detected.
	AlertOnNew bool `json:"alert_on_new"`

	// AlertOnClosed triggers an alert when previously open ports are closed.
	AlertOnClosed bool `json:"alert_on_closed"`

	// Timeout is the per-port connection timeout during scanning.
	Timeout time.Duration `json:"timeout"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		PortRange:     "1-1024",
		ScanInterval:  60 * time.Second,
		SnapshotPath:  ".portwatch_snapshot.json",
		AlertOnNew:    true,
		AlertOnClosed: false,
		Timeout:       500 * time.Millisecond,
	}
}

// Load reads a JSON config file from path and returns a Config.
// Missing fields fall back to Default values.
func Load(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the Config as JSON to path.
func (c *Config) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}
