package notifier_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notifier"
)

func sampleEvent() notifier.Event {
	return notifier.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Added:     []int{8080, 9090},
		Removed:   []int{3000},
	}
}

func TestConsoleNotifier_WritesAdded(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.NewConsole(&buf)
	if err := n.Notify(sampleEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT OPENED: 8080") {
		t.Errorf("expected PORT OPENED 8080, got: %s", out)
	}
	if !strings.Contains(out, "PORT OPENED: 9090") {
		t.Errorf("expected PORT OPENED 9090, got: %s", out)
	}
}

func TestConsoleNotifier_WritesRemoved(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.NewConsole(&buf)
	if err := n.Notify(sampleEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT CLOSED: 3000") {
		t.Errorf("expected PORT CLOSED 3000, got: %s", out)
	}
}

func TestConsoleNotifier_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic when w is nil.
	n := notifier.NewConsole(nil)
	if n == nil {
		t.Fatal("expected non-nil ConsoleNotifier")
	}
}

func TestWebhookNotifier_PostsJSON(t *testing.T) {
	var received notifier.Event
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notifier.NewWebhook(ts.URL)
	e := sampleEvent()
	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Added) != 2 {
		t.Errorf("expected 2 added ports, got %d", len(received.Added))
	}
	if len(received.Removed) != 1 {
		t.Errorf("expected 1 removed port, got %d", len(received.Removed))
	}
}

func TestWebhookNotifier_Non2xx_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notifier.NewWebhook(ts.URL)
	if err := n.Notify(sampleEvent()); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}
