package schedule_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

func TestNew_TicksAtInterval(t *testing.T) {
	s := schedule.New(20*time.Millisecond, 0)
	defer s.Stop()

	select {
	case <-s.C:
		// received a tick — pass
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected tick within 200ms")
	}
}

func TestStop_HaltsEmission(t *testing.T) {
	s := schedule.New(20*time.Millisecond, 0)
	s.Stop()

	// drain any tick that may have fired before Stop
	time.Sleep(10 * time.Millisecond)
	for len(s.C) > 0 {
		<-s.C
	}

	time.Sleep(60 * time.Millisecond)
	if len(s.C) > 0 {
		t.Fatal("scheduler emitted after Stop")
	}
}

func TestStop_Idempotent(t *testing.T) {
	s := schedule.New(50*time.Millisecond, 0)
	s.Stop()
	s.Stop() // must not panic
}

func TestRunWithContext_CallsFn(t *testing.T) {
	s := schedule.New(20*time.Millisecond, 0)
	defer s.Stop()

	var count atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	s.RunWithContext(ctx, func() { count.Add(1) })

	if count.Load() == 0 {
		t.Fatal("fn was never called")
	}
}

func TestRunWithContext_StopsOnCancel(t *testing.T) {
	s := schedule.New(10*time.Millisecond, 0)
	defer s.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		s.RunWithContext(ctx, func() {})
		close(done)
	}()

	select {
	case <-done:
		// returned promptly
	case <-time.After(200 * time.Millisecond):
		t.Fatal("RunWithContext did not return after cancel")
	}
}

func TestJitter_DoesNotPanic(t *testing.T) {
	s := schedule.New(20*time.Millisecond, 5*time.Millisecond)
	defer s.Stop()

	select {
	case <-s.C:
		// received tick with jitter enabled
	case <-time.After(300 * time.Millisecond):
		t.Fatal("expected tick within 300ms with jitter")
	}
}
