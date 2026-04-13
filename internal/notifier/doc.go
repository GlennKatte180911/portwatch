// Package notifier provides notification backends for portwatch.
//
// A Notifier receives a [monitor.Event] and dispatches an alert through
// a specific channel. Two implementations are included out of the box:
//
//   - ConsoleNotifier – writes human-readable messages to an io.Writer
//     (defaults to os.Stdout). Useful for interactive terminal sessions
//     and piping output to other tools.
//
//   - WebhookNotifier – HTTP-POSTs a JSON payload to a configured URL.
//     The payload contains the event timestamp, lists of added and removed
//     ports, and the hostname of the machine being monitored. Suitable for
//     integrating with Slack incoming webhooks, PagerDuty, or any custom
//     HTTP endpoint.
//
// # Implementing a custom Notifier
//
// Any type that satisfies the Notifier interface can be plugged into the
// monitor pipeline:
//
//	type Notifier interface {
//		Notify(ctx context.Context, event Event) error
//	}
//
// # Error handling
//
// Notify returns an error when the underlying transport fails (e.g. a
// webhook endpoint is unreachable). The monitor will log the error and
// continue scanning; it will not stop on a notification failure.
package notifier
