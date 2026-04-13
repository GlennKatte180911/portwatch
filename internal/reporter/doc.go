// Package reporter provides utilities for formatting and writing
// port-change reports produced by the portwatch monitor.
//
// Supported output formats:
//
//	FormatText  – human-readable timestamped lines (default)
//	FormatJSON  – newline-delimited JSON objects
//
// Usage:
//
//	r := reporter.New(os.Stdout, reporter.FormatText)
//	if err := r.Report(diff); err != nil {
//	    log.Fatal(err)
//	}
package reporter
