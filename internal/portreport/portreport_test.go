package portreport_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portreport"
)

func sampleEntry(port int, label, class, rank, policy string) portreport.Entry {
	return portreport.Entry{
		Port:      port,
		Label:     label,
		Class:     class,
		Rank:      rank,
		Policy:    policy,
		SeenCount: 3,
		FirstSeen: time.Now().Add(-time.Hour),
		LastSeen:  time.Now(),
	}
}

func TestBuild_SortsByPort(t *testing.T) {
	b := portreport.New()
	b.Add(sampleEntry(8080, "http-alt", "user", "medium", "allow"))
	b.Add(sampleEntry(22, "ssh", "system", "critical", "allow"))
	b.Add(sampleEntry(443, "https", "system", "critical", "allow"))

	r := b.Build()
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
	if r.Entries[0].Port != 22 || r.Entries[1].Port != 443 || r.Entries[2].Port != 8080 {
		t.Errorf("entries not sorted by port: %v", r.Entries)
	}
}

func TestBuild_GeneratedAtSet(t *testing.T) {
	b := portreport.New()
	r := b.Build()
	if r.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestBuild_EmptyEntries(t *testing.T) {
	b := portreport.New()
	r := b.Build()
	if len(r.Entries) != 0 {
		t.Errorf("expected no entries, got %d", len(r.Entries))
	}
}

func TestWriteText_ContainsPortAndLabel(t *testing.T) {
	b := portreport.New()
	b.Add(sampleEntry(80, "http", "system", "high", "allow"))
	r := b.Build()

	var buf bytes.Buffer
	if err := portreport.WriteText(&buf, r); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in text output")
	}
	if !strings.Contains(out, "http") {
		t.Error("expected label 'http' in text output")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	b := portreport.New()
	b.Add(sampleEntry(443, "https", "system", "critical", "allow"))
	r := b.Build()

	var buf bytes.Buffer
	if err := portreport.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out portreport.Report
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Entries) != 1 || out.Entries[0].Port != 443 {
		t.Errorf("unexpected JSON content: %+v", out)
	}
}
