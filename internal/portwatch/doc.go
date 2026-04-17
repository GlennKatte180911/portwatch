// Package portwatch provides the top-level scan-diff-notify pipeline for
// portwatch. It wires together a port scanner, a running port state tracker,
// and a notifier so that callers only need to construct a Config and call Run.
//
// Typical usage:
//
//	w, err := portwatch.New(portwatch.Config{
//		StartPort: 1,
//		EndPort:   1024,
//		Interval:  30 * time.Second,
//		Notifier:  myNotifier,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package portwatch
