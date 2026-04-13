package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestNew_CreatesPersistentBaseline(t *testing.T) {
	path := tempPath(t)
	b, err := baseline.New(path, []int{80, 443, 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("baseline file not written: %v", err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	path := tempPath(t)
	ports := []int{22, 80, 443}
	if _, err := baseline.New(path, ports); err != nil {
		t.Fatalf("setup error: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(loaded.Ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(loaded.Ports))
	}
}

func TestLoad_FileNotFound_ReturnsErrNoBaseline(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err != baseline.ErrNoBaseline {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestContains(t *testing.T) {
	b := &baseline.Baseline{Ports: []int{80, 443}}
	if !b.Contains(80) {
		t.Error("expected 80 to be contained")
	}
	if b.Contains(9999) {
		t.Error("expected 9999 to not be contained")
	}
}

func TestUnexpected_ReturnsNewPorts(t *testing.T) {
	path := tempPath(t)
	b, _ := baseline.New(path, []int{80, 443})

	unexpected := b.Unexpected([]int{80, 443, 8080, 9090})
	if len(unexpected) != 2 {
		t.Errorf("expected 2 unexpected ports, got %d: %v", len(unexpected), unexpected)
	}
}

func TestUnexpected_EmptyWhenAllKnown(t *testing.T) {
	path := tempPath(t)
	b, _ := baseline.New(path, []int{80, 443})

	unexpected := b.Unexpected([]int{80, 443})
	if len(unexpected) != 0 {
		t.Errorf("expected no unexpected ports, got %v", unexpected)
	}
}
