package semaphore_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/semaphore"
)

func TestNew_InvalidCapacity_ReturnsError(t *testing.T) {
	_, err := semaphore.New(0)
	if err == nil {
		t.Fatal("expected error for capacity 0, got nil")
	}
}

func TestNew_ValidCapacity_ReturnsSemaphore(t *testing.T) {
	s, err := semaphore.New(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Cap() != 4 {
		t.Fatalf("expected cap 4, got %d", s.Cap())
	}
}

func TestAcquireRelease_UpdatesAvailable(t *testing.T) {
	s, _ := semaphore.New(3)
	ctx := context.Background()

	if err := s.Acquire(ctx); err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if s.Available() != 2 {
		t.Fatalf("expected 2 available, got %d", s.Available())
	}
	s.Release()
	if s.Available() != 3 {
		t.Fatalf("expected 3 available after release, got %d", s.Available())
	}
}

func TestAcquire_BlocksAtCapacity(t *testing.T) {
	s, _ := semaphore.New(2)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = s.Acquire(context.Background())
	_ = s.Acquire(context.Background())

	err := s.Acquire(ctx)
	if err == nil {
		t.Fatal("expected context deadline error when semaphore is full")
	}
}

func TestAcquire_ContextCancelled_ReturnsError(t *testing.T) {
	s, _ := semaphore.New(1)
	_ = s.Acquire(context.Background()) // fill the slot

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := s.Acquire(ctx); err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestSemaphore_BoundsConcurrency(t *testing.T) {
	const cap = 3
	s, _ := semaphore.New(cap)

	var active atomic.Int64
	var maxSeen atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Acquire(context.Background())
			defer s.Release()

			cur := active.Add(1)
			for {
				prev := maxSeen.Load()
				if cur <= prev || maxSeen.CompareAndSwap(prev, cur) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			active.Add(-1)
		}()
	}
	wg.Wait()

	if got := maxSeen.Load(); got > int64(cap) {
		t.Fatalf("concurrency exceeded cap: got %d, want <= %d", got, cap)
	}
}
