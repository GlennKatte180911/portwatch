package circuitbreaker

import (
	"errors"
	"time"
)

// Config holds tunable parameters for constructing a Breaker.
type Config struct {
	// MaxFailures is the number of consecutive failures before the circuit opens.
	MaxFailures int `json:"max_failures"`
	// ResetTimeout is how long to wait before attempting recovery.
	ResetTimeout time.Duration `json:"reset_timeout"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxFailures:  5,
		ResetTimeout: 30 * time.Second,
	}
}

// Validate returns an error if the Config contains invalid values.
func (c Config) Validate() error {
	if c.MaxFailures <= 0 {
		return errors.New("circuitbreaker: MaxFailures must be greater than zero")
	}
	if c.ResetTimeout <= 0 {
		return errors.New("circuitbreaker: ResetTimeout must be greater than zero")
	}
	return nil
}

// Build validates the Config and returns a new Breaker.
func (c Config) Build() (*Breaker, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.MaxFailures, c.ResetTimeout), nil
}
