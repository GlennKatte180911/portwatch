package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	if cfg.MaxFailures != 5 {
		t.Errorf("expected MaxFailures=5, got %d", cfg.MaxFailures)
	}
	if cfg.ResetTimeout != 30*time.Second {
		t.Errorf("expected ResetTimeout=30s, got %v", cfg.ResetTimeout)
	}
}

func TestValidate_ZeroMaxFailures_ReturnsError(t *testing.T) {
	cfg := circuitbreaker.Config{MaxFailures: 0, ResetTimeout: time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero MaxFailures")
	}
}

func TestValidate_ZeroResetTimeout_ReturnsError(t *testing.T) {
	cfg := circuitbreaker.Config{MaxFailures: 3, ResetTimeout: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero ResetTimeout")
	}
}

func TestValidate_ValidConfig_NoError(t *testing.T) {
	cfg := circuitbreaker.Config{MaxFailures: 3, ResetTimeout: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_ReturnsBreaker(t *testing.T) {
	cfg := circuitbreaker.Config{MaxFailures: 2, ResetTimeout: 100 * time.Millisecond}
	b, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Breaker")
	}
	if b.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed state, got %v", b.State())
	}
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := circuitbreaker.Config{MaxFailures: -1, ResetTimeout: time.Second}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
