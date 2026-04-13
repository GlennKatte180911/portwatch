package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected Open, got %v", b.State())
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestRecordSuccess_ResetsToClosed(t *testing.T) {
	b := circuitbreaker.New(2, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected Open state")
	}
	// simulate recovery window passing
	b2 := circuitbreaker.New(2, time.Millisecond)
	b2.RecordFailure()
	b2.RecordFailure()
	time.Sleep(5 * time.Millisecond)
	if err := b2.Allow(); err != nil {
		t.Fatalf("expected HalfOpen to allow, got %v", err)
	}
	b2.RecordSuccess()
	if b2.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed after success, got %v", b2.State())
	}
}

func TestHalfOpen_TransitionOnTimeout(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected Open")
	}
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected HalfOpen to allow, got %v", err)
	}
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", b.State())
	}
}

func TestHalfOpen_FailureReopens(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = b.Allow() // transitions to HalfOpen
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected Open after HalfOpen failure, got %v", b.State())
	}
}
