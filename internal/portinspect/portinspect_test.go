package portinspect_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portinspect"
)

// --- stub implementations ---

type stubLabeler struct{ label string }

func (s stubLabeler) Label(_ int) string { return s.label }

type stubClassifier struct{ class string }

func (s stubClassifier) Classify(_ int) string { return s.class }

type stubRanker struct{ rank string }

func (s stubRanker) Rank(_ int) string { return s.rank }

// --- tests ---

func TestInspect_PopulatesAllFields(t *testing.T) {
	ins := portinspect.New(
		stubLabeler{"http"},
		stubClassifier{"web"},
		stubRanker{"high"},
	)
	s := ins.Inspect(80)

	if s.Port != 80 {
		t.Errorf("expected port 80, got %d", s.Port)
	}
	if s.Label != "http" {
		t.Errorf("expected label 'http', got %q", s.Label)
	}
	if s.Class != "web" {
		t.Errorf("expected class 'web', got %q", s.Class)
	}
	if s.Rank != "high" {
		t.Errorf("expected rank 'high', got %q", s.Rank)
	}
	if s.InspectedAt.IsZero() {
		t.Error("expected InspectedAt to be set")
	}
}

func TestInspect_NilProviders_UseFallbacks(t *testing.T) {
	ins := portinspect.New(nil, nil, nil)
	s := ins.Inspect(9999)

	if s.Port != 9999 {
		t.Errorf("expected port 9999, got %d", s.Port)
	}
	if s.Label == "" {
		t.Error("expected non-empty fallback label")
	}
	if s.Class != "unknown" {
		t.Errorf("expected fallback class 'unknown', got %q", s.Class)
	}
	if s.Rank != "low" {
		t.Errorf("expected fallback rank 'low', got %q", s.Rank)
	}
}

func TestInspectAll_ReturnsOneSummaryPerPort(t *testing.T) {
	ins := portinspect.New(nil, nil, nil)
	ports := []int{22, 80, 443}

	summaries := ins.InspectAll(ports)

	if len(summaries) != len(ports) {
		t.Fatalf("expected %d summaries, got %d", len(ports), len(summaries))
	}
	for i, s := range summaries {
		if s.Port != ports[i] {
			t.Errorf("index %d: expected port %d, got %d", i, ports[i], s.Port)
		}
	}
}

func TestInspectAll_EmptySlice_ReturnsEmpty(t *testing.T) {
	ins := portinspect.New(nil, nil, nil)
	result := ins.InspectAll([]int{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestSummary_String_ContainsPortAndLabel(t *testing.T) {
	ins := portinspect.New(stubLabeler{"ssh"}, nil, nil)
	s := ins.Inspect(22)
	str := s.String()

	if !strings.Contains(str, fmt.Sprintf("%d", 22)) {
		t.Errorf("String() missing port number: %q", str)
	}
	if !strings.Contains(str, "ssh") {
		t.Errorf("String() missing label: %q", str)
	}
}
