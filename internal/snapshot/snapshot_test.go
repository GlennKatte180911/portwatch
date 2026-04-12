package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNew_SetsTimestampAndPorts(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.New([]int{80, 443, 8080})
	after := time.Now().UTC()

	if s.Timestamp.Before(before) || s.Timestamp.After(after) {
		t.Errorf("unexpected timestamp: %v", s.Timestamp)
	}
	if len(s.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(s.Ports))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "snap.json")

	orig := snapshot.New([]int{22, 80, 443})
	if err := orig.Save(tmp); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("port count mismatch: got %d, want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestDiff_DetectsAddedAndRemoved(t *testing.T) {
	prev := &snapshot.Snapshot{Ports: []int{22, 80, 443}}
	curr := &snapshot.Snapshot{Ports: []int{80, 443, 8080}}

	added, removed := snapshot.Diff(prev, curr)

	if len(added) != 1 || added[0] != 8080 {
		t.Errorf("expected added=[8080], got %v", added)
	}
	if len(removed) != 1 || removed[0] != 22 {
		t.Errorf("expected removed=[22], got %v", removed)
	}
}

func TestDiff_NoDifference(t *testing.T) {
	prev := &snapshot.Snapshot{Ports: []int{80, 443}}
	curr := &snapshot.Snapshot{Ports: []int{80, 443}}

	added, removed := snapshot.Diff(prev, curr)

	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}

func TestDiff_EmptyPrevious(t *testing.T) {
	_ = os.Getenv("CI") // satisfy import
	prev := &snapshot.Snapshot{Ports: []int{}}
	curr := &snapshot.Snapshot{Ports: []int{80, 443}}

	added, removed := snapshot.Diff(prev, curr)

	if len(added) != 2 {
		t.Errorf("expected 2 added ports, got %d", len(added))
	}
	if len(removed) != 0 {
		t.Errorf("expected 0 removed ports, got %d", len(removed))
	}
}
