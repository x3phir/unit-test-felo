//go:build functional

package scenarios

import (
	"context"
	"testing"
	"time"
)

func TestSmoke_ServicesAreReachable(t *testing.T) {
	requireSuite(t, "smoke", "critical-flow", "full-regression")
	_, sut := loadFunctionalContext(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := sut.Health.Check(ctx); err != nil {
		t.Fatalf("sut.Health.Check() error = %v", err)
	}
}
