package porttrend_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/porttrend"
)

func TestRecord_TracksSeenCount(t *testing.T) {
	tr := porttrend.New()
	now := time.Now()

	tr.Record([]int{80, 443}, now)
	tr.Record([]int{80}, now.Add(time.Second))

	e, ok := tr.Get(80)
	if !ok {
		t.Fatal("expected entry for port 80")
	}
	if e.SeenCount != 2 {
		t.Errorf("SeenCount = %d, want 2", e.SeenCount)
	}

	e443, ok := tr.Get(443)
	if !ok {
		t.Fatal("expected entry for port 443")
	}
	if e443.SeenCount != 1 {
		t.Errorf("SeenCount = %d, want 1", e443.SeenCount)
	}
}

func TestRecord_SetsFirstAndLastSeen(t *testing.T) {
	tr := porttrend.New()
	first := time.Now()
	second := first.Add(5 * time.Second)

	tr.Record([]int{8080}, first)
	tr.Record([]int{8080}, second)

	e, _ := tr.Get(8080)
	if !e.FirstSeen.Equal(first) {
		t.Errorf("FirstSeen = %v, want %v", e.FirstSeen, first)
	}
	if !e.LastSeen.Equal(second) {
		t.Errorf("LastSeen = %v, want %v", e.LastSeen, second)
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	tr := porttrend.New()
	_, ok := tr.Get(9999)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := porttrend.New()
	now := time.Now()
	tr.Record([]int{22, 80, 443}, now)

	all := tr.All()
	if len(all) != 3 {
		t.Errorf("len(All()) = %d, want 3", len(all))
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	tr := porttrend.New()
	tr.Record([]int{22, 80}, time.Now())
	tr.Reset()

	if entries := tr.All(); len(entries) != 0 {
		t.Errorf("expected empty after Reset, got %d entries", len(entries))
	}

	if _, ok := tr.Get(22); ok {
		t.Error("expected Get to return false after Reset")
	}
}
