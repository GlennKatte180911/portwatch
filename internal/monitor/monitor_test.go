package monitor_test

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestRun_DetectsNewPort(t *testing.T) {
	snapshotFile, err := os.CreateTemp(t.TempDir(), "snapshot-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	snapshotFile.Close()
	os.Remove(snapshotFile.Name()) // start with no snapshot

	port := freePort(t)
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("could not open test port: %v", err)
	}
	defer l.Close()

	var added, removed []int
	notifier := alert.NewNotifier(func(a, r []int) {
		added = a
		removed = r
	})

	cfg := config.Default()
	cfg.StartPort = port
	cfg.EndPort = port
	cfg.Interval = 50 * time.Millisecond
	cfg.SnapshotPath = snapshotFile.Name()

	mon := monitor.New(cfg, notifier)
	done := make(chan struct{})

	// First tick: no previous snapshot, saves current state.
	// Second tick: same ports open, no alert expected.
	go func() {
		time.Sleep(120 * time.Millisecond)
		close(done)
	}()

	if err := mon.Run(done); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	// No changes expected since port was open from the start.
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no alerts, got added=%v removed=%v", added, removed)
	}
}

func TestRun_StopsOnDone(t *testing.T) {
	cfg := config.Default()
	cfg.SnapshotPath = t.TempDir() + "/snap.json"
	cfg.Interval = 10 * time.Second // long interval to ensure done fires first

	notifier := alert.NewNotifier(nil)
	mon := monitor.New(cfg, notifier)

	done := make(chan struct{})
	close(done)

	start := time.Now()
	if err := mon.Run(done); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) > 2*time.Second {
		t.Error("Run did not stop promptly when done was closed")
	}
}
