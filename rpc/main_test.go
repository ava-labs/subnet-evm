package rpc

import (
	"testing"

	"go.uber.org/goleak"
)

// TestMain uses goleak to verify tests in this package do not leak unexpected
// goroutines.
func TestMain(m *testing.M) {
	opts := []goleak.Option{
		// No good way to shut down these goroutines:
		goleak.IgnoreTopFunction("github.com/ava-labs/subnet-evm/metrics.(*meterArbiter).tick"),
	}
	goleak.VerifyTestMain(m, opts...)
}
