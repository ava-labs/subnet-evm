// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type predicateCheckTest struct {
	accessList       types.AccessList
	gas              uint64
	predicateContext *precompileconfig.PredicateContext
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
	predicateContext := &precompileconfig.PredicateContext{
		ProposerVMBlockCtx: &block.Context{
			PChainHeight: 10,
		},
	}
	for name, test := range map[string]predicateCheckTest{
		"no predicates, no access list, no context passes": {
			gas:              53000,
			predicateContext: nil,
			expectedRes:      make(map[common.Address][]byte),
			expectedErr:      nil,
		},
		"no predicates, no access list, with context passes": {
			gas:              53000,
			predicateContext: predicateContext,
			expectedRes:      make(map[common.Address][]byte),
			expectedErr:      nil,
		},
		"no predicates, with access list, no context passes": {
			gas:              57300,
			predicateContext: nil,
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
		"predicate, no access list, no context passes": {
			gas:              53000,
			predicateContext: nil,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
				}
			},
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"predicate, no access list, no block context passes": {
			gas: 53000,
			predicateContext: &precompileconfig.PredicateContext{
				ProposerVMBlockCtx: nil,
			},
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
				}
			},
			expectedRes: make(map[common.Address][]byte),
			expectedErr: nil,
		},
		"predicate named by access list, without context errors": {
			gas:              53000,
			predicateContext: nil,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicater.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(1)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
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
			expectedErr: ErrMissingPredicateContext,
		},
		"predicate named by access list, without block context errors": {
			gas: 53000,
			predicateContext: &precompileconfig.PredicateContext{
				ProposerVMBlockCtx: nil,
			},
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicater.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(1)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
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
			expectedErr: ErrMissingPredicateContext,
		},
		"predicate named by access list returns non-empty": {
			gas:              53000,
			predicateContext: predicateContext,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicater.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(2)
				predicater.EXPECT().VerifyPredicate(gomock.Any(), arg[:]).Return(true)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
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
				addr1: {}, // valid bytes
			},
			expectedErr: nil,
		},
		"predicate returns gas err": {
			gas:              53000,
			predicateContext: predicateContext,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicater.EXPECT().PredicateGas(arg[:]).Return(uint64(0), testErr)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
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
			gas:              53000,
			predicateContext: predicateContext,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicater.EXPECT().PredicateGas(arg[:]).Return(uint64(0), nil).Times(2)
				predicater.EXPECT().VerifyPredicate(gomock.Any(), arg[:]).Return(true)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
					addr2: predicater,
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
				addr1: {}, // valid bytes
			},
			expectedErr: nil,
		},
		"two predicates both named by access list returns non-empty": {
			gas:              53000,
			predicateContext: predicateContext,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				ctrl := gomock.NewController(t)
				predicate1 := precompileconfig.NewMockPredicater(ctrl)
				arg1 := common.Hash{1}
				predicate1.EXPECT().PredicateGas(arg1[:]).Return(uint64(0), nil).Times(2)
				predicate1.EXPECT().VerifyPredicate(gomock.Any(), arg1[:]).Return(true)
				predicate2 := precompileconfig.NewMockPredicater(ctrl)
				arg2 := common.Hash{2}
				predicate2.EXPECT().PredicateGas(arg2[:]).Return(uint64(0), nil).Times(2)
				predicate2.EXPECT().VerifyPredicate(gomock.Any(), arg2[:]).Return(false)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicate1,
					addr2: predicate2,
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
						{2},
					},
				},
			}),
			expectedRes: map[common.Address][]byte{
				addr1: {},  // valid bytes
				addr2: {1}, // invalid bytes
			},
			expectedErr: nil,
		},
		"two predicates niether named by access list": {
			gas:              61600,
			predicateContext: predicateContext,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
					addr2: predicater,
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
			gas:              53000,
			predicateContext: predicateContext,
			createPredicates: func(t testing.TB) map[common.Address]precompileconfig.Predicater {
				predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
				arg := common.Hash{1}
				predicater.EXPECT().PredicateGas(arg[:]).Return(uint64(1), nil)
				return map[common.Address]precompileconfig.Predicater{
					addr1: predicater,
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
					rules.Predicaters[address] = predicater
				}
			}

			// Specify only the access list, since this test should not depend on any other values
			tx := types.NewTx(&types.DynamicFeeTx{
				AccessList: test.accessList,
				Gas:        test.gas,
			})
			predicateRes, err := CheckPredicates(rules, test.predicateContext, tx)
			require.ErrorIs(err, test.expectedErr)
			if test.expectedErr != nil {
				return
			}
			require.Equal(test.expectedRes, predicateRes)
			intrinsicGas, err := IntrinsicGas(tx.Data(), tx.AccessList(), true, rules)
			require.NoError(err)
			require.Equal(tx.Gas(), intrinsicGas) // Require test specifies exact amount of gas consumed
		})
	}
}

var (
	validHash   = common.Hash{1}
	invalidHash = common.Hash{2}
)

func TestCheckPredicatesOutput(t *testing.T) {
	addr1 := common.HexToAddress("0xaa")
	addr2 := common.HexToAddress("0xbb")

	predicateContext := &precompileconfig.PredicateContext{
		ProposerVMBlockCtx: &block.Context{
			PChainHeight: 10,
		},
	}
	type testTuple struct {
		address          common.Address
		isValidPredicate bool
	}
	type resultTest struct {
		name        string
		expectedRes map[common.Address][]byte
		testTuple   []testTuple
	}
	tests := []resultTest{
		{name: "no predicates", expectedRes: map[common.Address][]byte{}},
		{
			name: "one address one predicate",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: true},
			},
			expectedRes: map[common.Address][]byte{addr1: {}},
		},
		{
			name: "one address one invalid predicate",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: false},
			},
			expectedRes: map[common.Address][]byte{addr1: {1}},
		},
		{
			name: "one address two invalid predicates",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: false},
				{address: addr1, isValidPredicate: false},
			},
			expectedRes: map[common.Address][]byte{addr1: {3}},
		},
		{
			name: "one address two mixed predicates",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: true},
				{address: addr1, isValidPredicate: false},
			},
			expectedRes: map[common.Address][]byte{addr1: {2}},
		},
		{
			name: "one address mixed predicates",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: true},
				{address: addr1, isValidPredicate: false},
				{address: addr1, isValidPredicate: false},
				{address: addr1, isValidPredicate: true},
			},
			expectedRes: map[common.Address][]byte{addr1: {6}},
		},
		{
			name: "two addresses mixed predicates",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: true},
				{address: addr2, isValidPredicate: false},
				{address: addr1, isValidPredicate: false},
				{address: addr1, isValidPredicate: false},
				{address: addr2, isValidPredicate: true},
				{address: addr2, isValidPredicate: true},
				{address: addr2, isValidPredicate: false},
				{address: addr2, isValidPredicate: true},
			},
			expectedRes: map[common.Address][]byte{addr1: {6}, addr2: {9}},
		},
		{
			name: "two addresses all valid predicates",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: true},
				{address: addr2, isValidPredicate: true},
				{address: addr1, isValidPredicate: true},
				{address: addr1, isValidPredicate: true},
			},
			expectedRes: map[common.Address][]byte{addr1: {}, addr2: {}},
		},
		{
			name: "two addresses all invalid predicates",
			testTuple: []testTuple{
				{address: addr1, isValidPredicate: false},
				{address: addr2, isValidPredicate: false},
				{address: addr1, isValidPredicate: false},
				{address: addr1, isValidPredicate: false},
			},
			expectedRes: map[common.Address][]byte{addr1: {7}, addr2: {1}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			// Create the rules from TestChainConfig and update the predicates based on the test params
			rules := params.TestChainConfig.AvalancheRules(common.Big0, 0)
			predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
			predicater.EXPECT().PredicateGas(gomock.Any()).Return(uint64(0), nil).Times(len(test.testTuple))
			validPredicateCount := 0

			var txAccessList types.AccessList
			for _, tuple := range test.testTuple {
				predicateHash := invalidHash
				if tuple.isValidPredicate {
					validPredicateCount++
					predicateHash = validHash
				}
				txAccessList = append(txAccessList, types.AccessTuple{
					Address: tuple.address,
					StorageKeys: []common.Hash{
						predicateHash,
					},
				})
			}

			invalidPredicateCount := len(test.testTuple) - validPredicateCount
			predicater.EXPECT().VerifyPredicate(gomock.Any(), validHash[:]).Return(true).Times(validPredicateCount)
			predicater.EXPECT().VerifyPredicate(gomock.Any(), invalidHash[:]).Return(false).Times(invalidPredicateCount)
			rules.Predicaters[addr1] = predicater
			rules.Predicaters[addr2] = predicater

			// Specify only the access list, since this test should not depend on any other values
			tx := types.NewTx(&types.DynamicFeeTx{
				AccessList: txAccessList,
				Gas:        53000,
			})

			oldPredicateRes, err := CheckPredicatesTest(predicateContext, tx)
			require.NoError(err)
			require.Equal(test.expectedRes, oldPredicateRes)

			predicateRes, err := CheckPredicates(rules, predicateContext, tx)
			require.NoError(err)
			require.Equal(test.expectedRes, predicateRes)
		})
	}
}

func VerifyPredicateTest(predicateContext *precompileconfig.PredicateContext, predicates [][]byte) []byte {
	resultBitSet := set.NewBits()

	for predicateIndex, predicateBytes := range predicates {
		if bytes.Equal(predicateBytes, invalidHash[:]) {
			resultBitSet.Add(predicateIndex)
		}
	}
	return resultBitSet.Bytes()
}

func CheckPredicatesTest(predicateContext *precompileconfig.PredicateContext, tx *types.Transaction) (map[common.Address][]byte, error) {
	// Check that the transaction can cover its IntrinsicGas (including the gas required by the predicate) before
	// verifying the predicate.

	predicateResults := make(map[common.Address][]byte)
	// Short circuit early if there are no precompile predicates to verify

	// Prepare the predicate storage slots from the transaction's access list
	predicateArguments := PreparePredicateStorageSlotsTest(tx.AccessList())

	// If there are no predicates to verify, return early and skip requiring the proposervm block
	// context to be populated.
	if len(predicateArguments) == 0 {
		return predicateResults, nil
	}

	if predicateContext == nil || predicateContext.ProposerVMBlockCtx == nil {
		return nil, ErrMissingPredicateContext
	}

	for address, predicates := range predicateArguments {
		// Since [address] is only added to [predicateArguments] when there's a valid predicate in the ruleset
		// there's no need to check if the predicate exists here.
		res := VerifyPredicateTest(predicateContext, predicates)
		log.Debug("predicate verify", "tx", tx.Hash(), "address", address, "res", res)
		predicateResults[address] = res
	}

	return predicateResults, nil
}

func PreparePredicateStorageSlotsTest(list types.AccessList) map[common.Address][][]byte {
	predicateStorageSlots := make(map[common.Address][][]byte)
	for _, el := range list {
		predicateStorageSlots[el.Address] = append(predicateStorageSlots[el.Address], utils.HashSliceToBytes(el.StorageKeys))
	}

	return predicateStorageSlots
}

func setupCheckPredicatesOutput(t *testing.B, addrSize int, tupleSizePerAddr int) (params.Rules, *precompileconfig.PredicateContext, *types.Transaction) {
	type testTuple struct {
		address          common.Address
		isValidPredicate bool
	}

	predicater := precompileconfig.NewMockPredicater(gomock.NewController(t))
	predicater.EXPECT().VerifyPredicate(gomock.Any(), validHash[:]).Return(true).AnyTimes()
	predicater.EXPECT().VerifyPredicate(gomock.Any(), invalidHash[:]).Return(false).AnyTimes()
	predicater.EXPECT().PredicateGas(gomock.Any()).Return(uint64(0), nil).AnyTimes()
	rules := params.TestChainConfig.AvalancheRules(common.Big0, 0)

	testTuples := make([]testTuple, 0)
	for i := 0; i < addrSize; i++ {
		bigIndex := big.NewInt(int64(i))
		addr := common.BigToAddress(bigIndex)
		rules.Predicaters[addr] = predicater
		for k := 0; k < tupleSizePerAddr; k++ {
			testTuples = append(testTuples, testTuple{
				address:          addr,
				isValidPredicate: k%2 == 0,
			})
		}
	}

	predicateContext := &precompileconfig.PredicateContext{
		ProposerVMBlockCtx: &block.Context{
			PChainHeight: 10,
		},
	}

	var txAccessList types.AccessList
	for _, tuple := range testTuples {
		// Create the rules from TestChainConfig and update the predicates based on the test params
		predicateHash := invalidHash
		if tuple.isValidPredicate {
			predicateHash = validHash
		}
		txAccessList = append(txAccessList, types.AccessTuple{
			Address: tuple.address,
			StorageKeys: []common.Hash{
				predicateHash,
			}})
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		AccessList: txAccessList,
		Gas:        53000,
	})

	return rules, predicateContext, tx
}

func runBenchmarkCheckPredicates(b *testing.B, rules params.Rules, predicateContext *precompileconfig.PredicateContext, tx *types.Transaction, useOld bool) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var predicateRes map[common.Address][]byte
		var err error
		if useOld {
			predicateRes, err = CheckPredicatesOld(rules, predicateContext, tx)
		} else {
			predicateRes, err = CheckPredicates(rules, predicateContext, tx)
		}
		require.NoError(b, err)
		require.NotNil(b, predicateRes)
	}
}

func BenchmarkPredicateCheck(b *testing.B) {
	// addrSize, predicatePerAddrNum
	addrSizes := []int{1, 10, 100, 1000}
	predicatePerAddr := []int{1, 10, 100, 1000}

	for _, addrSize := range addrSizes {
		for _, predicatePerAddrSize := range predicatePerAddr {
			b.Run(fmt.Sprintf("addressNum_%d_predicatePerAddr_%d_predicateCheck", addrSize, predicatePerAddrSize), func(b *testing.B) {
				rules, context, tx := setupCheckPredicatesOutput(b, addrSize, predicatePerAddrSize)
				runBenchmarkCheckPredicates(b, rules, context, tx, false)
			})
		}
	}
}

func BenchmarkPredicateCheckOld(b *testing.B) {
	// addrSize, predicatePerAddrNum
	addrSizes := []int{1, 10, 100, 1000}
	predicatePerAddr := []int{1, 10, 100, 1000}

	for _, addrSize := range addrSizes {
		for _, predicatePerAddrSize := range predicatePerAddr {
			b.Run(fmt.Sprintf("addressNum_%d_predicatePerAddr_%d_predicateCheck", addrSize, predicatePerAddrSize), func(b *testing.B) {
				rules, context, tx := setupCheckPredicatesOutput(b, addrSize, predicatePerAddrSize)
				runBenchmarkCheckPredicates(b, rules, context, tx, true)
			})
		}
	}
}
