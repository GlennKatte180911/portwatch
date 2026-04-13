// Package notifier provides pluggable notification backends for portwatch.
// It defines a common interface and a console (stdout) implementation,
// with optional support for webhook-based delivery.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Notify(event Event) error
}

// Event represents a port-change notification payload.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Added     []int     `json:"added"`
	Removed   []int     `json:"removed"`
}

// ConsoleNotifier writes events to an io.Writer (defaults to os.Stdout).
type ConsoleNotifier struct {
	out io.Writer
}

// NewConsole returns a ConsoleNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewConsole(w io.Writer) *ConsoleNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &ConsoleNotifier{out: w}
}

// Notify prints the event in a human-readable format.
func (c *ConsoleNotifier) Notify(e Event) error {
	for _, p := range e.Added {
		fmt.Fprintf(c.out, "[%s] PORT OPENED: %d\n", e.Timestamp.Format(time.RFC3339), p)
	}
	for _, p := range e.Removed {
		fmt.Fprintf(c.out, "[%s] PORT CLOSED: %d\n", e.Timestamp.Format(time.RFC3339), p)
	}
	return nil
}

// WebhookNotifier posts events as JSON to a remote URL.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhook returns a WebhookNotifier that posts to url.
func NewWebhook(url string) *WebhookNotifier {
	return &WebhookNotifier{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Notify marshals the event to JSON and HTTP-POSTs it to the webhook URL.
func (w *WebhookNotifier) Notify(e Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("notifier: marshal event: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
