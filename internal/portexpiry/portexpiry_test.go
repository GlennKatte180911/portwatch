package portexpiry

import (
	"testing"
	"time"
)

func TestObserve_SetsFirstSeen(t *testing.T) {
	e := New(time.Hour)
	before := time.Now()
	e.Observe(8080)
	age, ok := e.Age(8080)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age < 0 || age > time.Since(before)+time.Millisecond {
		t.Errorf("unexpected age: %v", age)
	}
}

func TestObserve_IdempotentFirstSeen(t *testing.T) {
	e := New(time.Hour)
	e.Observe(9090)
	age1, _ := e.Age(9090)
	time.Sleep(5 * time.Millisecond)
	e.Observe(9090)
	age2, _ := e.Age(9090)
	if age2 < age1 {
		t.Error("second Observe should not reset FirstSeen")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	e := New(time.Hour)
	e.Observe(443)
	e.Remove(443)
	_, ok := e.Age(443)
	if ok {
		t.Error("expected port to be removed")
	}
}

func TestExpired_ReturnsOverAgeEntries(t *testing.T) {
	e := New(10 * time.Millisecond)
	e.Observe(22)
	e.Observe(80)
	time.Sleep(20 * time.Millisecond)
	expired := e.Expired()
	if len(expired) != 2 {
		t.Fatalf("expected 2 expired ports, got %d", len(expired))
	}
}

func TestExpired_ExcludesRecentEntries(t *testing.T) {
	e := New(time.Hour)
	e.Observe(3000)
	expired := e.Expired()
	if len(expired) != 0 {
		t.Errorf("expected no expired ports, got %d", len(expired))
	}
}

func TestAge_UnknownPort_ReturnsFalse(t *testing.T) {
	e := New(time.Hour)
	_, ok := e.Age(9999)
	if ok {
		t.Error("expected ok=false for unknown port")
	}
}
