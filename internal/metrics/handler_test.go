package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	c := metrics.New()
	c.RecordScan(15 * time.Millisecond)
	c.RecordAlert()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	c.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}
}

func TestHandler_FieldsPresent(t *testing.T) {
	c := metrics.New()
	c.RecordScan(10 * time.Millisecond)
	c.RecordScan(20 * time.Millisecond)
	c.RecordAlert()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	c.Handler()(rec, req)

	var snap metrics.Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if snap.TotalScans != 2 {
		t.Fatalf("expected 2 scans, got %d", snap.TotalScans)
	}
	if snap.TotalAlerts != 1 {
		t.Fatalf("expected 1 alert, got %d", snap.TotalAlerts)
	}
	if snap.AvgScanTime != 15*time.Millisecond {
		t.Fatalf("expected 15ms avg, got %v", snap.AvgScanTime)
	}
	if snap.UptimeSeconds <= 0 {
		t.Fatal("uptime should be positive")
	}
}

func TestHandler_EmptyCollector_StillResponds(t *testing.T) {
	c := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	c.Handler()(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for empty collector, got %d", rec.Code)
	}
}
