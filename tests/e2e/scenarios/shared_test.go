//go:build e2e

package scenarios

import (
	"testing"

	"github.com/felo/felo-backend/tests/e2e/config"
	"github.com/felo/felo-backend/tests/e2e/harness"
)

func loadE2EContext(t *testing.T) (config.Config, harness.SystemUnderTest) {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("e2e tests skipped: %v", err)
	}
	sut, err := harness.NewSystemUnderTest(cfg)
	if err != nil {
		t.Skipf("e2e tests unavailable: %v", err)
	}
	return cfg, sut
}

func requireSuite(t *testing.T, allowed ...string) {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("e2e tests skipped: %v", err)
	}
	for _, suite := range allowed {
		if cfg.Suite == suite {
			return
		}
	}
	t.Skipf("suite %q is not enabled for this test", cfg.Suite)
}
