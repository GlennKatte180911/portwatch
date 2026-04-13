package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestBeat_UpdatesLastBeat(t *testing.T) {
	w := watchdog.New(time.Second, time.Second, nil)
	before := w.LastBeat()
	time.Sleep(5 * time.Millisecond)
	w.Beat()
	after := w.LastBeat()
	if !after.After(before) {
		t.Errorf("expected LastBeat to advance after Beat(); got before=%v after=%v", before, after)
	}
}

func TestRun_CallsOnStall_WhenNoBeats(t *testing.T) {
	var called int32
	var stalledAt time.Time

	onStall := func(last time.Time) {
		atomic.StoreInt32(&called, 1)
		stalledAt = last
	}

	// Very short interval so the test completes quickly.
	w := watchdog.New(10*time.Millisecond, 10*time.Millisecond, onStall)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go w.Run(ctx)

	// Wait long enough for the watchdog to fire.
	time.Sleep(150 * time.Millisecond)

	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected onStall to be called when no beats are sent")
	}
	if stalledAt.IsZero() {
		t.Error("expected stalledAt to be non-zero")
	}
}

func TestRun_DoesNotCallOnStall_WhenBeatsAreSent(t *testing.T) {
	var called int32

	onStall := func(_ time.Time) {
		atomic.StoreInt32(&called, 1)
	}

	w := watchdog.New(30*time.Millisecond, 30*time.Millisecond, onStall)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go w.Run(ctx)

	// Keep sending beats faster than the watchdog deadline.
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	done := time.After(180 * time.Millisecond)
loop:
	for {
		select {
		case <-ticker.C:
			w.Beat()
		case <-done:
			break loop
		}
	}

	if atomic.LoadInt32(&called) != 0 {
		t.Error("expected onStall NOT to be called while regular beats are sent")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	w := watchdog.New(50*time.Millisecond, 50*time.Millisecond, nil)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	cancel()

	select {
		case <-done:
		// ok
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Run did not return after context cancellation")
	}
}
