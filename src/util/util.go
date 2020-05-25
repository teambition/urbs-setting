package util

import (
	"context"
	"time"
)

// Go ...
func Go(du time.Duration, fn func(context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), du)

	go func() {
		fn(ctx)
		cancel()
	}()
}
