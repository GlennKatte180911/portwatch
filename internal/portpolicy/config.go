package portpolicy

import (
	"fmt"
)

// RuleConfig is a serialisable representation of a single rule.
type RuleConfig struct {
	Ports  []int  `json:"ports"`
	Action string `json:"action"`
	Reason string `json:"reason,omitempty"`
}

// Config holds the full policy configuration.
type Config struct {
	Default string       `json:"default"`
	Rules   []RuleConfig `json:"rules,omitempty"`
}

// DefaultConfig returns a permissive policy with no rules.
func DefaultConfig() Config {
	return Config{Default: string(ActionAllow)}
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	a := Action(c.Default)
	if a != ActionAllow && a != ActionDeny {
		return fmt.Errorf("portpolicy: invalid default action %q", c.Default)
	}
	for i, r := range c.Rules {
		ra := Action(r.Action)
		if ra != ActionAllow && ra != ActionDeny {
			return fmt.Errorf("portpolicy: rule[%d] invalid action %q", i, r.Action)
		}
		if len(r.Ports) == 0 {
			return fmt.Errorf("portpolicy: rule[%d] has no ports", i)
		}
	}
	return nil
}

// Build constructs a Policy from the configuration.
func (c Config) Build() (*Policy, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	p, err := New(Action(c.Default))
	if err != nil {
		return nil, err
	}
	for _, rc := range c.Rules {
		if err := p.AddRule(Rule{
			Ports:  rc.Ports,
			Action: Action(rc.Action),
			Reason: rc.Reason,
		}); err != nil {
			return nil, err
		}
	}
	return p, nil
}
