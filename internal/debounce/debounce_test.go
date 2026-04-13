package debounce_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

const shortWait = 50 * time.Millisecond

func TestPush_EmitsAfterQuietPeriod(t *testing.T) {
	d, ch := debounce.New(shortWait, 4)
	defer d.Stop()

	d.Push("8080")

	select {
	case ev := <-ch:
		if ev.Key != "8080" {
			t.Fatalf("expected key 8080, got %s", ev.Key)
		}
	case <-time.After(shortWait * 5):
		t.Fatal("timed out waiting for debounced event")
	}
}

func TestPush_ResetsTimerOnRapidEvents(t *testing.T) {
	d, ch := debounce.New(shortWait, 4)
	defer d.Stop()

	// Push three times in quick succession; only one event should be emitted.
	d.Push("9090")
	time.Sleep(shortWait / 4)
	d.Push("9090")
	time.Sleep(shortWait / 4)
	d.Push("9090")

	// Collect events for 3× the wait window.
	var count int
	timeout := time.After(shortWait * 6)
loop:
	for {
		select {
		case <-ch:
			count++
		case <-timeout:
			break loop
		}
	}

	if count != 1 {
		t.Fatalf("expected exactly 1 event, got %d", count)
	}
}

func TestPush_IndependentKeysEachEmit(t *testing.T) {
	d, ch := debounce.New(shortWait, 8)
	defer d.Stop()

	d.Push("80")
	d.Push("443")

	received := make(map[string]bool)
	timeout := time.After(shortWait * 5)
loop:
	for len(received) < 2 {
		select {
		case ev := <-ch:
			received[ev.Key] = true
		case <-timeout:
			break loop
		}
	}

	if !received["80"] || !received["443"] {
		t.Fatalf("expected events for both keys, got %v", received)
	}
}

func TestStop_CancelsPendingTimers(t *testing.T) {
	d, ch := debounce.New(shortWait*10, 4)

	d.Push("3000")
	d.Stop()

	// Channel should be closed with no event emitted.
	select {
	case ev, ok := <-ch:
		if ok {
			t.Fatalf("unexpected event after Stop: %v", ev)
		}
		// closed — expected
	case <-time.After(shortWait * 3):
		t.Fatal("channel was not closed after Stop")
	}
}
