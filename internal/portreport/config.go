package portreport

import "fmt"

// Format controls the output format of a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Config holds configuration for report generation.
type Config struct {
	Format      Format `json:"format"`
	IncludeNotes bool  `json:"include_notes"`
	MaxEntries  int    `json:"max_entries"` // 0 means unlimited
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Format:       FormatText,
		IncludeNotes: true,
		MaxEntries:   0,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	switch c.Format {
	case FormatText, FormatJSON:
		// valid
	default:
		return fmt.Errorf("portreport: unknown format %q; valid values are text, json", c.Format)
	}
	if c.MaxEntries < 0 {
		return fmt.Errorf("portreport: max_entries must be >= 0, got %d", c.MaxEntries)
	}
	return nil
}

// Apply returns a copy of the config with zero-value fields filled from defaults.
func (c Config) Apply() Config {
	d := DefaultConfig()
	if c.Format == "" {
		c.Format = d.Format
	}
	if c.MaxEntries == 0 {
		c.MaxEntries = d.MaxEntries
	}
	return c
}
