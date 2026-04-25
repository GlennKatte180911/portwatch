package portannounce

import "fmt"

// Config holds configuration for building an Announcer.
type Config struct {
	// MaxSubscribers limits the total number of concurrent subscribers.
	// Zero means unlimited.
	MaxSubscribers int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxSubscribers: 0,
	}
}

// Validate checks that the Config fields are valid.
func (c Config) Validate() error {
	if c.MaxSubscribers < 0 {
		return fmt.Errorf("portannounce: MaxSubscribers must be >= 0, got %d", c.MaxSubscribers)
	}
	return nil
}

// Build constructs an Announcer from the Config.
func (c Config) Build() (*Announcer, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(), nil
}
