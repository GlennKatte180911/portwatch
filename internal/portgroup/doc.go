// Package portgroup groups ports into named buckets based on their human-readable
// labels, making structured reporting and summarisation easier.
//
// A Grouper wraps a label function (for example from internal/portlabel) and
// partitions any slice of port numbers into []Group values. Ports whose labels
// share the same leading word are merged into a single bucket; unlabelled ports
// fall into the reserved "other" bucket.
//
// Example
//
//	labeller := portlabel.New(nil)
//	g := portgroup.New(labeller.Label)
//
//	groups := g.Group([]int{80, 443, 8080, 5432})
//	for _, gr := range groups {
//		fmt.Printf("%s: %v\n", gr.Name, gr.Ports)
//	}
//
//	// http: [80 8080]
//	// https: [443]
//	// postgres: [5432]
package portgroup
