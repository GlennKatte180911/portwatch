package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/snapshot"
)

func TestReport_TextFormat_Added(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	d := snapshot.Diff{Added: []int{8080, 9090}, Removed: []int{}}

	if err := r.Report(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "OPENED  port 8080") {
		t.Errorf("expected OPENED 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "OPENED  port 9090") {
		t.Errorf("expected OPENED 9090 in output, got: %s", out)
	}
}

func TestReport_TextFormat_Removed(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	d := snapshot.Diff{Added: []int{}, Removed: []int{3000}}

	if err := r.Report(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "CLOSED  port 3000") {
		t.Errorf("expected CLOSED 3000 in output, got: %s", out)
	}
}

func TestReport_NoDiff_WritesNothing(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	d := snapshot.Diff{Added: []int{}, Removed: []int{}}

	if err := r.Report(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestReport_JSONFormat_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	d := snapshot.Diff{Added: []int{443}, Removed: []int{80}}

	if err := r.Report(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"timestamp", "added", "removed", "443", "80"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got: %s", want, out)
		}
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Ensure New(nil, ...) does not panic.
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
