package portreport_test

import (
	"testing"

	"github.com/user/portwatch/internal/portreport"
)

func TestDefaultConfig_Values(t *testing.T) {
	c := portreport.DefaultConfig()
	if c.Format != portreport.FormatText {
		t.Errorf("expected format text, got %q", c.Format)
	}
	if !c.IncludeNotes {
		t.Error("expected IncludeNotes to be true")
	}
	if c.MaxEntries != 0 {
		t.Errorf("expected MaxEntries 0, got %d", c.MaxEntries)
	}
}

func TestValidate_InvalidFormat_ReturnsError(t *testing.T) {
	c := portreport.Config{Format: "xml"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestValidate_NegativeMaxEntries_ReturnsError(t *testing.T) {
	c := portreport.Config{Format: portreport.FormatJSON, MaxEntries: -1}
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative max_entries")
	}
}

func TestValidate_ValidConfig_NoError(t *testing.T) {
	c := portreport.Config{Format: portreport.FormatJSON, MaxEntries: 100}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestApply_FillsEmptyFormat(t *testing.T) {
	c := portreport.Config{}
	applied := c.Apply()
	if applied.Format != portreport.FormatText {
		t.Errorf("expected format text after Apply, got %q", applied.Format)
	}
}

func TestApply_PreservesExplicitFormat(t *testing.T) {
	c := portreport.Config{Format: portreport.FormatJSON}
	applied := c.Apply()
	if applied.Format != portreport.FormatJSON {
		t.Errorf("expected format json to be preserved, got %q", applied.Format)
	}
}
