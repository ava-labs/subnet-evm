// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"errors"
	"testing"

	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TODO: re-write these tests with mocks and allow for multiple precompile addresses
var _ precompileconfig.Predicater = (*mockPredicater)(nil)

type mockPredicater struct {
	predicateFunc    func(*precompileconfig.PredicateContext, [][]byte) []byte
	predicateGasFunc func([]byte) (uint64, error)
}

func (m *mockPredicater) VerifyPredicate(predicateContext *precompileconfig.PredicateContext, b [][]byte) []byte {
	return m.predicateFunc(predicateContext, b)
}

func (m *mockPredicater) PredicateGas(b []byte) (uint64, error) {
	if m.predicateGasFunc == nil {
		return 0, nil
	}
	return m.predicateGasFunc(b)
}

type predicateCheckTest struct {
	address     common.Address
	predicater  precompileconfig.Predicater
	accessList  types.AccessList
	gas         uint64
	expectedRes map[common.Address][]byte
	expectedErr error
}

func TestCheckPredicate(t *testing.T) {
	testErr := errors.New("test error")
	addr1 := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	addr2 := common.HexToAddress("0xb94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	predicateResultBytes := []byte{1, 2, 3}
	for name, test := range map[string]predicateCheckTest{
		// TODO: add test for multiple precompile predicates checked
		"no predicates, no access list passes": {
			gas:         53000,
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"no predicates, with access list passes": {
			gas: 57300,
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"predicate no access list passes": {
			address:     addr1,
			gas:         53000,
			predicater:  &mockPredicater{predicateFunc: func(*precompileconfig.PredicateContext, [][]byte) []byte { return nil }},
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"predicate valid access list passes": {
			address:    addr1,
			gas:        53000,
			predicater: &mockPredicater{predicateFunc: func(*precompileconfig.PredicateContext, [][]byte) []byte { return nil }},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedRes: map[common.Address][]byte{
				addr1: nil,
			},
			expectedErr: nil,
		},
		"predicate access list does not name precompile": {
			address:    addr1,
			gas:        57300,
			predicater: &mockPredicater{predicateFunc: func(*precompileconfig.PredicateContext, [][]byte) []byte { return nil }},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr2,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"predicate valid access list returns non-empty passes": {
			address:    addr1,
			gas:        53000,
			predicater: &mockPredicater{predicateFunc: func(*precompileconfig.PredicateContext, [][]byte) []byte { return predicateResultBytes }},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedErr: nil,
			expectedRes: map[common.Address][]byte{
				addr1: predicateResultBytes,
			},
		},
		"predicate invalid access list gas err": {
			address: addr1,
			gas:     53000,
			predicater: &mockPredicater{
				predicateGasFunc: func(b []byte) (uint64, error) { return 0, testErr },
				predicateFunc:    func(*precompileconfig.PredicateContext, [][]byte) []byte { return nil },
			},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedErr: testErr,
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			// Create the rules from TestChainConfig and update the predicates based on the test params
			rules := params.TestChainConfig.AvalancheRules(common.Big0, 0)
			if test.predicater != nil {
				rules.PredicatePrecompiles[test.address] = test.predicater
			}

			// Specify only the access list, since this test should not depend on any other values
			tx := types.NewTx(&types.DynamicFeeTx{
				AccessList: test.accessList,
				Gas:        test.gas,
			})
			predicateContext := &precompileconfig.PredicateContext{
				ProposerVMBlockCtx: &block.Context{
					PChainHeight: 10,
				},
			}
			predicateRes, err := CheckPredicates(rules, predicateContext, tx)
			if test.expectedErr == nil {
				require.NoError(err)
			} else {
				require.ErrorIs(err, test.expectedErr)
				return
			}
			require.Equal(test.expectedRes, predicateRes)
			intrinsicGas, err := IntrinsicGas(tx.Data(), tx.AccessList(), true, rules)
			require.NoError(err)
			require.Equal(tx.Gas(), intrinsicGas) // Require test specifies exact amount of gas consumed
		})
	}
}
