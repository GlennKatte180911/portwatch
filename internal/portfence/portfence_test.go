package portfence_test

import (
	"testing"

	"github.com/user/portwatch/internal/portfence"
)

func TestNew_ValidPorts_NoError(t *testing.T) {
	f, err := portfence.New([]int{80, 443, 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Fence")
	}
}

func TestNew_InvalidPort_ReturnsError(t *testing.T) {
	_, err := portfence.New([]int{80, 0})
	if err == nil {
		t.Fatal("expected error for port 0")
	}
	_, err = portfence.New([]int{65536})
	if err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestCheck_AllowsPermittedPort(t *testing.T) {
	f, _ := portfence.New([]int{22, 80})
	if got := f.Check(80); got != portfence.Allow {
		t.Errorf("expected Allow, got %s", got)
	}
}

func TestCheck_DeniesUnpermittedPort(t *testing.T) {
	f, _ := portfence.New([]int{22, 80})
	if got := f.Check(9999); got != portfence.Deny {
		t.Errorf("expected Deny, got %s", got)
	}
}

func TestPermit_AddsPort(t *testing.T) {
	f, _ := portfence.New([]int{80})
	if err := f.Permit(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := f.Check(443); got != portfence.Allow {
		t.Errorf("expected Allow after Permit, got %s", got)
	}
}

func TestRevoke_RemovesPort(t *testing.T) {
	f, _ := portfence.New([]int{80, 443})
	f.Revoke(80)
	if got := f.Check(80); got != portfence.Deny {
		t.Errorf("expected Deny after Revoke, got %s", got)
	}
}

func TestEvaluate_ReturnsViolations(t *testing.T) {
	f, _ := portfence.New([]int{80, 443})
	violations := f.Evaluate([]int{80, 22, 443, 3306})
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
	for _, v := range violations {
		if v.Verdict != portfence.Deny {
			t.Errorf("expected Deny verdict, got %s", v.Verdict)
		}
	}
}

func TestEvaluate_NoViolations_WhenAllAllowed(t *testing.T) {
	f, _ := portfence.New([]int{80, 443})
	violations := f.Evaluate([]int{80, 443})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestPermitted_ReturnsSortedPorts(t *testing.T) {
	f, _ := portfence.New([]int{443, 22, 80})
	got := f.Permitted()
	expected := []int{22, 80, 443}
	for i, p := range expected {
		if got[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, got[i])
		}
	}
}
