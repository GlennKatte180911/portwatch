package portguard_test

import (
	"testing"

	"github.com/user/portwatch/internal/portguard"
)

func TestEvaluate_EmptyGuard_AllowsAll(t *testing.T) {
	g := portguard.New()
	for _, port := range []int{22, 80, 443, 8080} {
		if got := g.Evaluate(port); got != portguard.Allow {
			t.Errorf("port %d: expected Allow, got %v", port, got)
		}
	}
}

func TestEvaluate_DenylistTakesPrecedence(t *testing.T) {
	g := portguard.New()
	_ = g.Permit(80)
	_ = g.Block(80)
	if got := g.Evaluate(80); got != portguard.Deny {
		t.Errorf("expected Deny, got %v", got)
	}
}

func TestEvaluate_AllowlistRestrictsUnlistedPorts(t *testing.T) {
	g := portguard.New()
	_ = g.Permit(443)
	if got := g.Evaluate(443); got != portguard.Allow {
		t.Errorf("port 443: expected Allow, got %v", got)
	}
	if got := g.Evaluate(80); got != portguard.Deny {
		t.Errorf("port 80: expected Deny, got %v", got)
	}
}

func TestBlock_InvalidPort_ReturnsError(t *testing.T) {
	g := portguard.New()
	if err := g.Block(0); err == nil {
		t.Error("expected error for port 0, got nil")
	}
	if err := g.Block(65536); err == nil {
		t.Error("expected error for port 65536, got nil")
	}
}

func TestPermit_InvalidPort_ReturnsError(t *testing.T) {
	g := portguard.New()
	if err := g.Permit(-1); err == nil {
		t.Error("expected error for port -1, got nil")
	}
}

func TestApply_FiltersPortsCorrectly(t *testing.T) {
	g := portguard.New()
	_ = g.Permit(22)
	_ = g.Permit(443)

	input := []int{22, 80, 443, 8080}
	got := g.Apply(input)

	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d: %v", len(got), got)
	}
	if got[0] != 22 || got[1] != 443 {
		t.Errorf("unexpected ports: %v", got)
	}
}

func TestApply_EmptyGuard_ReturnsAll(t *testing.T) {
	g := portguard.New()
	input := []int{22, 80, 443}
	got := g.Apply(input)
	if len(got) != len(input) {
		t.Errorf("expected %d ports, got %d", len(input), len(got))
	}
}

func TestApply_AllDenied_ReturnsEmpty(t *testing.T) {
	g := portguard.New()
	_ = g.Block(22)
	_ = g.Block(80)
	got := g.Apply([]int{22, 80})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
