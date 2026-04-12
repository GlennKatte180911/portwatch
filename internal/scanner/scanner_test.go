package scanner

import (
	"net"
	"strconv"
	"testing"
)

// startTestServer opens a TCP listener on a random port and returns the port number and a stop function.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	// TCPAddr.Port is already an int
	port = ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	s := New("127.0.0.1")

	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}
	if ports[0].Number != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Number)
	}
	if ports[0].State != "open" {
		t.Errorf("expected state 'open', got %s", ports[0].State)
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := New("127.0.0.1")

	_, err := s.Scan(9000, 8000)
	if err == nil {
		t.Error("expected error for invalid port range, got nil")
	}

	_, err = s.Scan(0, 100)
	if err == nil {
		t.Error("expected error for port 0, got nil")
	}
}

func TestScan_NoOpenPorts(t *testing.T) {
	s := New("127.0.0.1")
	// Port 1 is almost certainly closed in a test environment
	ports, err := s.Scan(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(ports))
	}
}
