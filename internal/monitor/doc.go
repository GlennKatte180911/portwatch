// Package monitor provides the core monitoring loop for portwatch.
//
// It ties together the scanner, snapshot, and alert packages to periodically
// scan a configured port range, compare results against the previous snapshot,
// and emit notifications when ports are opened or closed.
//
// Usage:
//
//	cfg := config.Default()
//	notifier := alert.NewNotifier(nil)
//	mon := monitor.New(cfg, notifier)
//
//	done := make(chan struct{})
//	go func() {
//		<-signalChan
//		close(done)
//	}()
//
//	if err := mon.Run(done); err != nil {
//		log.Fatal(err)
//	}
package monitor
