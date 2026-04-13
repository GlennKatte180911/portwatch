package dedup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/dedup"
	"github.com/user/portwatch/internal/snapshot"
)

func sampleDiff(added, removed []int) snapshot.Diff {
	return snapshot.Diff{Added: added, Removed: removed}
}

func TestIsDuplicate_FirstCallReturnsFalse(t *testing.T) {
	d := dedup.New(5 * time.Second)
	if d.IsDuplicate(sampleDiff([]int{8080}, nil)) {
		t.Fatal("expected first occurrence to not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinWindowReturnsTrue(t *testing.T) {
	d := dedup.New(5 * time.Second)
	diff := sampleDiff([]int{8080}, nil)
	d.IsDuplicate(diff)
	if !d.IsDuplicate(diff) {
		t.Fatal("expected second identical diff within window to be a duplicate")
	}
}

func TestIsDuplicate_DifferentDiffsAreIndependent(t *testing.T) {
	d := dedup.New(5 * time.Second)
	d.IsDuplicate(sampleDiff([]int{8080}, nil))
	if d.IsDuplicate(sampleDiff([]int{9090}, nil)) {
		t.Fatal("expected different diff to not be a duplicate")
	}
}

func TestIsDuplicate_EmptyDiffNeverDuplicate(t *testing.T) {
	d := dedup.New(5 * time.Second)
	empty := sampleDiff(nil, nil)
	if d.IsDuplicate(empty) {
		t.Fatal("empty diff should never be considered a duplicate")
	}
	if d.IsDuplicate(empty) {
		t.Fatal("empty diff should never be considered a duplicate on second call")
	}
}

func TestIsDuplicate_AllowsAfterWindowExpires(t *testing.T) {
	d := dedup.New(50 * time.Millisecond)
	diff := sampleDiff([]int{8080}, nil)
	d.IsDuplicate(diff)
	time.Sleep(80 * time.Millisecond)
	if d.IsDuplicate(diff) {
		t.Fatal("expected diff to be allowed after window expiry")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	d := dedup.New(5 * time.Second)
	diff := sampleDiff([]int{443}, nil)
	d.IsDuplicate(diff)
	d.Reset()
	if d.IsDuplicate(diff) {
		t.Fatal("expected diff to pass through after Reset")
	}
}
