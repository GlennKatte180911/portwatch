package baseline_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func TestEnsure_CreatesBaselineWhenMissing(t *testing.T) {
	var buf bytes.Buffer
	path := filepath.Join(t.TempDir(), "baseline.json")
	m := baseline.NewManager(path, &buf)

	b, err := m.Ensure([]int{80, 443})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(b.Ports))
	}
	if !strings.Contains(buf.String(), "created") {
		t.Errorf("expected 'created' in output, got: %s", buf.String())
	}
}

func TestEnsure_LoadsExistingBaseline(t *testing.T) {
	var buf bytes.Buffer
	path := filepath.Join(t.TempDir(), "baseline.json")

	// Pre-create a baseline.
	if _, err := baseline.New(path, []int{22, 80}); err != nil {
		t.Fatalf("setup error: %v", err)
	}

	m := baseline.NewManager(path, &buf)
	b, err := m.Ensure([]int{9999})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should load existing — not overwrite with 9999.
	if len(b.Ports) != 2 {
		t.Errorf("expected 2 ports from existing baseline, got %d", len(b.Ports))
	}
	if !strings.Contains(buf.String(), "loaded") {
		t.Errorf("expected 'loaded' in output, got: %s", buf.String())
	}
}

func TestReset_OverwritesExistingBaseline(t *testing.T) {
	var buf bytes.Buffer
	path := filepath.Join(t.TempDir(), "baseline.json")

	if _, err := baseline.New(path, []int{22, 80}); err != nil {
		t.Fatalf("setup error: %v", err)
	}

	m := baseline.NewManager(path, &buf)
	b, err := m.Reset([]int{443, 8080, 9090})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.Ports) != 3 {
		t.Errorf("expected 3 ports after reset, got %d", len(b.Ports))
	}
	if !strings.Contains(buf.String(), "reset") {
		t.Errorf("expected 'reset' in output, got: %s", buf.String())
	}
}
