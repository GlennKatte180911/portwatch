package portpolicy_test

import (
	"testing"

	"github.com/user/portwatch/internal/portpolicy"
)

func TestDefaultConfig_Values(t *testing.T) {
	c := portpolicy.DefaultConfig()
	if c.Default != "allow" {
		t.Errorf("expected default 'allow', got %s", c.Default)
	}
	if len(c.Rules) != 0 {
		t.Errorf("expected no rules, got %d", len(c.Rules))
	}
}

func TestValidate_InvalidDefault_ReturnsError(t *testing.T) {
	c := portpolicy.Config{Default: "maybe"}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for invalid default")
	}
}

func TestValidate_InvalidRuleAction_ReturnsError(t *testing.T) {
	c := portpolicy.Config{
		Default: "allow",
		Rules: []portpolicy.RuleConfig{
			{Ports: []int{22}, Action: "skip"},
		},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for invalid rule action")
	}
}

func TestValidate_EmptyPorts_ReturnsError(t *testing.T) {
	c := portpolicy.Config{
		Default: "deny",
		Rules: []portpolicy.RuleConfig{
			{Ports: []int{}, Action: "allow"},
		},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for rule with no ports")
	}
}

func TestValidate_ValidConfig_NoError(t *testing.T) {
	c := portpolicy.Config{
		Default: "deny",
		Rules: []portpolicy.RuleConfig{
			{Ports: []int{80, 443}, Action: "allow", Reason: "web"},
		},
	}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_ReturnsPolicy(t *testing.T) {
	c := portpolicy.Config{
		Default: "deny",
		Rules: []portpolicy.RuleConfig{
			{Ports: []int{443}, Action: "allow", Reason: "https"},
		},
	}
	p, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := p.Evaluate(443)
	if r.Action != portpolicy.ActionAllow {
		t.Errorf("expected allow for 443, got %s", r.Action)
	}
	r2 := p.Evaluate(22)
	if r2.Action != portpolicy.ActionDeny {
		t.Errorf("expected deny for 22, got %s", r2.Action)
	}
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	c := portpolicy.Config{Default: "bad"}
	_, err := c.Build()
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
