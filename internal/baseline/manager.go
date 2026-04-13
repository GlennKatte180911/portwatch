package baseline

import (
	"errors"
	"fmt"
	"io"
)

// Manager coordinates baseline lifecycle: loading, creating, and reporting.
type Manager struct {
	path   string
	output io.Writer
}

// NewManager returns a Manager that stores the baseline at path and writes
// status messages to output.
func NewManager(path string, output io.Writer) *Manager {
	return &Manager{path: path, output: output}
}

// Ensure loads an existing baseline or creates one from the provided ports
// if none exists. It returns the loaded or newly created baseline.
func (m *Manager) Ensure(currentPorts []int) (*Baseline, error) {
	b, err := Load(m.path)
	if err == nil {
		fmt.Fprintf(m.output, "baseline: loaded from %s (created %s)\n",
			m.path, b.CreatedAt.Format("2006-01-02 15:04:05"))
		return b, nil
	}
	if !errors.Is(err, ErrNoBaseline) {
		return nil, fmt.Errorf("baseline: load error: %w", err)
	}

	b, err = New(m.path, currentPorts)
	if err != nil {
		return nil, fmt.Errorf("baseline: create error: %w", err)
	}
	fmt.Fprintf(m.output, "baseline: created at %s with %d port(s)\n",
		m.path, len(currentPorts))
	return b, nil
}

// Reset overwrites the existing baseline with the provided ports.
func (m *Manager) Reset(currentPorts []int) (*Baseline, error) {
	b, err := New(m.path, currentPorts)
	if err != nil {
		return nil, fmt.Errorf("baseline: reset error: %w", err)
	}
	fmt.Fprintf(m.output, "baseline: reset at %s with %d port(s)\n",
		m.path, len(currentPorts))
	return b, nil
}
