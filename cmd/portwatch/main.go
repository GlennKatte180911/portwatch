package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/portwatch/internal/scanner"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Host to scan")
	startPort := flag.Int("start", 1, "Start of port range")
	endPort := flag.Int("end", 1024, "End of port range")
	flag.Parse()

	fmt.Printf("Scanning %s ports %d-%d...\n", *host, *startPort, *endPort)

	s := scanner.New(*host)
	ports, err := s.Scan(*startPort, *endPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(ports) == 0 {
		fmt.Println("No open ports found.")
		return
	}

	fmt.Printf("Found %d open port(s):\n", len(ports))
	for _, p := range ports {
		fmt.Printf("  %-6d %s\t[%s]\n", p.Number, p.Protocol, p.State)
	}
}
