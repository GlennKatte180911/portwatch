package healthcheck_test

import (
	"testing"

	"github.com/user/portwatch/internal/healthcheck"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := healthcheck.DefaultConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if cfg.Addr != ":9090" {
		t.Errorf("expected default addr :9090, got %s", cfg.Addr)
	}
}

func TestValidate_EmptyAddr_ReturnsError(t *testing.T) {
	cfg := healthcheck.Config{Enabled: true, Addr: ""}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty addr")
	}
}

func TestValidate_ValidConfig_NoError(t *testing.T) {
	cfg := healthcheck.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_ReturnsChecker(t *testing.T) {
	cfg := healthcheck.DefaultConfig()
	checker, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if checker == nil {
		t.Fatal("expected non-nil checker")
	}
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := healthcheck.Config{Enabled: true, Addr: ""}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
