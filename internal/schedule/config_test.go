package schedule_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := schedule.DefaultConfig()
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", cfg.Interval)
	}
	if cfg.Jitter != 5*time.Second {
		t.Errorf("expected 5s jitter, got %v", cfg.Jitter)
	}
}

func TestValidate_ZeroInterval_ReturnsError(t *testing.T) {
	cfg := schedule.Config{Interval: 0, Jitter: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidate_NegativeJitter_ReturnsError(t *testing.T) {
	cfg := schedule.Config{Interval: time.Second, Jitter: -time.Millisecond}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative jitter")
	}
}

func TestValidate_ValidConfig_NoError(t *testing.T) {
	cfg := schedule.Config{Interval: 10 * time.Second, Jitter: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_ReturnsScheduler(t *testing.T) {
	cfg := schedule.Config{Interval: 20 * time.Millisecond, Jitter: 0}
	s, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer s.Stop()

	select {
	case <-s.C:
		// scheduler is running
	case <-time.After(200 * time.Millisecond):
		t.Fatal("scheduler did not emit tick")
	}
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := schedule.Config{Interval: 0}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error from Build with invalid config")
	}
}
