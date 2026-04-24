package portpolicy_test

import (
	"testing"

	"github.com/user/portwatch/internal/portpolicy"
)

func TestNew_InvalidAction_ReturnsError(t *testing.T) {
	_, err := portpolicy.New("unknown")
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestNew_ValidAction_ReturnsPolicy(t *testing.T) {
	p, err := portpolicy.New(portpolicy.ActionAllow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil policy")
	}
}

func TestEvaluate_DefaultAllow_WhenNoRules(t *testing.T) {
	p, _ := portpolicy.New(portpolicy.ActionAllow)
	r := p.Evaluate(8080)
	if r.Action != portpolicy.ActionAllow {
		t.Errorf("expected allow, got %s", r.Action)
	}
	if r.Reason != "default" {
		t.Errorf("expected reason 'default', got %s", r.Reason)
	}
}

func TestEvaluate_DefaultDeny_WhenNoRules(t *testing.T) {
	p, _ := portpolicy.New(portpolicy.ActionDeny)
	r := p.Evaluate(22)
	if r.Action != portpolicy.ActionDeny {
		t.Errorf("expected deny, got %s", r.Action)
	}
}

func TestEvaluate_MatchingRuleWins(t *testing.T) {
	p, _ := portpolicy.New(portpolicy.ActionAllow)
	_ = p.AddRule(portpolicy.Rule{
		Ports:  []int{22, 23},
		Action: portpolicy.ActionDeny,
		Reason: "block telnet and ssh",
	})
	r := p.Evaluate(22)
	if r.Action != portpolicy.ActionDeny {
		t.Errorf("expected deny for port 22, got %s", r.Action)
	}
	if r.Reason != "block telnet and ssh" {
		t.Errorf("unexpected reason: %s", r.Reason)
	}
}

func TestEvaluate_NonMatchingPort_UsesDefault(t *testing.T) {
	p, _ := portpolicy.New(portpolicy.ActionDeny)
	_ = p.AddRule(portpolicy.Rule{
		Ports:  []int{443},
		Action: portpolicy.ActionAllow,
		Reason: "https ok",
	})
	r := p.Evaluate(8080)
	if r.Action != portpolicy.ActionDeny {
		t.Errorf("expected deny for port 8080, got %s", r.Action)
	}
}

func TestApply_ReturnsDeniedPortsOnly(t *testing.T) {
	p, _ := portpolicy.New(portpolicy.ActionAllow)
	_ = p.AddRule(portpolicy.Rule{
		Ports:  []int{23, 3389},
		Action: portpolicy.ActionDeny,
		Reason: "legacy",
	})
	denied := p.Apply([]int{80, 443, 23, 3389, 8080})
	if len(denied) != 2 {
		t.Fatalf("expected 2 denied, got %d", len(denied))
	}
	for _, r := range denied {
		if r.Action != portpolicy.ActionDeny {
			t.Errorf("expected deny, got %s", r.Action)
		}
	}
}

func TestAddRule_InvalidAction_ReturnsError(t *testing.T) {
	p, _ := portpolicy.New(portpolicy.ActionAllow)
	err := p.AddRule(portpolicy.Rule{Ports: []int{80}, Action: "maybe"})
	if err == nil {
		t.Fatal("expected error for invalid rule action")
	}
}
