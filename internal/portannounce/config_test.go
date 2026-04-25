package portannounce_test

import (
	"testing"

	"github.com/user/portwatch/internal/portannounce"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := portannounce.DefaultConfig()
	if cfg.MaxSubscribers != 0 {
		t.Errorf("expected MaxSubscribers=0, got %d", cfg.MaxSubscribers)
	}
}

func TestValidate_NegativeMaxSubscribers_ReturnsError(t *testing.T) {
	cfg := portannounce.Config{MaxSubscribers: -1}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative MaxSubscribers")
	}
}

func TestValidate_ZeroMaxSubscribers_NoError(t *testing.T) {
	cfg := portannounce.Config{MaxSubscribers: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_PositiveMaxSubscribers_NoError(t *testing.T) {
	cfg := portannounce.Config{MaxSubscribers: 10}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_ReturnsAnnouncer(t *testing.T) {
	cfg := portannounce.DefaultConfig()
	a, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Error("expected non-nil Announcer")
	}
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := portannounce.Config{MaxSubscribers: -5}
	_, err := cfg.Build()
	if err == nil {
		t.Error("expected error from invalid config")
	}
}
