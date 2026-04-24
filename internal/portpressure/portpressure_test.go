package portpressure_test

import (
	"testing"

	"github.com/user/portwatch/internal/portpressure"
)

func TestObserve_IncreasesScoreForSeenPort(t *testing.T) {
	tr := portpressure.New(0.5)
	tr.Observe([]int{8080})

	e, ok := tr.Get(8080)
	if !ok {
		t.Fatal("expected entry for port 8080")
	}
	if e.Score <= 0 {
		t.Errorf("expected positive score, got %f", e.Score)
	}
}

func TestObserve_DecreasesScoreForAbsentPort(t *testing.T) {
	tr := portpressure.New(0.5)
	// First observation establishes a score.
	tr.Observe([]int{9090})
	initial, _ := tr.Get(9090)

	// Second observation without the port should lower the score.
	tr.Observe([]int{})
	after, _ := tr.Get(9090)

	if after.Score >= initial.Score {
		t.Errorf("expected score to decrease: initial=%f after=%f", initial.Score, after.Score)
	}
}

func TestObserve_ScanCountIncrementsOnlyWhenPresent(t *testing.T) {
	tr := portpressure.New(0.3)
	tr.Observe([]int{443})
	tr.Observe([]int{443})
	tr.Observe([]int{}) // absent

	e, ok := tr.Get(443)
	if !ok {
		t.Fatal("expected entry for port 443")
	}
	if e.ScanCount != 2 {
		t.Errorf("expected ScanCount=2, got %d", e.ScanCount)
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	tr := portpressure.New(0.3)
	_, ok := tr.Get(12345)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestAbove_ReturnsPortsOverThreshold(t *testing.T) {
	tr := portpressure.New(0.9)
	// Drive port 80 to a high score with repeated observations.
	for i := 0; i < 10; i++ {
		tr.Observe([]int{80})
	}
	// Port 22 observed only once — low score.
	tr.Observe([]int{22})

	high := tr.Above(0.8)
	found := false
	for _, e := range high {
		if e.Port == 80 {
			found = true
		}
		if e.Port == 22 {
			t.Error("port 22 should not exceed threshold")
		}
	}
	if !found {
		t.Error("expected port 80 to be above threshold")
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	tr := portpressure.New(0.3)
	tr.Observe([]int{80, 443, 8080})
	tr.Reset()

	_, ok := tr.Get(80)
	if ok {
		t.Error("expected no entry after reset")
	}
	if above := tr.Above(0); len(above) != 0 {
		t.Errorf("expected empty result after reset, got %d entries", len(above))
	}
}

func TestNew_InvalidAlpha_UsesDefault(t *testing.T) {
	// Should not panic; invalid alpha falls back to 0.3.
	tr := portpressure.New(-1)
	tr.Observe([]int{53})
	e, ok := tr.Get(53)
	if !ok || e.Score <= 0 {
		t.Error("expected valid entry with default alpha")
	}
}
