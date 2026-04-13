package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/ratelimit"
)

func TestAllow_FirstEventAlwaysPasses(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	if !l.Allow(8080) {
		t.Fatal("expected first event for port 8080 to be allowed")
	}
}

func TestAllow_SecondEventBlockedWithinCooldown(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	l.Allow(8080)
	if l.Allow(8080) {
		t.Fatal("expected second event within cooldown to be blocked")
	}
}

func TestAllow_EventAllowedAfterCooldown(t *testing.T) {
	l := ratelimit.New(20 * time.Millisecond)
	l.Allow(9090)
	time.Sleep(30 * time.Millisecond)
	if !l.Allow(9090) {
		t.Fatal("expected event to be allowed after cooldown elapsed")
	}
}

func TestAllow_IndependentPerPort(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	l.Allow(1000)
	if !l.Allow(2000) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	l.Allow(3000)
	l.Reset(3000)
	if !l.Allow(3000) {
		t.Fatal("expected Allow to pass after Reset")
	}
}

func TestResetAll_ClearsAllPorts(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	l.Allow(4000)
	l.Allow(5000)
	l.ResetAll()
	if !l.Allow(4000) {
		t.Fatal("expected port 4000 to be allowed after ResetAll")
	}
	if !l.Allow(5000) {
		t.Fatal("expected port 5000 to be allowed after ResetAll")
	}
}
