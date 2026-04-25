package portbudget_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/portbudget"
)

func TestNew_ValidLimit_ReturnsBudget(t *testing.T) {
	b, err := portbudget.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Limit() != 10 {
		t.Fatalf("expected limit 10, got %d", b.Limit())
	}
}

func TestNew_ZeroLimit_ReturnsError(t *testing.T) {
	_, err := portbudget.New(0)
	if err == nil {
		t.Fatal("expected error for zero limit, got nil")
	}
}

func TestNew_NegativeLimit_ReturnsError(t *testing.T) {
	_, err := portbudget.New(-5)
	if err == nil {
		t.Fatal("expected error for negative limit, got nil")
	}
}

func TestCheck_NoViolation_WhenUnderLimit(t *testing.T) {
	b, _ := portbudget.New(5)
	ports := []int{80, 443, 8080}
	if err := b.Check(ports); err != nil {
		t.Fatalf("unexpected violation: %v", err)
	}
}

func TestCheck_NoViolation_WhenAtLimit(t *testing.T) {
	b, _ := portbudget.New(3)
	ports := []int{80, 443, 8080}
	if err := b.Check(ports); err != nil {
		t.Fatalf("unexpected violation at exact limit: %v", err)
	}
}

func TestCheck_ReturnsViolation_WhenOverLimit(t *testing.T) {
	b, _ := portbudget.New(2)
	ports := []int{80, 443, 8080}
	err := b.Check(ports)
	if err == nil {
		t.Fatal("expected violation, got nil")
	}
	var v portbudget.Violation
	if !errors.As(err, &v) {
		t.Fatalf("expected Violation type, got %T", err)
	}
	if v.Limit != 2 || v.Actual != 3 || v.Excess != 1 {
		t.Fatalf("unexpected violation fields: %+v", v)
	}
}

func TestWithin_ReturnsTrueUnderLimit(t *testing.T) {
	b, _ := portbudget.New(10)
	if !b.Within([]int{22, 80}) {
		t.Fatal("expected Within to return true")
	}
}

func TestWithin_ReturnsFalseOverLimit(t *testing.T) {
	b, _ := portbudget.New(1)
	if b.Within([]int{22, 80}) {
		t.Fatal("expected Within to return false")
	}
}

func TestSetLimit_UpdatesLimit(t *testing.T) {
	b, _ := portbudget.New(5)
	if err := b.SetLimit(20); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Limit() != 20 {
		t.Fatalf("expected limit 20, got %d", b.Limit())
	}
}

func TestSetLimit_ZeroReturnsError(t *testing.T) {
	b, _ := portbudget.New(5)
	if err := b.SetLimit(0); err == nil {
		t.Fatal("expected error for zero limit update, got nil")
	}
}
