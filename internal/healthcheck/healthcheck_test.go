package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func TestNew_InitialisesChecker(t *testing.T) {
	c := healthcheck.New()
	s := c.Status()
	if !s.OK {
		t.Fatal("expected OK to be true")
	}
	if s.Scans != 0 {
		t.Fatalf("expected 0 scans, got %d", s.Scans)
	}
	if s.StartedAt.IsZero() {
		t.Fatal("expected StartedAt to be set")
	}
}

func TestRecordScan_IncrementsCounter(t *testing.T) {
	c := healthcheck.New()
	c.RecordScan()
	c.RecordScan()
	s := c.Status()
	if s.Scans != 2 {
		t.Fatalf("expected 2 scans, got %d", s.Scans)
	}
}

func TestRecordScan_UpdatesLastScan(t *testing.T) {
	c := healthcheck.New()
	before := time.Now()
	c.RecordScan()
	s := c.Status()
	if s.LastScan.Before(before) {
		t.Fatal("LastScan should be after the scan was recorded")
	}
}

func TestHandler_ReturnsJSON(t *testing.T) {
	c := healthcheck.New()
	c.RecordScan()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}

	var payload healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !payload.OK {
		t.Fatal("expected ok=true in JSON payload")
	}
	if payload.Scans != 1 {
		t.Fatalf("expected 1 scan in payload, got %d", payload.Scans)
	}
}

func TestHandler_UptimeNonEmpty(t *testing.T) {
	c := healthcheck.New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c.Handler()(rec, req)

	var payload healthcheck.Status
	_ = json.NewDecoder(rec.Body).Decode(&payload)
	if payload.Uptime == "" {
		t.Fatal("expected non-empty uptime string")
	}
}
