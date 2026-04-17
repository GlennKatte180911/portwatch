package portscope_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscope"
)

func TestDefaultConfig_NoExtraScopes(t *testing.T) {
	c := portscope.DefaultConfig()
	if len(c.ExtraScopes) != 0 {
		t.Errorf("expected no extra scopes, got %d", len(c.ExtraScopes))
	}
}

func TestValidate_InvalidRange_ReturnsError(t *testing.T) {
	c := portscope.Config{
		ExtraScopes: []portscope.Scope{
			{Name: "bad", Ranges: [][2]int{{500, 100}}},
		},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestValidate_OutOfBounds_ReturnsError(t *testing.T) {
	c := portscope.Config{
		ExtraScopes: []portscope.Scope{
			{Name: "oob", Ranges: [][2]int{{0, 100}}},
		},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for out-of-bounds range")
	}
}

func TestValidate_ValidConfig_NoError(t *testing.T) {
	c := portscope.Config{
		ExtraScopes: []portscope.Scope{
			{Name: "web", Ranges: [][2]int{{80, 80}, {443, 443}}},
		},
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_ValidConfig_ReturnsRegistry(t *testing.T) {
	c := portscope.Config{
		ExtraScopes: []portscope.Scope{
			{Name: "db", Ranges: [][2]int{{5432, 5432}}},
		},
	}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, err := r.Get("db")
	if err != nil {
		t.Fatalf("expected db scope: %v", err)
	}
	if !s.Contains(5432) {
		t.Error("expected port 5432 in db scope")
	}
}

func TestBuild_InvalidConfig_ReturnsError(t *testing.T) {
	c := portscope.Config{
		ExtraScopes: []portscope.Scope{
			{Name: "", Ranges: [][2]int{{80, 80}}},
		},
	}
	if _, err := c.Build(); err == nil {
		t.Fatal("expected error for invalid config")
	}
}
