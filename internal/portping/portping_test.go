package portping_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
)

func startTCPServer(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func TestProbe_ReachablePort(t *testing.T) {
	port := startTCPServer(t)
	p := portping.New("127.0.0.1", time.Second)
	res := p.Probe(context.Background(), port)
	if !res.Reachable {
		t.Fatalf("expected reachable, got err: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbe_UnreachablePort(t *testing.T) {
	p := portping.New("127.0.0.1", 100*time.Millisecond)
	res := p.Probe(context.Background(), 1)
	if res.Reachable {
		t.Fatal("expected unreachable")
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestProbeAll_ReturnsResultPerPort(t *testing.T) {
	port := startTCPServer(t)
	p := portping.New("127.0.0.1", time.Second)
	results := p.ProbeAll(context.Background(), []int{port, 1})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Reachable {
		t.Error("first port should be reachable")
	}
	if results[1].Reachable {
		t.Error("second port should be unreachable")
	}
}

func TestProbeAll_CancelledContext_StopsEarly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := portping.New("127.0.0.1", time.Second)
	results := p.ProbeAll(ctx, []int{80, 443, 8080})
	if len(results) != 0 {
		t.Errorf("expected 0 results with cancelled context, got %d", len(results))
	}
}

func TestNew_DefaultTimeout(t *testing.T) {
	p := portping.New("localhost", 0)
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}
