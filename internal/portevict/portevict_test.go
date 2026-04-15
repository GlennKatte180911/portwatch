package portevict_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portevict"
)

func TestMark_TracksPendingPort(t *testing.T) {
	e := portevict.New(10 * time.Second)
	e.Mark(8080)
	if got := e.Pending(); got != 1 {
		t.Fatalf("expected 1 pending, got %d", got)
	}
}

func TestMark_IdempotentForSamePort(t *testing.T) {
	e := portevict.New(10 * time.Second)
	e.Mark(8080)
	e.Mark(8080)
	if got := e.Pending(); got != 1 {
		t.Fatalf("expected 1 pending after duplicate mark, got %d", got)
	}
}

func TestLift_RemovesPort(t *testing.T) {
	e := portevict.New(10 * time.Second)
	e.Mark(9090)
	e.Lift(9090)
	if got := e.Pending(); got != 0 {
		t.Fatalf("expected 0 pending after lift, got %d", got)
	}
}

func TestConfirmed_EmptyBeforeGraceElapsed(t *testing.T) {
	e := portevict.New(10 * time.Second)
	e.Mark(443)
	confirmed := e.Confirmed()
	if len(confirmed) != 0 {
		t.Fatalf("expected no confirmed ports within grace period, got %v", confirmed)
	}
}

func TestConfirmed_ReturnsPortAfterGraceElapsed(t *testing.T) {
	e := portevict.New(10 * time.Millisecond)
	e.Mark(22)
	time.Sleep(30 * time.Millisecond)
	confirmed := e.Confirmed()
	if len(confirmed) != 1 || confirmed[0] != 22 {
		t.Fatalf("expected port 22 confirmed, got %v", confirmed)
	}
	if got := e.Pending(); got != 0 {
		t.Fatalf("expected 0 pending after confirmation, got %d", got)
	}
}

func TestConfirmed_PartialExpiry(t *testing.T) {
	e := portevict.New(20 * time.Millisecond)
	e.Mark(80)
	time.Sleep(30 * time.Millisecond)
	e.Mark(443) // marked later, still in grace
	confirmed := e.Confirmed()
	if len(confirmed) != 1 || confirmed[0] != 80 {
		t.Fatalf("expected only port 80 confirmed, got %v", confirmed)
	}
	if got := e.Pending(); got != 1 {
		t.Fatalf("expected 1 still pending, got %d", got)
	}
}
