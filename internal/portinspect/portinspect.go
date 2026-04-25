// Package portinspect provides detailed inspection of a single port,
// combining label, classification, rank, and policy information into
// a unified summary suitable for display or structured output.
package portinspect

import (
	"fmt"
	"time"
)

// Inspector aggregates per-port metadata from multiple providers.
type Inspector struct {
	labeler    Labeler
	classifier Classifier
	ranker      Ranker
}

// Labeler returns a human-readable name for a port.
type Labeler interface {
	Label(port int) string
}

// Classifier returns the classification category for a port.
type Classifier interface {
	Classify(port int) string
}

// Ranker returns the priority rank for a port.
type Ranker interface {
	Rank(port int) string
}

// Summary holds the inspection result for a single port.
type Summary struct {
	Port       int       `json:"port"`
	Label      string    `json:"label"`
	Class      string    `json:"class"`
	Rank       string    `json:"rank"`
	InspectedAt time.Time `json:"inspected_at"`
}

// String returns a compact human-readable representation of the summary.
func (s Summary) String() string {
	return fmt.Sprintf("port=%d label=%q class=%s rank=%s",
		s.Port, s.Label, s.Class, s.Rank)
}

// New creates an Inspector with the provided metadata providers.
// Any nil provider is replaced with a no-op fallback so callers may
// omit providers they do not need.
func New(l Labeler, c Classifier, r Ranker) *Inspector {
	if l == nil {
		l = noopLabeler{}
	}
	if c == nil {
		c = noopClassifier{}
	}
	if r == nil {
		r = noopRanker{}
	}
	return &Inspector{labeler: l, classifier: c, ranker: r}
}

// Inspect returns a Summary for the given port.
func (ins *Inspector) Inspect(port int) Summary {
	return Summary{
		Port:        port,
		Label:       ins.labeler.Label(port),
		Class:       ins.classifier.Classify(port),
		Rank:        ins.ranker.Rank(port),
		InspectedAt: time.Now().UTC(),
	}
}

// InspectAll returns a Summary for each port in the provided slice.
func (ins *Inspector) InspectAll(ports []int) []Summary {
	out := make([]Summary, len(ports))
	for i, p := range ports {
		out[i] = ins.Inspect(p)
	}
	return out
}

// --- no-op fallbacks ---

type noopLabeler struct{}

func (noopLabeler) Label(port int) string { return fmt.Sprintf("%d", port) }

type noopClassifier struct{}

func (noopClassifier) Classify(_ int) string { return "unknown" }

type noopRanker struct{}

func (noopRanker) Rank(_ int) string { return "low" }
