package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port with its protocol and state.
type Port struct {
	Number   int
	Protocol string
	State    string
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner with sensible defaults.
func New(host string) *Scanner {
	return &Scanner{
		Host:    host,
		Timeout: 500 * time.Millisecond,
	}
}

// Scan checks a range of ports and returns those that are open.
func (s *Scanner) Scan(startPort, endPort int) ([]Port, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	var openPorts []Port

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		if err == nil {
			conn.Close()
			openPorts = append(openPorts, Port{
				Number:   port,
				Protocol: "tcp",
				State:    "open",
			})
		}
	}

	return openPorts, nil
}
