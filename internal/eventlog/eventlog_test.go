package eventlog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/eventlog"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "events.jsonl")
}

func TestAppend_PersistsEntry(t *testing.T) {
	log := eventlog.New(tempLogPath(t))
	if err := log.Append([]int{8080}, []int{22}); err != nil {
		t.Fatalf("Append: %v", err)
	}
	entries, err := log.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Added) != 1 || entries[0].Added[0] != 8080 {
		t.Errorf("unexpected Added: %v", entries[0].Added)
	}
	if len(entries[0].Removed) != 1 || entries[0].Removed[0] != 22 {
		t.Errorf("unexpected Removed: %v", entries[0].Removed)
	}
}

func TestAppend_SkipsEmptyDiff(t *testing.T) {
	path := tempLogPath(t)
	log := eventlog.New(path)
	if err := log.Append(nil, nil); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected no file to be created for empty diff")
	}
}

func TestLoad_FileNotFound_ReturnsEmpty(t *testing.T) {
	log := eventlog.New(tempLogPath(t))
	entries, err := log.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestAppend_MultipleEntries_PreservesOrder(t *testing.T) {
	log := eventlog.New(tempLogPath(t))
	_ = log.Append([]int{80}, nil)
	_ = log.Append([]int{443}, nil)
	_ = log.Append(nil, []int{80})

	entries, err := log.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Added[0] != 80 {
		t.Errorf("first entry: expected Added[0]=80, got %d", entries[0].Added[0])
	}
	if entries[2].Removed[0] != 80 {
		t.Errorf("third entry: expected Removed[0]=80, got %d", entries[2].Removed[0])
	}
}

func TestAppend_TimestampIsSet(t *testing.T) {
	log := eventlog.New(tempLogPath(t))
	_ = log.Append([]int{9000}, nil)
	entries, _ := log.Load()
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
