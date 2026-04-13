package aggregator

import (
	"errors"
	"time"
)

// Config holds tunable parameters for an Aggregator.
type Config struct {
	// Window is the duration over which diffs are accumulated before a
	// batched Event is emitted. Must be positive.
	Window time.Duration `json:"window"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 5 * time.Second,
	}
}

// Validate returns an error if any field is out of range.
func (c Config) Validate() error {
	if c.Window <= 0 {
		return errors.New("aggregator: window must be positive")
	}
	return nil
}

// Build constructs an Aggregator from the Config after validating it.
func (c Config) Build() (*Aggregator, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Window), nil
}
