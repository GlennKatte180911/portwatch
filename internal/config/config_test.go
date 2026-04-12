package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault_ReturnsExpectedValues(t *testing.T) {
	cfg := Default()

	if cfg.PortRange != "1-1024" {
		t.Errorf("expected port range '1-1024', got %q", cfg.PortRange)
	}
	if cfg.ScanInterval != 60*time.Second {
		t.Errorf("expected scan interval 60s, got %v", cfg.ScanInterval)
	}
	if !cfg.AlertOnNew {
		t.Error("expected AlertOnNew to be true by default")
	}
	if cfg.AlertOnClosed {
		t.Error("expected AlertOnClosed to be false by default")
	}
}

func TestLoad_FileNotFound_ReturnsDefaults(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	def := Default()
	if cfg.PortRange != def.PortRange {
		t.Errorf("expected default port range, got %q", cfg.PortRange)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	orig := &Config{
		PortRange:     "8000-9000",
		ScanInterval:  30 * time.Second,
		SnapshotPath:  "/tmp/snap.json",
		AlertOnNew:    false,
		AlertOnClosed: true,
		Timeout:       200 * time.Millisecond,
	}

	if err := orig.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.PortRange != orig.PortRange {
		t.Errorf("PortRange: want %q, got %q", orig.PortRange, loaded.PortRange)
	}
	if loaded.ScanInterval != orig.ScanInterval {
		t.Errorf("ScanInterval: want %v, got %v", orig.ScanInterval, loaded.ScanInterval)
	}
	if loaded.AlertOnNew != orig.AlertOnNew {
		t.Errorf("AlertOnNew: want %v, got %v", orig.AlertOnNew, loaded.AlertOnNew)
	}
	if loaded.AlertOnClosed != orig.AlertOnClosed {
		t.Errorf("AlertOnClosed: want %v, got %v", orig.AlertOnClosed, loaded.AlertOnClosed)
	}
	if loaded.Timeout != orig.Timeout {
		t.Errorf("Timeout: want %v, got %v", orig.Timeout, loaded.Timeout)
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(path, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
