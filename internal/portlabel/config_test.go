package portlabel_test

import (
	"testing"

	"github.com/user/portwatch/internal/portlabel"
)

func TestDefaultConfig_HasEmptyOverrides(t *testing.T) {
	cfg := portlabel.DefaultConfig()
	if len(cfg.Overrides) != 0 {
		t.Fatalf("expected empty overrides, got %d entries", len(cfg.Overrides))
	}
}

func TestValidate_InvalidPort_ReturnsError(t *testing.T) {
	cfg := portlabel.Config{Overrides: map[int]string{0: "zero"}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 0, got nil")
	}
}

func TestValidate_ValidPorts_NoError(t *testing.T) {
	cfg := portlabel.Config{Overrides: map[int]string{8080: "dev-server"}}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_AppliesOverrides(t *testing.T) {
	cfg := portlabel.Config{Overrides: map[int]string{9090: "prometheus"}}
	r, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Label(9090); got != "prometheus" {
		t.Fatalf("expected prometheus, got %q", got)
	}
}

func TestBuild_EmptyLabel_ReturnsError(t *testing.T) {
	cfg := portlabel.Config{Overrides: map[int]string{8080: ""}}
	if _, err := cfg.Build(); err == nil {
		t.Fatal("expected error for empty label, got nil")
	}
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := portlabel.Config{Overrides: map[int]string{99999: "oob"}}
	if _, err := cfg.Build(); err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
