// Package portgroup groups ports by label prefix for structured reporting.
package portgroup

import (
	"fmt"
	"sort"
	"strings"
)

// Group holds a named collection of ports.
type Group struct {
	Name  string
	Ports []int
}

// Grouper partitions a slice of ports into named groups using a label lookup
// function. Ports whose label shares the same prefix word are placed together.
type Grouper struct {
	label func(port int) string
}

// New returns a Grouper that uses the provided label function.
func New(label func(port int) string) *Grouper {
	return &Grouper{label: label}
}

// Group partitions ports into named groups. Ports with the same first word of
// their label are combined. Ports with no label are placed in "other".
func (g *Grouper) Group(ports []int) []Group {
	buckets := make(map[string][]int)
	for _, p := range ports {
		lbl := g.label(p)
		key := bucketKey(lbl)
		buckets[key] = append(buckets[key], p)
	}

	groups := make([]Group, 0, len(buckets))
	for name, ps := range buckets {
		sort.Ints(ps)
		groups = append(groups, Group{Name: name, Ports: ps})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}

// Summary returns a human-readable one-line summary of the groups.
func (g *Grouper) Summary(ports []int) string {
	groups := g.Group(ports)
	parts := make([]string, 0, len(groups))
	for _, gr := range groups {
		parts = append(parts, fmt.Sprintf("%s(%d)", gr.Name, len(gr.Ports)))
	}
	return strings.Join(parts, ", ")
}

func bucketKey(label string) string {
	if label == "" {
		return "other"
	}
	parts := strings.Fields(label)
	return strings.ToLower(parts[0])
}
