package aggregator_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/aggregator"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := aggregator.DefaultConfig()
	if cfg.Window != 5*time.Second {
		t.Fatalf("expected 5s window, got %v", cfg.Window)
	}
}

func TestValidate_ZeroWindow_ReturnsError(t *testing.T) {
	cfg := aggregator.Config{Window: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestValidate_NegativeWindow_ReturnsError(t *testing.T) {
	cfg := aggregator.Config{Window: -1 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestValidate_PositiveWindow_NoError(t *testing.T) {
	cfg := aggregator.Config{Window: 100 * time.Millisecond}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_ReturnsAggregator(t *testing.T) {
	cfg := aggregator.Config{Window: 50 * time.Millisecond}
	agg, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agg == nil {
		t.Fatal("expected non-nil aggregator")
	}
	agg.Stop()
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := aggregator.Config{Window: 0}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error from Build with zero window")
	}
}
