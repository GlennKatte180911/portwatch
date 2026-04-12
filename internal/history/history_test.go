package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "history.json")
}

func TestRecord_PersistsEntry(t *testing.T) {
	h := history.New(tempHistoryPath(t))
	if err := h.Record([]int{8080}, []int{22}); err != nil {
		t.Fatalf("Record() error: %v", err)
	}
	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	e := h.Entries[0]
	if len(e.Added) != 1 || e.Added[0] != 8080 {
		t.Errorf("unexpected Added: %v", e.Added)
	}
	if len(e.Removed) != 1 || e.Removed[0] != 22 {
		t.Errorf("unexpected Removed: %v", e.Removed)
	}
}

func TestRecord_SkipsEmptyDiff(t *testing.T) {
	h := history.New(tempHistoryPath(t))
	if err := h.Record(nil, nil); err != nil {
		t.Fatalf("Record() error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected 0 entries for empty diff, got %d", len(h.Entries))
	}
}

func TestLoad_FileNotFound_ReturnsEmpty(t *testing.T) {
	h := history.New(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err := h.Load(); err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(h.Entries))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempHistoryPath(t)
	h1 := history.New(path)
	_ = h1.Record([]int{443}, nil)
	_ = h1.Record([]int{8443}, []int{80})

	h2 := history.New(path)
	if err := h2.Load(); err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(h2.Entries) != 2 {
		t.Fatalf("expected 2 entries after reload, got %d", len(h2.Entries))
	}
}

func TestLast_ReturnsNewestFirst(t *testing.T) {
	h := history.New(tempHistoryPath(t))
	old := time.Now().Add(-2 * time.Hour)
	recent := time.Now()
	h.Entries = []history.Entry{
		{Timestamp: old, Added: []int{22}, Removed: nil},
		{Timestamp: recent, Added: []int{8080}, Removed: nil},
	}
	last := h.Last(1)
	if len(last) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(last))
	}
	if last[0].Added[0] != 8080 {
		t.Errorf("expected most recent entry first, got added=%v", last[0].Added)
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempHistoryPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	h := history.New(path)
	if err := h.Load(); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
