package portlabel_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portlabel"
)

func TestLabel_KnownPort_ReturnsName(t *testing.T) {
	r := portlabel.New()
	if got := r.Label(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestLabel_UnknownPort_ReturnsFallback(t *testing.T) {
	r := portlabel.New()
	got := r.Label(9999)
	if got != "port/9999" {
		t.Fatalf("expected port/9999, got %q", got)
	}
}

func TestSet_OverwritesExistingLabel(t *testing.T) {
	r := portlabel.New()
	r.Set(80, "my-web")
	if got := r.Label(80); got != "my-web" {
		t.Fatalf("expected my-web, got %q", got)
	}
}

func TestSet_AddsNewLabel(t *testing.T) {
	r := portlabel.New()
	r.Set(12345, "custom-svc")
	if got := r.Label(12345); got != "custom-svc" {
		t.Fatalf("expected custom-svc, got %q", got)
	}
}

func TestAnnotate_ReturnsPairs(t *testing.T) {
	r := portlabel.New()
	pairs := r.Annotate([]int{22, 80, 9999})
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
	if pairs[0].Port != 22 || pairs[0].Label != "ssh" {
		t.Errorf("unexpected pair[0]: %+v", pairs[0])
	}
	if pairs[1].Port != 80 || pairs[1].Label != "http" {
		t.Errorf("unexpected pair[1]: %+v", pairs[1])
	}
	if !strings.HasPrefix(pairs[2].Label, "port/") {
		t.Errorf("expected fallback label for unknown port, got %q", pairs[2].Label)
	}
}

func TestAnnotate_EmptySlice_ReturnsEmpty(t *testing.T) {
	r := portlabel.New()
	if got := r.Annotate(nil); len(got) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(got))
	}
}

func TestNew_ContainsWellKnownPorts(t *testing.T) {
	r := portlabel.New()
	expected := map[int]string{
		443:  "https",
		3306: "mysql",
		6379: "redis",
	}
	for port, want := range expected {
		if got := r.Label(port); got != want {
			t.Errorf("port %d: expected %q, got %q", port, want, got)
		}
	}
}
