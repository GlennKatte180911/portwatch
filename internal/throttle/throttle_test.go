package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_FirstEventAlwaysPasses(t *testing.T) {
	th := throttle.New(time.Second, 3)
	if !th.Allow("port:80") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_BlocksAfterMaxReached(t *testing.T) {
	th := throttle.New(time.Second, 2)
	key := "port:443"

	if !th.Allow(key) {
		t.Fatal("expected event 1 to be allowed")
	}
	if !th.Allow(key) {
		t.Fatal("expected event 2 to be allowed")
	}
	if th.Allow(key) {
		t.Fatal("expected event 3 to be blocked")
	}
}

func TestAllow_AllowsAgainAfterWindowExpires(t *testing.T) {
	th := throttle.New(50*time.Millisecond, 1)
	key := "port:8080"

	if !th.Allow(key) {
		t.Fatal("expected first event to pass")
	}
	if th.Allow(key) {
		t.Fatal("expected second event to be blocked")
	}

	time.Sleep(60 * time.Millisecond)

	if !th.Allow(key) {
		t.Fatal("expected event to pass after window expired")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := throttle.New(time.Second, 1)

	if !th.Allow("port:80") {
		t.Fatal("expected port:80 to be allowed")
	}
	if !th.Allow("port:443") {
		t.Fatal("expected port:443 to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	th := throttle.New(time.Second, 1)
	key := "port:22"

	th.Allow(key)
	if th.Allow(key) {
		t.Fatal("expected key to be throttled before reset")
	}

	th.Reset(key)
	if !th.Allow(key) {
		t.Fatal("expected key to be allowed after reset")
	}
}

func TestRemaining_DecreasesWithEvents(t *testing.T) {
	th := throttle.New(time.Second, 3)
	key := "port:9000"

	if got := th.Remaining(key); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}

	th.Allow(key)
	if got := th.Remaining(key); got != 2 {
		t.Fatalf("expected 2 remaining, got %d", got)
	}

	th.Allow(key)
	th.Allow(key) // hits limit
	th.Allow(key) // blocked
	if got := th.Remaining(key); got != 0 {
		t.Fatalf("expected 0 remaining, got %d", got)
	}
}
