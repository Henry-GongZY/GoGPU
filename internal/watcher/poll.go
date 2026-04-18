package watcher

import (
	"context"
	"time"
)

// Poll triggers the provided fetch function periodically, pushing results to the returned channel.
// It stops and gracefully closes the channel when the provided context is canceled.
func Poll[T any](ctx context.Context, interval time.Duration, fetch func() T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		// Fire immediately on start
		select {
		case out <- fetch():
		case <-ctx.Done():
			return
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case out <- fetch():
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}
