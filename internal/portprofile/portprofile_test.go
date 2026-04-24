package portprofile_test

import (
	"testing"

	"github.com/user/portwatch/internal/portprofile"
)

func TestBuild_NoProviders_ReturnsEmptyFields(t *testing.T) {
	p := portprofile.New()
	pr := p.Build(80)
	if pr.Port != 80 {
		t.Fatalf("expected port 80, got %d", pr.Port)
	}
	if pr.Label != "" || pr.Rank != "" || pr.Class != "" || pr.Scope != "" {
		t.Errorf("expected empty fields without providers, got %+v", pr)
	}
}

func TestBuild_WithAllProviders(t *testing.T) {
	p := portprofile.New().
		WithLabeler(func(port int) string { return "http" }).
		WithRanker(func(port int) string { return "high" }).
		WithClasser(func(port int) string { return "web" }).
		WithScoper(func(port int) string { return "public" })

	pr := p.Build(80)
	if pr.Label != "http" {
		t.Errorf("expected label 'http', got %q", pr.Label)
	}
	if pr.Rank != "high" {
		t.Errorf("expected rank 'high', got %q", pr.Rank)
	}
	if pr.Class != "web" {
		t.Errorf("expected class 'web', got %q", pr.Class)
	}
	if pr.Scope != "public" {
		t.Errorf("expected scope 'public', got %q", pr.Scope)
	}
}

func TestBuild_CriticalRank_AddsNote(t *testing.T) {
	p := portprofile.New().
		WithRanker(func(int) string { return "critical" })

	pr := p.Build(22)
	if len(pr.Notes) == 0 {
		t.Fatal("expected at least one note for critical rank")
	}
	if pr.Notes[0] == "" {
		t.Error("note should not be empty")
	}
}

func TestBuild_SystemClassHighPort_AddsNote(t *testing.T) {
	p := portprofile.New().
		WithClasser(func(int) string { return "system" })

	pr := p.Build(8080)
	found := false
	for _, n := range pr.Notes {
		if n != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected a note for system-class port above 1023")
	}
}

func TestBuildAll_ReturnsProfilePerPort(t *testing.T) {
	p := portprofile.New().
		WithLabeler(func(port int) string {
			if port == 443 {
				return "https"
			}
			return "other"
		})

	profiles := p.BuildAll([]int{80, 443, 8080})
	if len(profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(profiles))
	}
	if profiles[1].Label != "https" {
		t.Errorf("expected 'https' for port 443, got %q", profiles[1].Label)
	}
}

func TestBuildAll_EmptySlice_ReturnsEmpty(t *testing.T) {
	p := portprofile.New()
	profiles := p.BuildAll(nil)
	if len(profiles) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(profiles))
	}
}
