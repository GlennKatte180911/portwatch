package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
)

func labelFunc(port int) string {
	switch port {
	case 80, 8080:
		return "HTTP server"
	case 443:
		return "HTTPS server"
	case 5432:
		return "Postgres database"
	case 6379:
		return "Redis cache"
	default:
		return ""
	}
}

func TestGroup_PartitionsByLabelPrefix(t *testing.T) {
	g := portgroup.New(labelFunc)
	groups := g.Group([]int{80, 443, 8080, 5432})

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	// groups are sorted by name: http, https, postgres
	if groups[0].Name != "http" {
		t.Errorf("expected first group 'http', got %q", groups[0].Name)
	}
	if len(groups[0].Ports) != 2 {
		t.Errorf("expected 2 ports in http group, got %d", len(groups[0].Ports))
	}
}

func TestGroup_UnknownPortsGoToOther(t *testing.T) {
	g := portgroup.New(labelFunc)
	groups := g.Group([]int{9999, 10000})

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "other" {
		t.Errorf("expected group 'other', got %q", groups[0].Name)
	}
}

func TestGroup_EmptySlice_ReturnsNoGroups(t *testing.T) {
	g := portgroup.New(labelFunc)
	groups := g.Group([]int{})
	if len(groups) != 0 {
		t.Errorf("expected no groups, got %d", len(groups))
	}
}

func TestGroup_PortsWithinGroupAreSorted(t *testing.T) {
	g := portgroup.New(labelFunc)
	groups := g.Group([]int{8080, 80})

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Ports[0] != 80 || groups[0].Ports[1] != 8080 {
		t.Errorf("expected ports sorted [80 8080], got %v", groups[0].Ports)
	}
}

func TestSummary_FormatIsCorrect(t *testing.T) {
	g := portgroup.New(labelFunc)
	summary := g.Summary([]int{80, 5432, 6379})

	if summary == "" {
		t.Error("expected non-empty summary")
	}
	// should contain group names
	for _, want := range []string{"http", "postgres", "redis"} {
		if !containsStr(summary, want) {
			t.Errorf("summary %q missing %q", summary, want)
		}
	}
}

func TestSummary_EmptyPorts_ReturnsEmpty(t *testing.T) {
	g := portgroup.New(labelFunc)
	if got := g.Summary([]int{}); got != "" {
		t.Errorf("expected empty summary, got %q", got)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub ||
		len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
