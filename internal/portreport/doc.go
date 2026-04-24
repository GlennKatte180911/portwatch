// Package portreport assembles and renders structured summaries of the
// current port landscape observed by portwatch.
//
// A Builder collects Entry values — each describing a single port with its
// label, classification, rank, policy decision, and trend data — and produces
// a Report that can be serialised to either human-readable text or JSON.
//
// Typical usage:
//
//	b := portreport.New()
//	for _, port := range activePorts {
//		b.Add(portreport.Entry{
//			Port:   port,
//			Label:  labeler.Label(port),
//			Class:  classifier.Classify(port).String(),
//			Rank:   ranker.Get(port).String(),
//			Policy: policy.Evaluate(port).String(),
//		})
//	}
//	report := b.Build()
//	portreport.WriteText(os.Stdout, report)
//
// The Config type controls output format and optional entry limits.
package portreport
