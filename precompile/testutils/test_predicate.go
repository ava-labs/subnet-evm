// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package testutils

import (
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/stretchr/testify/require"
)

type PredicateTest struct {
	Config precompileconfig.Config

	ProposerVMBlockContext *block.Context
	SnowContext            *snow.Context

	StorageSlots         []byte
	Gas                  uint64
	GasErr, PredicateErr error
}

func (test PredicateTest) Run(t testing.TB) {
	var (
		gas                  uint64
		gasErr, predicateErr error
	)
	switch predicate := test.Config.(type) {
	case precompileconfig.PrecompilePredicater:
		gas, gasErr = predicate.PredicateGas(test.StorageSlots)
		if gasErr == nil {
			predicateErr = predicate.VerifyPredicate(
				&precompileconfig.PrecompilePredicateContext{
					SnowCtx: test.SnowContext,
				},
				test.StorageSlots,
			)
		}
	case precompileconfig.ProposerPredicater:
		gas, gasErr = predicate.PredicateGas(test.StorageSlots)
		if gasErr == nil {
			predicateErr = predicate.VerifyPredicate(
				&precompileconfig.ProposerPredicateContext{
					PrecompilePredicateContext: precompileconfig.PrecompilePredicateContext{
						SnowCtx: test.SnowContext,
					},
					ProposerVMBlockCtx: test.ProposerVMBlockContext,
				},
				test.StorageSlots,
			)
		}
	default:
		t.Fatal("ran predicate test with precompileconfig that does not support a predicate")
	}

	if test.GasErr != nil {
		// If an error occurs here, the test finishes here
		require.ErrorIs(t, gasErr, test.GasErr)
		return
	}
	require.Equal(t, test.Gas, gas)
	if test.PredicateErr == nil {
		require.NoError(t, predicateErr)
	} else {
		require.ErrorIs(t, predicateErr, test.PredicateErr)
	}
}

func RunPredicateTests(t *testing.T, predicateTests map[string]PredicateTest) {
	t.Helper()

	for name, test := range predicateTests {
		t.Run(name, func(t *testing.T) {
			test.Run(t)
		})
	}
}

func RunPredicateBenchmarks(b *testing.B, predicateTests map[string]PredicateTest) {
	b.Helper()

	for name, test := range predicateTests {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			start := time.Now()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				test.Run(b)
			}
			b.StopTimer()
			elapsed := uint64(time.Since(start))
			if elapsed < 1 {
				elapsed = 1
			}

			b.ReportMetric(float64(test.Gas), "gas/op")
			// Keep it as uint64, multiply 100 to get two digit float later
			mgasps := (100 * 1000 * test.Gas) / elapsed
			b.ReportMetric(float64(mgasps)/100, "mgas/s")
		})
	}
}
