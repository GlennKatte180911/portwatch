package portstate_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/portstate"
)

func TestNew_InitialisesEmpty(t *testing.T) {
	s := portstate.New()
	if s.Ports() != nil {
		t.Fatal("expected nil ports on new state")
	}
	if s.ScanCount() != 0 {
		t.Fatalf("expected scan count 0, got %d", s.ScanCount())
	}
}

func TestUpdate_ReplacesPorts(t *testing.T) {
	s := portstate.New()
	s.Update([]int{80, 443})

	ports := s.Ports()
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}
}

func TestUpdate_IncrementsScanCount(t *testing.T) {
	s := portstate.New()
	s.Update([]int{22})
	s.Update([]int{22, 80})

	if s.ScanCount() != 2 {
		t.Fatalf("expected scan count 2, got %d", s.ScanCount())
	}
}

func TestUpdate_SetsUpdatedAt(t *testing.T) {
	s := portstate.New()
	before := time.Now()
	s.Update([]int{8080})
	after := time.Now()

	at := s.UpdatedAt()
	if at.Before(before) || at.After(after) {
		t.Fatalf("UpdatedAt %v not within expected range [%v, %v]", at, before, after)
	}
}

func TestPorts_ReturnsCopy(t *testing.T) {
	s := portstate.New()
	s.Update([]int{9000})

	p1 := s.Ports()
	p1[0] = 1234
	p2 := s.Ports()

	if p2[0] == 1234 {
		t.Fatal("Ports should return an independent copy")
	}
}

func TestContains_ReturnsTrueForKnownPort(t *testing.T) {
	s := portstate.New()
	s.Update([]int{22, 80, 443})

	if !s.Contains(80) {
		t.Fatal("expected Contains(80) to return true")
	}
}

func TestContains_ReturnsFalseForUnknownPort(t *testing.T) {
	s := portstate.New()
	s.Update([]int{22})

	if s.Contains(9999) {
		t.Fatal("expected Contains(9999) to return false")
	}
}

func TestUpdate_EmptySlice_ClearsPorts(t *testing.T) {
	s := portstate.New()
	s.Update([]int{80, 443})
	s.Update([]int{})

	if len(s.Ports()) != 0 {
		t.Fatal("expected ports to be cleared after empty update")
	}
}
