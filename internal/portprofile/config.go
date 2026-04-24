package portprofile

import "errors"

// Config holds optional static overrides applied during profile construction.
type Config struct {
	// DefaultScope is used when the scoper returns an empty string.
	DefaultScope string `json:"default_scope"`
	// DefaultRank is used when the ranker returns an empty string.
	DefaultRank string `json:"default_rank"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DefaultScope: "unknown",
		DefaultRank:  "low",
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.DefaultScope == "" {
		return errors.New("portprofile: default_scope must not be empty")
	}
	if c.DefaultRank == "" {
		return errors.New("portprofile: default_rank must not be empty")
	}
	return nil
}

// Apply returns a new Profiler that fills empty scope/rank fields using
// the values from cfg.
func (c Config) Apply(base *Profiler) *Profiler {
	origScoper := base.scoper
	origRanker := base.ranker

	return base.
		WithScoper(func(port int) string {
			if v := origScoper(port); v != "" {
				return v
			}
			return c.DefaultScope
		}).
		WithRanker(func(port int) string {
			if v := origRanker(port); v != "" {
				return v
			}
			return c.DefaultRank
		})
}
