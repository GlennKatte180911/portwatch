package portlabel

import (
	"errors"
	"fmt"
)

// Config holds configuration for building a Resolver with optional overrides.
type Config struct {
	// Overrides maps port numbers to custom service labels. These are applied
	// on top of the built-in well-known port list after construction.
	Overrides map[int]string `json:"overrides,omitempty"`
}

// DefaultConfig returns a Config with no overrides.
func DefaultConfig() Config {
	return Config{Overrides: make(map[int]string)}
}

// Validate checks that all override port numbers are within the valid range.
func (c Config) Validate() error {
	for port := range c.Overrides {
		if port < 1 || port > 65535 {
			return fmt.Errorf("portlabel: invalid port number in overrides: %d", port)
		}
	}
	return nil
}

// Build validates the config and returns a fully initialised Resolver.
func (c Config) Build() (*Resolver, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	r := New()
	for port, label := range c.Overrides {
		if label == "" {
			return nil, errors.New("portlabel: override label must not be empty")
		}
		r.Set(port, label)
	}
	return r, nil
}
