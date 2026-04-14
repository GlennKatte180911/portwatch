package tunables_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/tunables"
)

func TestDefaults_ReturnsExpectedValues(t *testing.T) {
	tn := tunables.Defaults()

	if tn.ScanInterval() != 30*time.Second {
		t.Errorf("expected 30s scan interval, got %v", tn.ScanInterval())
	}
	if tn.ScanTimeout() != 2*time.Second {
		t.Errorf("expected 2s scan timeout, got %v", tn.ScanTimeout())
	}
	if tn.AlertCooldown() != 5*time.Minute {
		t.Errorf("expected 5m alert cooldown, got %v", tn.AlertCooldown())
	}
	if tn.MaxAlerts() != 10 {
		t.Errorf("expected max alerts 10, got %d", tn.MaxAlerts())
	}
}

func TestSetScanInterval_UpdatesValue(t *testing.T) {
	tn := tunables.Defaults()
	if err := tn.SetScanInterval(1 * time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tn.ScanInterval() != 1*time.Minute {
		t.Errorf("expected 1m, got %v", tn.ScanInterval())
	}
}

func TestSetScanInterval_ZeroReturnsError(t *testing.T) {
	tn := tunables.Defaults()
	if err := tn.SetScanInterval(0); err == nil {
		t.Fatal("expected error for zero interval, got nil")
	}
}

func TestSetScanTimeout_NegativeReturnsError(t *testing.T) {
	tn := tunables.Defaults()
	if err := tn.SetScanTimeout(-1 * time.Second); err == nil {
		t.Fatal("expected error for negative timeout, got nil")
	}
}

func TestSetAlertCooldown_ZeroIsAllowed(t *testing.T) {
	tn := tunables.Defaults()
	if err := tn.SetAlertCooldown(0); err != nil {
		t.Fatalf("zero cooldown should be allowed, got: %v", err)
	}
	if tn.AlertCooldown() != 0 {
		t.Errorf("expected 0, got %v", tn.AlertCooldown())
	}
}

func TestSetAlertCooldown_NegativeReturnsError(t *testing.T) {
	tn := tunables.Defaults()
	if err := tn.SetAlertCooldown(-time.Second); err == nil {
		t.Fatal("expected error for negative cooldown, got nil")
	}
}

func TestSetMaxAlerts_UpdatesValue(t *testing.T) {
	tn := tunables.Defaults()
	if err := tn.SetMaxAlerts(25); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tn.MaxAlerts() != 25 {
		t.Errorf("expected 25, got %d", tn.MaxAlerts())
	}
}

func TestSetMaxAlerts_ZeroReturnsError(t *testing.T) {
	tn := tunables.Defaults()
	if err := tn.SetMaxAlerts(0); err == nil {
		t.Fatal("expected error for zero max alerts, got nil")
	}
}

func TestConcurrentAccess_DoesNotRace(t *testing.T) {
	tn := tunables.Defaults()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			_ = tn.SetScanInterval(time.Duration(i+1) * time.Millisecond)
		}
		close(done)
	}()
	for i := 0; i < 100; i++ {
		_ = tn.ScanInterval()
	}
	<-done
}
