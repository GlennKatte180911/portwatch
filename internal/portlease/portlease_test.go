package portlease

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestGrant_StoresLease(t *testing.T) {
	r := New()
	if err := r.Grant(8080, time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Active(8080) {
		t.Error("expected port 8080 to have an active lease")
	}
}

func TestGrant_InvalidPort_ReturnsError(t *testing.T) {
	r := New()
	if err := r.Grant(0, time.Minute); err == nil {
		t.Error("expected error for port 0")
	}
	if err := r.Grant(65536, time.Minute); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestGrant_ZeroDuration_ReturnsError(t *testing.T) {
	r := New()
	if err := r.Grant(443, 0); err == nil {
		t.Error("expected error for zero duration")
	}
}

func TestActive_ReturnsFalseForUnknownPort(t *testing.T) {
	r := New()
	if r.Active(9999) {
		t.Error("expected inactive for unknown port")
	}
}

func TestActive_ReturnsFalseAfterExpiry(t *testing.T) {
	base := time.Now()
	r := New()
	r.now = fixedClock(base)
	_ = r.Grant(8080, time.Second)
	r.now = fixedClock(base.Add(2 * time.Second))
	if r.Active(8080) {
		t.Error("expected lease to be expired")
	}
}

func TestRevoke_RemovesLease(t *testing.T) {
	r := New()
	_ = r.Grant(8080, time.Minute)
	r.Revoke(8080)
	if r.Active(8080) {
		t.Error("expected lease to be revoked")
	}
}

func TestExpired_ReturnsExpiredLeases(t *testing.T) {
	base := time.Now()
	r := New()
	r.now = fixedClock(base)
	_ = r.Grant(80, time.Second)
	_ = r.Grant(443, time.Minute)
	r.now = fixedClock(base.Add(2 * time.Second))
	expired := r.Expired()
	if len(expired) != 1 || expired[0].Port != 80 {
		t.Errorf("expected only port 80 expired, got %v", expired)
	}
}

func TestPurge_RemovesExpiredLeases(t *testing.T) {
	base := time.Now()
	r := New()
	r.now = fixedClock(base)
	_ = r.Grant(80, time.Second)
	_ = r.Grant(443, time.Minute)
	r.now = fixedClock(base.Add(2 * time.Second))
	n := r.Purge()
	if n != 1 {
		t.Errorf("expected 1 purged, got %d", n)
	}
	if r.Active(80) {
		t.Error("expected port 80 to be purged")
	}
	if !r.Active(443) {
		t.Error("expected port 443 to still be active")
	}
}
