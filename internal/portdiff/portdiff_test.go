package portdiff_test

import (
	"testing"

	"github.com/example/portwatch/internal/portdiff"
)

func TestCompute_DetectsAdded(t *testing.T) {
	d := portdiff.Compute([]int{80, 443}, []int{80, 443, 8080})
	if len(d.Added) != 1 || d.Added[0] != 8080 {
		t.Fatalf("expected [8080] added, got %v", d.Added)
	}
	if len(d.Removed) != 0 {
		t.Fatalf("expected no removed, got %v", d.Removed)
	}
}

func TestCompute_DetectsRemoved(t *testing.T) {
	d := portdiff.Compute([]int{80, 443, 8080}, []int{80, 443})
	if len(d.Removed) != 1 || d.Removed[0] != 8080 {
		t.Fatalf("expected [8080] removed, got %v", d.Removed)
	}
	if len(d.Added) != 0 {
		t.Fatalf("expected no added, got %v", d.Added)
	}
}

func TestCompute_NoChange(t *testing.T) {
	d := portdiff.Compute([]int{80, 443}, []int{443, 80})
	if !d.IsEmpty() {
		t.Fatalf("expected empty diff, got added=%v removed=%v", d.Added, d.Removed)
	}
}

func TestCompute_BothEmpty(t *testing.T) {
	d := portdiff.Compute(nil, nil)
	if !d.IsEmpty() {
		t.Fatal("expected empty diff for nil inputs")
	}
}

func TestCompute_ResultIsSorted(t *testing.T) {
	d := portdiff.Compute([]int{}, []int{9000, 80, 443})
	for i := 1; i < len(d.Added); i++ {
		if d.Added[i] < d.Added[i-1] {
			t.Fatalf("added ports not sorted: %v", d.Added)
		}
	}
}

func TestSummary_NoChanges(t *testing.T) {
	s := portdiff.Summary(portdiff.Diff{})
	if s != "no changes" {
		t.Fatalf("unexpected summary: %q", s)
	}
}

func TestSummary_AddedAndRemoved(t *testing.T) {
	d := portdiff.Diff{Added: []int{8080}, Removed: []int{443, 80}}
	s := portdiff.Summary(d)
	if s == "" || s == "no changes" {
		t.Fatalf("expected non-empty summary, got %q", s)
	}
}
