// Package portping provides TCP reachability probing for individual ports.
//
// It complements the port scanner by performing targeted latency-aware probes
// against known-open ports, useful for confirming liveness after a scan diff
// or verifying that a newly detected port is genuinely accepting connections.
//
// Basic usage:
//
//	prober := portping.New("127.0.0.1", 2*time.Second)
//	result := prober.Probe(ctx, 8080)
//	if result.Reachable {
//		fmt.Printf("port 8080 up, latency %v\n", result.Latency)
//	}
package portping
