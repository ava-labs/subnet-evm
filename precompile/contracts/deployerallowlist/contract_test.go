// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/ethdb/memorydb"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contracts/test_utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestContractDeployerAllowListRun(t *testing.T) {
	// TODO: add module specific tests
	tests := make(map[string]test_utils.PrecompileTest)
	for name, test := range allowlist.AddAllowListTests(t, Module, tests) {
		t.Run(name, func(t *testing.T) {
			db := memorydb.New()
			stateDB, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			test.Run(t, Module, stateDB)
		})
	}
}
