// Package retry implements exponential-backoff retries for portwatch operations
// that may encounter transient failures.
//
// Basic usage:
//
//	r := retry.New(retry.DefaultConfig())
//	err := r.Do(ctx, func() error {
//		return sendWebhook(payload)
//	})
//	if errors.Is(err, retry.ErrMaxAttempts) {
//		log.Println("webhook delivery failed after all retries")
//	}
//
// The delay between attempts grows exponentially starting from BaseDelay and
// is capped at MaxDelay. Context cancellation is respected between attempts,
// allowing clean shutdown without waiting for the next backoff window.
package retry
