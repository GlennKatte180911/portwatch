package healthcheck

import "fmt"

// Config holds configuration for the health-check HTTP server.
type Config struct {
	// Enabled controls whether the health-check server is started.
	Enabled bool `json:"enabled"`

	// Addr is the TCP address the server listens on, e.g. ":9090".
	Addr string `json:"addr"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		Addr:    ":9090",
	}
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("healthcheck: addr must not be empty")
	}
	return nil
}

// Build validates the config and returns a ready-to-use Checker.
// The caller is responsible for starting the HTTP server via ListenAndServe.
func (c Config) Build() (*Checker, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(), nil
}
