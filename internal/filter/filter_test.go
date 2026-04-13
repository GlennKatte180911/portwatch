package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestNew_ValidRanges(t *testing.T) {
	f, err := filter.New([]string{"22", "8000-9000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNew_InvalidRange_ReturnsError(t *testing.T) {
	cases := []string{"abc", "0", "70000", "9000-8000", ""}
	for _, c := range cases {
		_, err := filter.New([]string{c})
		if err == nil {
			t.Errorf("expected error for range %q, got nil", c)
		}
	}
}

func TestAllowed_ExcludesMatchingPort(t *testing.T) {
	f, _ := filter.New([]string{"22", "8000-8080"})

	if f.Allowed(22) {
		t.Error("port 22 should be excluded")
	}
	if f.Allowed(8040) {
		t.Error("port 8040 should be excluded")
	}
	if !f.Allowed(80) {
		t.Error("port 80 should be allowed")
	}
	if !f.Allowed(8081) {
		t.Error("port 8081 should be allowed")
	}
}

func TestApply_FiltersSlice(t *testing.T) {
	f, _ := filter.New([]string{"22", "8000-8080"})

	input := []int{22, 80, 443, 8000, 8080, 8081, 9000}
	got := f.Apply(input)
	want := []int{80, 443, 8081, 9000}

	if len(got) != len(want) {
		t.Fatalf("Apply returned %v, want %v", got, want)
	}
	for i, p := range want {
		if got[i] != p {
			t.Errorf("Apply[%d] = %d, want %d", i, got[i], p)
		}
	}
}

func TestApply_EmptyFilter_AllowsAll(t *testing.T) {
	f, _ := filter.New(nil)
	input := []int{22, 80, 443}
	got := f.Apply(input)
	if len(got) != len(input) {
		t.Errorf("expected all ports allowed, got %v", got)
	}
}
