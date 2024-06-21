// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
)

func TestTxAllowListRun(t *testing.T) {
	allowlist.RunPrecompileWithAllowListTests(t, Module, state.NewTestStateDB, nil)
}

func BenchmarkTxAllowList(b *testing.B) {
	allowlist.BenchPrecompileWithAllowList(b, Module, state.NewTestStateDB, nil)
}
