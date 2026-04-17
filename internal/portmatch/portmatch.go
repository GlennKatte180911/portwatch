// Package portmatch provides pattern-based port matching using glob-style
// expressions and named aliases (e.g. "web", "db").
package portmatch

import (
	"fmt"
	"strconv"
	"strings"
)

// Matcher checks whether a port satisfies one or more patterns.
type Matcher struct {
	patterns []pattern
}

type pattern struct {
	raw string
	low int
	high int
}

var aliases = map[string][]string{
	"web":      {"80", "443", "8080", "8443"},
	"db":       {"3306", "5432", "6379", "27017"},
	"mail":     {"25", "465", "587", "993", "995"},
	"ssh":      {"22"},
	"dns":      {"53"},
	"system":   {"1-1023"},
}

// New creates a Matcher from a slice of pattern strings.
// Patterns may be:
//   - a single port number: "80"
//   - a range: "8000-9000"
//   - a named alias: "web", "db"
func New(patterns []string) (*Matcher, error) {
	m := &Matcher{}
	for _, p := range patterns {
		expanded := expand(p)
		for _, e := range expanded {
			pat, err := parsePattern(e)
			if err != nil {
				return nil, fmt.Errorf("portmatch: invalid pattern %q: %w", p, err)
			}
			m.patterns = append(m.patterns, pat)
		}
	}
	return m, nil
}

// Match returns true if port satisfies any pattern.
func (m *Matcher) Match(port int) bool {
	for _, p := range m.patterns {
		if port >= p.low && port <= p.high {
			return true
		}
	}
	return false
}

// Filter returns only the ports that match.
func (m *Matcher) Filter(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if m.Match(p) {
			out = append(out, p)
		}
	}
	return out
}

func expand(p string) []string {
	if v, ok := aliases[strings.ToLower(p)]; ok {
		return v
	}
	return []string{p}
}

func parsePattern(s string) (pattern, error) {
	if idx := strings.Index(s, "-"); idx != -1 {
		lo, err1 := strconv.Atoi(s[:idx])
		hi, err2 := strconv.Atoi(s[idx+1:])
		if err1 != nil || err2 != nil {
			return pattern{}, fmt.Errorf("bad range")
		}
		if lo > hi {
			return pattern{}, fmt.Errorf("low > high")
		}
		return pattern{raw: s, low: lo, high: hi}, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return pattern{}, fmt.Errorf("not a number")
	}
	return pattern{raw: s, low: n, high: n}, nil
}
