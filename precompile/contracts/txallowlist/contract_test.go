// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
)

func TestTxAllowListRun(t *testing.T) {
	// TODO: add module specific tests
	allowlist.RunTestsWithAllowListSetup(t, Module, state.NewTestStateDB, allowlist.AllowListTests(Module))
}
