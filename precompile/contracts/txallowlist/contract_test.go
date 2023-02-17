// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/ethdb/memorydb"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func newStateDB(t *testing.T) contract.StateDB {
	db := memorydb.New()
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
	require.NoError(t, err)
	return stateDB
}

func TestTxAllowListRun(t *testing.T) {
	// TODO: add module specific tests
	allowlist.RunTestsWithAllowListSetup(t, Module, newStateDB, allowlist.AllowListTests(Module))
}
