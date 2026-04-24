package portquota_test

import (
	"testing"

	"github.com/user/portwatch/internal/portquota"
)

func TestSet_ValidLimit_NoError(t *testing.T) {
	q := portquota.New()
	if err := q.Set("web", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSet_ZeroLimit_ReturnsError(t *testing.T) {
	q := portquota.New()
	if err := q.Set("web", 0); err == nil {
		t.Fatal("expected error for zero limit, got nil")
	}
}

func TestSet_NegativeLimit_ReturnsError(t *testing.T) {
	q := portquota.New()
	if err := q.Set("db", -1); err == nil {
		t.Fatal("expected error for negative limit, got nil")
	}
}

func TestCheck_NoViolation_WhenUnderLimit(t *testing.T) {
	q := portquota.New()
	_ = q.Set("web", 10)
	violations := q.Check(map[string]int{"web": 5})
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_ReturnsViolation_WhenOverLimit(t *testing.T) {
	q := portquota.New()
	_ = q.Set("web", 3)
	violations := q.Check(map[string]int{"web": 7})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	v := violations[0]
	if v.Group != "web" {
		t.Errorf("expected group 'web', got %q", v.Group)
	}
	if v.Limit != 3 {
		t.Errorf("expected limit 3, got %d", v.Limit)
	}
	if v.Actual != 7 {
		t.Errorf("expected actual 7, got %d", v.Actual)
	}
}

func TestCheck_IgnoresGroupsWithNoLimit(t *testing.T) {
	q := portquota.New()
	violations := q.Check(map[string]int{"unknown": 999})
	if len(violations) != 0 {
		t.Fatalf("expected no violations for unconfigured group, got %d", len(violations))
	}
}

func TestRemove_DeletesLimit(t *testing.T) {
	q := portquota.New()
	_ = q.Set("db", 2)
	q.Remove("db")
	violations := q.Check(map[string]int{"db": 100})
	if len(violations) != 0 {
		t.Fatal("expected no violations after removal")
	}
}

func TestLimits_ReturnsCopy(t *testing.T) {
	q := portquota.New()
	_ = q.Set("web", 5)
	_ = q.Set("db", 2)
	limits := q.Limits()
	if len(limits) != 2 {
		t.Fatalf("expected 2 limits, got %d", len(limits))
	}
	// Mutating the returned map must not affect the quota.
	delete(limits, "web")
	if len(q.Limits()) != 2 {
		t.Error("mutation of returned map affected internal state")
	}
}

func TestViolation_String_ContainsGroupAndCounts(t *testing.T) {
	v := portquota.Violation{Group: "web", Limit: 3, Actual: 7}
	s := v.String()
	for _, want := range []string{"web", "3", "7"} {
		if !containsStr(s, want) {
			t.Errorf("Violation.String() missing %q: %s", want, s)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
