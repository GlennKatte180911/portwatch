package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retry"
)

var errTemp = errors.New("temporary error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	r := retry.New(retry.DefaultConfig())
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUntilSuccess(t *testing.T) {
	cfg := retry.Config{MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}
	r := retry.New(cfg)
	var calls int32
	err := r.Do(context.Background(), func() error {
		if atomic.AddInt32(&calls, 1) < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after eventual success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ReturnsErrMaxAttemptsWhenExhausted(t *testing.T) {
	cfg := retry.Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
	r := retry.New(cfg)
	var calls int32
	err := r.Do(context.Background(), func() error {
		atomic.AddInt32(&calls, 1)
		return errTemp
	})
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	cfg := retry.Config{MaxAttempts: 10, BaseDelay: 100 * time.Millisecond, MaxDelay: time.Second}
	r := retry.New(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	var calls int32
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()
	err := r.Do(ctx, func() error {
		atomic.AddInt32(&calls, 1)
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls == 0 {
		t.Fatal("expected at least one call before cancellation")
	}
}

func TestDefaultConfig_Values(t *testing.T) {
	cfg := retry.DefaultConfig()
	if cfg.MaxAttempts != 4 {
		t.Errorf("MaxAttempts: want 4, got %d", cfg.MaxAttempts)
	}
	if cfg.BaseDelay != 250*time.Millisecond {
		t.Errorf("BaseDelay: want 250ms, got %v", cfg.BaseDelay)
	}
	if cfg.MaxDelay != 10*time.Second {
		t.Errorf("MaxDelay: want 10s, got %v", cfg.MaxDelay)
	}
}
