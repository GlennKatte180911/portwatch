package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNew_InitialisesCollector(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.TotalScans != 0 {
		t.Fatalf("expected 0 scans, got %d", s.TotalScans)
	}
	if s.TotalAlerts != 0 {
		t.Fatalf("expected 0 alerts, got %d", s.TotalAlerts)
	}
	if s.UptimeSeconds < 0 {
		t.Fatal("uptime should be non-negative")
	}
}

func TestRecordScan_IncrementsCounter(t *testing.T) {
	c := metrics.New()
	c.RecordScan(10 * time.Millisecond)
	c.RecordScan(20 * time.Millisecond)
	s := c.Snapshot()
	if s.TotalScans != 2 {
		t.Fatalf("expected 2, got %d", s.TotalScans)
	}
}

func TestRecordScan_UpdatesLastScanAt(t *testing.T) {
	c := metrics.New()
	before := time.Now()
	c.RecordScan(5 * time.Millisecond)
	s := c.Snapshot()
	if s.LastScanAt.Before(before) {
		t.Fatal("LastScanAt should be after test start")
	}
}

func TestRecordAlert_IncrementsCounter(t *testing.T) {
	c := metrics.New()
	c.RecordAlert()
	c.RecordAlert()
	c.RecordAlert()
	s := c.Snapshot()
	if s.TotalAlerts != 3 {
		t.Fatalf("expected 3, got %d", s.TotalAlerts)
	}
}

func TestAvgScanTime_ComputedCorrectly(t *testing.T) {
	c := metrics.New()
	c.RecordScan(10 * time.Millisecond)
	c.RecordScan(30 * time.Millisecond)
	s := c.Snapshot()
	want := 20 * time.Millisecond
	if s.AvgScanTime != want {
		t.Fatalf("expected %v, got %v", want, s.AvgScanTime)
	}
}

func TestReset_ZeroesCounters(t *testing.T) {
	c := metrics.New()
	c.RecordScan(5 * time.Millisecond)
	c.RecordAlert()
	c.Reset()
	s := c.Snapshot()
	if s.TotalScans != 0 || s.TotalAlerts != 0 {
		t.Fatal("expected zeroed counters after Reset")
	}
}

func TestSnapshot_IsConcurrentlySafe(t *testing.T) {
	c := metrics.New()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			c.RecordScan(time.Millisecond)
		}
		close(done)
	}()
	for i := 0; i < 50; i++ {
		_ = c.Snapshot()
	}
	<-done
}
