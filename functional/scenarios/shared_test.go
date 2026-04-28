//go:build functional

package scenarios

import (
	"testing"

	"github.com/felo/felo-backend/functional/config"
	"github.com/felo/felo-backend/functional/harness"
)

func loadFunctionalContext(t *testing.T) (config.Config, harness.SystemUnderTest) {
	t.Helper()

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("functional tests skipped: %v", err)
	}

	sut, err := harness.NewSystemUnderTest(cfg)
	if err != nil {
		t.Skipf("functional tests scaffold ready, adapter wiring pending: %v", err)
	}

	return cfg, sut
}

func requireSuite(t *testing.T, allowed ...string) {
	t.Helper()

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("functional tests skipped: %v", err)
	}
	for _, suite := range allowed {
		if cfg.Suite == suite {
			return
		}
	}

	t.Skipf("suite %q is not enabled for this test", cfg.Suite)
}
