package portscope_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscope"
)

func TestNew_ContainsDefaultScopes(t *testing.T) {
	r := portscope.New()
	for _, name := range []string{"system", "registered", "dynamic", "all"} {
		if _, err := r.Get(name); err != nil {
			t.Errorf("expected default scope %q to exist, got error: %v", name, err)
		}
	}
}

func TestGet_UnknownScope_ReturnsError(t *testing.T) {
	r := portscope.New()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown scope, got nil")
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	r := portscope.New()
	err := r.Register(portscope.Scope{Name: "", Ranges: [][2]int{{80, 80}}})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestContains_PortInRange(t *testing.T) {
	r := portscope.New()
	s, _ := r.Get("system")
	if !s.Contains(80) {
		t.Error("expected port 80 to be in system scope")
	}
	if s.Contains(8080) {
		t.Error("expected port 8080 to be outside system scope")
	}
}

func TestNames_ReturnsSorted(t *testing.T) {
	r := portscope.New()
	names := r.Names()
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("names not sorted: %v", names)
		}
	}
}

func TestRegister_OverwritesExisting(t *testing.T) {
	r := portscope.New()
	_ = r.Register(portscope.Scope{Name: "custom", Ranges: [][2]int{{100, 200}}})
	_ = r.Register(portscope.Scope{Name: "custom", Ranges: [][2]int{{300, 400}}})
	s, _ := r.Get("custom")
	if s.Contains(150) {
		t.Error("expected old range to be overwritten")
	}
	if !s.Contains(350) {
		t.Error("expected new range to be active")
	}
}
