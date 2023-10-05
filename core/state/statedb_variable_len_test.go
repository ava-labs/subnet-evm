// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVariableLengthKeys(t *testing.T) {
	require := require.New(t)

	// Create an empty state database
	db := rawdb.NewMemoryDatabase()
	state, err := New(common.Hash{}, NewDatabase(db), nil)
	require.NoError(err)

	var (
		address = common.Address{0x01}
		key     = "key1"
		val     = string([]byte{0, 1})
	)

	result := state.GetStateVariableLength(address, key)
	require.Equal(string(EmptyVal), result)

	state.SetStateVariableLength(address, key, val)
	result = state.GetStateVariableLength(address, key)
	require.Equal(val, result)

	// Test that the value is persisted
	root, err := state.Commit(false, false)
	require.NoError(err)
	err = state.Database().TrieDB().Commit(root, false)
	require.NoError(err)

	state, err = New(root, NewDatabase(db), nil)
	require.NoError(err)

	result = state.GetStateVariableLength(address, key)
	require.Equal(val, result)
}
