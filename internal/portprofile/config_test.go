package portprofile_test

import (
	"testing"

	"github.com/user/portwatch/internal/portprofile"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := portprofile.DefaultConfig()
	if cfg.DefaultScope != "unknown" {
		t.Errorf("expected default_scope 'unknown', got %q", cfg.DefaultScope)
	}
	if cfg.DefaultRank != "low" {
		t.Errorf("expected default_rank 'low', got %q", cfg.DefaultRank)
	}
}

func TestValidate_EmptyScope_ReturnsError(t *testing.T) {
	cfg := portprofile.DefaultConfig()
	cfg.DefaultScope = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty default_scope")
	}
}

func TestValidate_EmptyRank_ReturnsError(t *testing.T) {
	cfg := portprofile.DefaultConfig()
	cfg.DefaultRank = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty default_rank")
	}
}

func TestValidate_ValidConfig_NoError(t *testing.T) {
	cfg := portprofile.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestApply_FillsEmptyScope(t *testing.T) {
	cfg := portprofile.DefaultConfig()
	p := cfg.Apply(portprofile.New())
	pr := p.Build(9999)
	if pr.Scope != "unknown" {
		t.Errorf("expected scope 'unknown', got %q", pr.Scope)
	}
}

func TestApply_FillsEmptyRank(t *testing.T) {
	cfg := portprofile.DefaultConfig()
	p := cfg.Apply(portprofile.New())
	pr := p.Build(9999)
	if pr.Rank != "low" {
		t.Errorf("expected rank 'low', got %q", pr.Rank)
	}
}

func TestApply_DoesNotOverrideExistingScope(t *testing.T) {
	cfg := portprofile.DefaultConfig()
	base := portprofile.New().WithScoper(func(int) string { return "internal" })
	p := cfg.Apply(base)
	pr := p.Build(80)
	if pr.Scope != "internal" {
		t.Errorf("expected scope 'internal', got %q", pr.Scope)
	}
}
