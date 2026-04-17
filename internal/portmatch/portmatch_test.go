package portmatch_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmatch"
)

func TestMatch_SinglePort(t *testing.T) {
	m, err := portmatch.New([]string{"80"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(80) {
		t.Error("expected port 80 to match")
	}
	if m.Match(81) {
		t.Error("expected port 81 not to match")
	}
}

func TestMatch_Range(t *testing.T) {
	m, err := portmatch.New([]string{"8000-8100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{8000, 8050, 8100} {
		if !m.Match(p) {
			t.Errorf("expected port %d to match", p)
		}
	}
	if m.Match(7999) || m.Match(8101) {
		t.Error("out-of-range ports should not match")
	}
}

func TestMatch_Alias_Web(t *testing.T) {
	m, err := portmatch.New([]string{"web"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{80, 443, 8080, 8443} {
		if !m.Match(p) {
			t.Errorf("expected web alias to match port %d", p)
		}
	}
}

func TestMatch_Alias_System(t *testing.T) {
	m, err := portmatch.New([]string{"system"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(22) || !m.Match(1) || !m.Match(1023) {
		t.Error("system alias should cover 1-1023")
	}
	if m.Match(1024) {
		t.Error("port 1024 should not be in system range")
	}
}

func TestFilter_ReturnsMatchingPorts(t *testing.T) {
	m, _ := portmatch.New([]string{"80", "443", "8080-8090"})
	input := []int{22, 80, 443, 8085, 9000}
	got := m.Filter(input)
	want := []int{80, 443, 8085}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, p := range want {
		if got[i] != p {
			t.Errorf("index %d: want %d, got %d", i, p, got[i])
		}
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := portmatch.New([]string{"notaport"})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestNew_InvalidRange_LowGreaterThanHigh(t *testing.T) {
	_, err := portmatch.New([]string{"9000-8000"})
	if err == nil {
		t.Error("expected error for inverted range")
	}
}
