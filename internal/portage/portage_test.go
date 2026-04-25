package portage

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestObserve_SetsFirstAndLastSeen(t *testing.T) {
	base := time.Now()
	tr := New(fixedNow(base))
	tr.Observe([]int{80})
	e, ok := tr.Get(80)
	if !ok {
		t.Fatal("expected entry for port 80")
	}
	if !e.FirstSeen.Equal(base) {
		t.Errorf("FirstSeen = %v; want %v", e.FirstSeen, base)
	}
	if !e.LastSeen.Equal(base) {
		t.Errorf("LastSeen = %v; want %v", e.LastSeen, base)
	}
}

func TestObserve_UpdatesLastSeenOnly(t *testing.T) {
	base := time.Now()
	later := base.Add(5 * time.Second)
	tr := New(fixedNow(base))
	tr.Observe([]int{443})
	tr.now = fixedNow(later)
	tr.Observe([]int{443})
	e, _ := tr.Get(443)
	if !e.FirstSeen.Equal(base) {
		t.Errorf("FirstSeen changed; want %v got %v", base, e.FirstSeen)
	}
	if !e.LastSeen.Equal(later) {
		t.Errorf("LastSeen = %v; want %v", e.LastSeen, later)
	}
}

func TestAge_ComputedCorrectly(t *testing.T) {
	base := time.Now()
	tr := New(fixedNow(base))
	tr.Observe([]int{22})
	tr.now = fixedNow(base.Add(10 * time.Second))
	tr.Observe([]int{22})
	e, _ := tr.Get(22)
	if e.Age() != 10*time.Second {
		t.Errorf("Age = %v; want 10s", e.Age())
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	tr := New(nil)
	tr.Observe([]int{8080})
	tr.Remove(8080)
	_, ok := tr.Get(8080)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	tr := New(nil)
	_, ok := tr.Get(9999)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestOlderThan_ReturnsMatchingEntries(t *testing.T) {
	base := time.Now()
	tr := New(fixedNow(base))
	tr.Observe([]int{80, 443})
	tr.now = fixedNow(base.Add(30 * time.Second))
	tr.Observe([]int{80, 443})
	results := tr.OlderThan(20 * time.Second)
	if len(results) != 2 {
		t.Errorf("OlderThan returned %d entries; want 2", len(results))
	}
}

func TestOlderThan_ExcludesRecentEntries(t *testing.T) {
	base := time.Now()
	tr := New(fixedNow(base))
	tr.Observe([]int{8080})
	tr.now = fixedNow(base.Add(5 * time.Second))
	tr.Observe([]int{8080})
	results := tr.OlderThan(10 * time.Second)
	if len(results) != 0 {
		t.Errorf("expected no results; got %d", len(results))
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := New(nil)
	tr.Observe([]int{22, 80, 443})
	all := tr.All()
	if len(all) != 3 {
		t.Errorf("All returned %d entries; want 3", len(all))
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	tr := New(nil)
	tr.Observe([]int{22, 80})
	tr.Reset()
	if len(tr.All()) != 0 {
		t.Error("expected empty tracker after Reset")
	}
}
