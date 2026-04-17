package portscope

import "fmt"

// Config holds configuration for building a scope-aware registry.
type Config struct {
	// ExtraScopes are user-defined scopes merged with the defaults.
	ExtraScopes []Scope
}

// DefaultConfig returns a Config with no extra scopes.
func DefaultConfig() Config {
	return Config{}
}

// Validate returns an error if any extra scope is misconfigured.
func (c Config) Validate() error {
	for _, s := range c.ExtraScopes {
		if s.Name == "" {
			return fmt.Errorf("portscope: extra scope must have a name")
		}
		for _, rng := range s.Ranges {
			if rng[0] > rng[1] {
				return fmt.Errorf("portscope: invalid range [%d, %d] in scope %q", rng[0], rng[1], s.Name)
			}
			if rng[0] < 1 || rng[1] > 65535 {
				return fmt.Errorf("portscope: range [%d, %d] out of bounds in scope %q", rng[0], rng[1], s.Name)
			}
		}
	}
	return nil
}

// Build validates the config and returns a populated Registry.
func (c Config) Build() (*Registry, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	r := New()
	for _, s := range c.ExtraScopes {
		_ = r.Register(s)
	}
	return r, nil
}
