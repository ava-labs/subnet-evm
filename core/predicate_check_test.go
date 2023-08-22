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
	"go.uber.org/mock/gomock"
)

type predicateCheckTest struct {
	accessList       types.AccessList
	gas              uint64
	createPredicates func(t testing.TB) map[common.Address]precompileconfig.Predicater
	expectedRes      map[common.Address][]byte
	expectedErr      error
}

func TestCheckPredicate(t *testing.T) {
	testErr := errors.New("test error")
	addr1 := common.HexToAddress("0xaa")
	addr2 := common.HexToAddress("0xbb")
	addr3 := common.HexToAddress("0xcc")
	addr4 := common.HexToAddress("0xdd")
	predicateResultBytes := []byte{1, 2, 3}
	for name, test := range map[string]predicateCheckTest{
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
			gas: 53000,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
				}
			},
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"predicate named by access list returns empty": {
			gas: 53000,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicate.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(2)
				predicate.EXPECT().VerifyPredicate(gomock.Any(), [][]byte{arg[:]}).Return(nil)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
				}
			},
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
		"predicate named by access list returns non-empty": {
			gas: 53000,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicate.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(2)
				predicate.EXPECT().VerifyPredicate(gomock.Any(), [][]byte{arg[:]}).Return(predicateResultBytes)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
				}
			},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedRes: map[common.Address][]byte{
				addr1: predicateResultBytes,
			},
			expectedErr: nil,
		},
		"predicate returns gas err": {
			gas: 53000,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicate.EXPECT().PredicateGas(arg[:]).Return(uint64(0), testErr)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
				}
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
		"two predicates one named by access list returns non-empty": {
			gas: 53000,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicate.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(2)
				predicate.EXPECT().VerifyPredicate(gomock.Any(), [][]byte{arg[:]}).Return(predicateResultBytes)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
					addr2: predicate,
				}
			},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedRes: map[common.Address][]byte{
				addr1: predicateResultBytes,
			},
			expectedErr: nil,
		},
		"two predicates both named by access list returns non-empty": {
			gas: 53000,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicate.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(4)
				predicate.EXPECT().VerifyPredicate(gomock.Any(), [][]byte{arg[:]}).Return(predicateResultBytes).Times(2)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
					addr2: predicate,
				}
			},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
				{
					Address: addr2,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedRes: map[common.Address][]byte{
				addr1: predicateResultBytes,
				addr2: predicateResultBytes,
			},
			expectedErr: nil,
		},
		"two predicates niether named by access list": {
			gas: 61600,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
					addr2: predicate,
				}
			},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr3,
					StorageKeys: []common.Hash{
						{1},
					},
				},
				{
					Address: addr4,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"insufficient gas": {
			gas: 53000,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicate := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicate.EXPECT().PredicateGas(arg[:]).Return(uint64(1), nil)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate,
				}
			},
			accessList: types.AccessList([]types.AccessTuple{
				{
					Address: addr1,
					StorageKeys: []common.Hash{
						{1},
					},
				},
			}),
			expectedErr: ErrIntrinsicGas,
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			// Create the rules from TestChainConfig and update the predicates based on the test params
			rules := params.TestChainConfig.AvalancheRules(common.Big0, 0)
			if test.createPredicates != nil {
				for address, predicater := range test.createPredicates(t) {
					rules.PredicatePrecompiles[address] = predicater
				}
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
