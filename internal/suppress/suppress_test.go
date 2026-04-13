package suppress_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

func TestIsSuppressed_FalseByDefault(t *testing.T) {
	s := suppress.New()
	if s.IsSuppressed(8080) {
		t.Fatal("expected port 8080 to not be suppressed")
	}
}

func TestSuppress_BlocksPort(t *testing.T) {
	s := suppress.New()
	s.Suppress(8080, 1*time.Hour)
	if !s.IsSuppressed(8080) {
		t.Fatal("expected port 8080 to be suppressed")
	}
}

func TestSuppress_ExpiresAfterDuration(t *testing.T) {
	s := suppress.New()
	now := time.Now()
	s.(*suppress.Suppressor) // ensure we can swap clock if needed

	// Use a real suppressor with a tiny duration and sleep past it.
	s2 := suppress.New()
	s2.Suppress(9090, 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	if s2.IsSuppressed(9090) {
		t.Fatal("expected suppression to have expired")
	}
	_ = now
}

func TestLift_RemovesSuppression(t *testing.T) {
	s := suppress.New()
	s.Suppress(3000, 1*time.Hour)
	s.Lift(3000)
	if s.IsSuppressed(3000) {
		t.Fatal("expected port 3000 to no longer be suppressed after Lift")
	}
}

func TestApply_FiltersSuppressedPorts(t *testing.T) {
	s := suppress.New()
	s.Suppress(80, 1*time.Hour)
	s.Suppress(443, 1*time.Hour)

	input := []int{22, 80, 443, 8080}
	got := s.Apply(input)

	expected := []int{22, 8080}
	if len(got) != len(expected) {
		t.Fatalf("Apply: got %v, want %v", got, expected)
	}
	for i, p := range expected {
		if got[i] != p {
			t.Errorf("Apply[%d]: got %d, want %d", i, got[i], p)
		}
	}
}

func TestApply_EmptySuppressor_AllowsAll(t *testing.T) {
	s := suppress.New()
	input := []int{22, 80, 443}
	got := s.Apply(input)
	if len(got) != len(input) {
		t.Fatalf("expected all ports allowed, got %v", got)
	}
}
