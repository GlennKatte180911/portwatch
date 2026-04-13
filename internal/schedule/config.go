package schedule

import (
	"errors"
	"time"
)

// Config holds parameters for constructing a Scheduler.
type Config struct {
	// Interval is the base period between scan ticks.
	Interval time.Duration `json:"interval"`
	// Jitter adds a random offset in [0, Jitter) to each interval.
	// Set to 0 to disable.
	Jitter time.Duration `json:"jitter"`
}

// DefaultConfig returns a Config with sensible defaults suitable for
// interactive use: 60-second interval with 5-second jitter.
func DefaultConfig() Config {
	return Config{
		Interval: 60 * time.Second,
		Jitter:   5 * time.Second,
	}
}

// Validate returns an error if the Config contains invalid values.
func (c Config) Validate() error {
	if c.Interval <= 0 {
		return errors.New("schedule: interval must be greater than zero")
	}
	if c.Jitter < 0 {
		return errors.New("schedule: jitter must be non-negative")
	}
	return nil
}

// Build constructs and starts a Scheduler from the Config.
// Returns an error if Validate fails.
func (c Config) Build() (*Scheduler, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Interval, c.Jitter), nil
}
