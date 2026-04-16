// Package portping provides lightweight TCP reachability probing for open ports.
package portping

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single probe.
type Result struct {
	Port    int
	Reachable bool
	Latency time.Duration
	Err     error
}

// Prober probes TCP ports for reachability.
type Prober struct {
	host    string
	timeout time.Duration
}

// New returns a Prober targeting host with the given dial timeout.
func New(host string, timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Prober{host: host, timeout: timeout}
}

// Probe attempts a TCP connection to the given port and returns a Result.
func (p *Prober) Probe(ctx context.Context, port int) Result {
	addr := fmt.Sprintf("%s:%d", p.host, port)
	start := time.Now()

	dialCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(dialCtx, "tcp", addr)
	latency := time.Since(start)

	if err != nil {
		return Result{Port: port, Reachable: false, Latency: latency, Err: err}
	}
	conn.Close()
	return Result{Port: port, Reachable: true, Latency: latency}
}

// ProbeAll probes each port in the slice and returns all results.
func (p *Prober) ProbeAll(ctx context.Context, ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		if ctx.Err() != nil {
			break
		}
		results = append(results, p.Probe(ctx, port))
	}
	return results
}
