// Package semaphore provides a simple counting semaphore used by the port
// scanner to cap the number of goroutines that may perform TCP dial attempts
// simultaneously.
//
// # Usage
//
//	s, err := semaphore.New(256)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Inside each scan goroutine:
//	if err := s.Acquire(ctx); err != nil {
//		return err // context cancelled
//	}
//	defer s.Release()
//
// Calling Release without a prior Acquire panics to surface programming errors
// early.
package semaphore
