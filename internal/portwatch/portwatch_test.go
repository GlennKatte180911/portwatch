package portwatch_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portdiff"
	"github.com/user/portwatch/internal/portwatch"
)

// fakeNotifier captures the last diff it received.
type fakeNotifier struct {
	events []portdiff.Diff
}

func (f *fakeNotifier) Notify(_ context.Context, d portdiff.Diff) error {
	f.events = append(f.events, d)
	return nil
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestRun_DetectsOpenPort(t *testing.T) {
	port := freePort(t)
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	fn := &fakeNotifier{}
	w, err := portwatch.New(portwatch.Config{
		StartPort: port,
		EndPort:   port,
		Interval:  20 * time.Millisecond,
		Notifier:  fn,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Run(ctx) //nolint:errcheck

	if len(fn.events) == 0 {
		t.Fatal("expected at least one diff event")
	}
	found := false
	for _, p := range fn.events[0].Added {
		if p == port {
			found = true
		}
	}
	if !found {
		t.Errorf("port %d not in first Added diff", port)
	}
}

func TestRun_StopsOnCancel(t *testing.T) {
	fn := &fakeNotifier{}
	w, err := portwatch.New(portwatch.Config{
		StartPort: 19000,
		EndPort:   19001,
		Interval:  50 * time.Millisecond,
		Notifier:  fn,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := w.Run(ctx); err == nil {
		t.Error("expected non-nil error on cancelled context")
	}
}
