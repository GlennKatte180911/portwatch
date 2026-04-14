package metrics

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.HandlerFunc that serves the current metrics
// snapshot as JSON. It is intended to be mounted alongside the health
// check endpoint exposed by internal/healthcheck.
func (c *Collector) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := c.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "failed to encode metrics", http.StatusInternalServerError)
		}
	}
}
