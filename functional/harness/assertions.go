//go:build functional

package harness

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func Eventually(ctx context.Context, interval time.Duration, fn func() error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		err := fn()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return errors.Join(ctx.Err(), err)
		case <-ticker.C:
		}
	}
}

func ExpectEqual(got string, want string) error {
	if got != want {
		return fmt.Errorf("got %q want %q", got, want)
	}
	return nil
}
