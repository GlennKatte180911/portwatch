package portsampler_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portsampler"
)

func TestObserve_TracksCount(t *testing.T) {
	s := portsampler.New(time.Minute, 2)
	s.Observe([]int{80, 443})
	s.Observe([]int{80})

	sm, ok := s.Get(80)
	if !ok {
		t.Fatal("expected sample for port 80")
	}
	if sm.Count != 2 {
		t.Fatalf("expected count 2, got %d", sm.Count)
	}

	sm443, ok := s.Get(443)
	if !ok {
		t.Fatal("expected sample for port 443")
	}
	if sm443.Count != 1 {
		t.Fatalf("expected count 1, got %d", sm443.Count)
	}
}

func TestStable_ReturnsPortsAtThreshold(t *testing.T) {
	s := portsampler.New(time.Minute, 3)
	for i := 0; i < 3; i++ {
		s.Observe([]int{8080})
	}
	s.Observe([]int{9090}) // only once — below threshold

	stable := s.Stable()
	if len(stable) != 1 || stable[0] != 8080 {
		t.Fatalf("expected [8080], got %v", stable)
	}
}

func TestStable_EmptyWhenNoneReachThreshold(t *testing.T) {
	s := portsampler.New(time.Minute, 5)
	s.Observe([]int{22, 80})

	if got := s.Stable(); len(got) != 0 {
		t.Fatalf("expected empty stable list, got %v", got)
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	s := portsampler.New(time.Minute, 1)
	_, ok := s.Get(12345)
	if ok {
		t.Fatal("expected false for unknown port")
	}
}

func TestObserve_EvictsExpiredSamples(t *testing.T) {
	s := portsampler.New(50*time.Millisecond, 1)
	s.Observe([]int{3306})

	time.Sleep(80 * time.Millisecond)

	// Trigger eviction via a new observation.
	s.Observe([]int{5432})

	if _, ok := s.Get(3306); ok {
		t.Fatal("expected port 3306 to be evicted after window expiry")
	}
	if _, ok := s.Get(5432); !ok {
		t.Fatal("expected port 5432 to be present")
	}
}

func TestReset_ClearsAllSamples(t *testing.T) {
	s := portsampler.New(time.Minute, 1)
	s.Observe([]int{80, 443, 8080})
	s.Reset()

	if _, ok := s.Get(80); ok {
		t.Fatal("expected no samples after Reset")
	}
	if got := s.Stable(); len(got) != 0 {
		t.Fatalf("expected empty stable list after Reset, got %v", got)
	}
}

func TestNew_ThresholdBelowOne_ClampedToOne(t *testing.T) {
	s := portsampler.New(time.Minute, 0)
	s.Observe([]int{80})

	if stable := s.Stable(); len(stable) != 1 {
		t.Fatalf("expected threshold clamped to 1, got stable=%v", stable)
	}
}
