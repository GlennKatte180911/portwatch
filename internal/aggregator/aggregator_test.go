package aggregator_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/aggregator"
	"github.com/user/portwatch/internal/snapshot"
)

func TestPush_BatchesWithinWindow(t *testing.T) {
	a := aggregator.New(50 * time.Millisecond)
	defer a.Stop()

	a.Push(snapshot.Diff{Added: []int{8080}})
	a.Push(snapshot.Diff{Added: []int{9090}})

	select {
	case ev := <-a.Events():
		if len(ev.Added) != 2 {
			t.Fatalf("expected 2 added ports, got %d", len(ev.Added))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for aggregated event")
	}
}

func TestPush_CancelsOpposites(t *testing.T) {
	a := aggregator.New(50 * time.Millisecond)
	defer a.Stop()

	a.Push(snapshot.Diff{Added: []int{8080}})
	a.Push(snapshot.Diff{Removed: []int{8080}})

	select {
	case ev := <-a.Events():
		t.Fatalf("expected no event, got added=%v removed=%v", ev.Added,case <-time.After(120 * time.Millisecond):
		// correct: net change is zero, nothing emitted
	}
}

func TestStop_FlushesPending(t *testing.T) {
	a := aggregator.New(10 * time.Second) // long window

	a.Push(snapshot.Diff{Added: []int{3000}})
	a.Stop()

	select {
	case ev, ok := <-a.Events():
		if !ok {
			t.Fatal("channel closed before event was delivered")
		}
		if len(ev.Added) != 1 || ev.Added[0] != 3000 {
			t.Fatalf("unexpected event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for flush on Stop")
	}
}

func TestEvents_EmptyDiff_NotEmitted(t *testing.T) {
	a := aggregator.New(40 * time.Millisecond)
	defer a.Stop()

	// push an empty diff
	a.Push(snapshot.Diff{})

	select {
	case ev := <-a.Events():
		t.Fatalf("expected no event for empty diff, got %+v", ev)
	case <-time.After(120 * time.Millisecond):
		// correct
	}
}

func TestEvent_TimestampSet(t *testing.T) {
	a := aggregator.New(40 * time.Millisecond)
	defer a.Stop()

	before := time.Now()
	a.Push(snapshot.Diff{Removed: []int{22}})

	select {
	case ev := <-a.Events():
		if ev.At.Before(before) {
			t.Fatalf("event timestamp %v is before push time %v", ev.At, before)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out")
	}
}
