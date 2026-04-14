package portaudit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/portaudit"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.json")
}

func TestNew_FileNotFound_ReturnsEmptyAuditor(t *testing.T) {
	a, err := portaudit.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(a.Entries()); got != 0 {
		t.Errorf("expected 0 entries, got %d", got)
	}
}

func TestRecord_PersistsEntry(t *testing.T) {
	path := tempPath(t)
	a, _ := portaudit.New(path)

	if err := a.Record("manual", []int{8080}, nil); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	entries := a.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Trigger != "manual" {
		t.Errorf("expected trigger 'manual', got %q", entries[0].Trigger)
	}
	if len(entries[0].Added) != 1 || entries[0].Added[0] != 8080 {
		t.Errorf("unexpected added ports: %v", entries[0].Added)
	}
}

func TestRecord_SkipsEmptyDiff(t *testing.T) {
	path := tempPath(t)
	a, _ := portaudit.New(path)

	if err := a.Record("scheduler", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(a.Entries()); got != 0 {
		t.Errorf("expected 0 entries, got %d", got)
	}
	// File should not be created for empty diffs.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to not exist for empty diff")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempPath(t)
	a1, _ := portaudit.New(path)
	_ = a1.Record("scheduler", []int{443}, []int{80})

	a2, err := portaudit.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	entries := a2.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(entries))
	}
	if entries[0].Trigger != "scheduler" {
		t.Errorf("trigger mismatch: got %q", entries[0].Trigger)
	}
	if len(entries[0].Removed) != 1 || entries[0].Removed[0] != 80 {
		t.Errorf("removed mismatch: %v", entries[0].Removed)
	}
}

func TestRecord_MultipleEntries_PreservesOrder(t *testing.T) {
	a, _ := portaudit.New(tempPath(t))
	_ = a.Record("first", []int{22}, nil)
	_ = a.Record("second", []int{443}, nil)

	entries := a.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Trigger != "first" || entries[1].Trigger != "second" {
		t.Errorf("order not preserved: %v", entries)
	}
}
