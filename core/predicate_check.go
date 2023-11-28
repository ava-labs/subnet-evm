// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"errors"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var ErrMissingPredicateContext = errors.New("missing predicate context")

// CheckPredicates verifies the predicates of [tx] and returns the result. Returning an error invalidates the block.
func CheckPredicates(rules params.Rules, predicateContext *precompileconfig.PredicateContext, tx *types.Transaction) (map[common.Address][]byte, error) {
	// Check that the transaction can cover its IntrinsicGas (including the gas required by the predicate) before
	// verifying the predicate.
	intrinsicGas, err := IntrinsicGas(tx.Data(), tx.AccessList(), tx.To() == nil, rules)
	if err != nil {
		return nil, err
	}
	if tx.Gas() < intrinsicGas {
		return nil, fmt.Errorf("%w for predicate verification (%d) < intrinsic gas (%d)", ErrIntrinsicGas, tx.Gas(), intrinsicGas)
	}

	predicateResults := make(map[common.Address][]byte)
	// Short circuit early if there are no precompile predicates to verify
	if !rules.PredicatersExist() {
		return predicateResults, nil
	}

	predicateIndexes := make(map[common.Address]int)
	for _, al := range tx.AccessList() {
		address := al.Address
		predicaterContract, exists := rules.Predicaters[address]
		if !exists {
			continue
		}
		// Invariant: We should return this error only if there is a predicate in txs.
		// If there is no predicate in txs, we should just return an empty result with no error.
		if predicateContext == nil || predicateContext.ProposerVMBlockCtx == nil {
			return nil, ErrMissingPredicateContext
		}
		verified := predicaterContract.VerifyPredicate(predicateContext, utils.HashSliceToBytes(al.StorageKeys))
		log.Debug("predicate verify", "tx", tx.Hash(), "address", address, "verified", verified)
		// Add bitset only if predicate is not verified
		if !verified {
			resultBitSet := set.NewBits()
			currentResult, ok := predicateResults[address]
			if ok {
				resultBitSet = set.BitsFromBytes(currentResult)
			}
			// this will default to 0 if the address is not in the map
			currentIndex := predicateIndexes[address]
			resultBitSet.Add(currentIndex)
			predicateResults[address] = resultBitSet.Bytes()
		}
		// add an empty byte to indicate that the predicate was verified
		// for the address
		if _, ok := predicateResults[address]; !ok {
			predicateResults[address] = []byte{}
		}
		predicateIndexes[address]++
	}

	return predicateResults, nil
}
