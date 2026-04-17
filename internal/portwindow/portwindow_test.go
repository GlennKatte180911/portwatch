package portwindow_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwindow"
)

func TestObserve_CountIncreases(t *testing.T) {
	w := portwindow.New(5 * time.Second)
	w.Observe(8080)
	w.Observe(8080)
	if got := w.Count(8080); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_UnknownPort_ReturnsZero(t *testing.T) {
	w := portwindow.New(5 * time.Second)
	if got := w.Count(9999); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestObserve_ExpiredEntriesNotCounted(t *testing.T) {
	w := portwindow.New(50 * time.Millisecond)
	w.Observe(443)
	time.Sleep(80 * time.Millisecond)
	if got := w.Count(443); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestActive_ReturnsObservedPorts(t *testing.T) {
	w := portwindow.New(5 * time.Second)
	w.Observe(80)
	w.Observe(443)
	w.Observe(80)
	active := w.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active ports, got %d", len(active))
	}
}

func TestActive_EmptyAfterExpiry(t *testing.T) {
	w := portwindow.New(50 * time.Millisecond)
	w.Observe(8080)
	time.Sleep(80 * time.Millisecond)
	if got := w.Active(); len(got) != 0 {
		t.Fatalf("expected empty active list, got %v", got)
	}
}

func TestReset_ClearsAllObservations(t *testing.T) {
	w := portwindow.New(5 * time.Second)
	w.Observe(22)
	w.Observe(80)
	w.Reset()
	if got := w.Count(22); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if active := w.Active(); len(active) != 0 {
		t.Fatalf("expected empty after reset, got %v", active)
	}
}
