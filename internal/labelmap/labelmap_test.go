package labelmap_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/labelmap"
)

func TestNew_ContainsWellKnownPorts(t *testing.T) {
	lm := labelmap.New()

	cases := map[int]string{
		22:  "ssh",
		80:  "http",
		443: "https",
	}
	for port, want := range cases {
		got, ok := lm.Get(port)
		if !ok {
			t.Errorf("port %d: expected entry, got none", port)
		}
		if got != want {
			t.Errorf("port %d: got %q, want %q", port, got, want)
		}
	}
}

func TestSet_OverwritesExistingLabel(t *testing.T) {
	lm := labelmap.New()
	lm.Set(80, "custom-http")

	got, ok := lm.Get(80)
	if !ok || got != "custom-http" {
		t.Errorf("expected %q, got %q (ok=%v)", "custom-http", got, ok)
	}
}

func TestLabel_FallbackForUnknownPort(t *testing.T) {
	lm := labelmap.New()
	got := lm.Label(9999)
	if got != "port/9999" {
		t.Errorf("expected fallback label, got %q", got)
	}
}

func TestLabel_ReturnsNameForKnownPort(t *testing.T) {
	lm := labelmap.New()
	got := lm.Label(22)
	if got != "ssh" {
		t.Errorf("expected %q, got %q", "ssh", got)
	}
}

func TestLoadFile_MergesEntries(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "labels.json")

	data := map[int]string{9000: "myapp", 9001: "myapp-metrics"}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}

	lm := labelmap.New()
	if err := lm.LoadFile(path); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	for port, want := range data {
		got, ok := lm.Get(port)
		if !ok || got != want {
			t.Errorf("port %d: got %q (ok=%v), want %q", port, got, ok, want)
		}
	}
	// Well-known entries must still be present.
	if _, ok := lm.Get(443); !ok {
		t.Error("expected well-known port 443 to still be present after LoadFile")
	}
}

func TestLoadFile_FileNotFound_ReturnsError(t *testing.T) {
	lm := labelmap.New()
	if err := lm.LoadFile("/nonexistent/labels.json"); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	lm := labelmap.New()
	snap := lm.Snapshot()
	snap[22] = "tampered"

	got, _ := lm.Get(22)
	if got == "tampered" {
		t.Error("Snapshot should return an independent copy")
	}
}
