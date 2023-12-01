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
	"github.com/ava-labs/subnet-evm/predicate"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
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
		//	Add bitset only if predicate is not verified
		resultBitSet := set.NewBits()
		if !verified {
			currentResult, ok := predicateResults[address]
			if ok {
				resultBitSet = set.BitsFromBytes(currentResult)
			}
			// this will default to 0 if the address is not in the map
			currentIndex := predicateIndexes[address]
			resultBitSet.Add(currentIndex)
		}
		// add an empty byte to indicate that the predicate was verified
		// for the address
		res := resultBitSet.Bytes()
		predicateResults[address] = res
		//	log.Debug("predicate verify", "tx", tx.Hash(), "address", address, "res", res)
		predicateIndexes[address]++
	}

	return predicateResults, nil
}

func CheckPredicatesOld(rules params.Rules, predicateContext *precompileconfig.PredicateContext, tx *types.Transaction) (map[common.Address][]byte, error) {
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

	// Prepare the predicate storage slots from the transaction's access list
	predicateArguments := predicate.PreparePredicateStorageSlots(rules, tx.AccessList())

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
		predicaterContract := rules.Predicaters[address]
		bitset := set.NewBits()
		for i, predicate := range predicates {
			if !predicaterContract.VerifyPredicate(predicateContext, predicate) {
				bitset.Add(i)
			}
		}
		res := bitset.Bytes()
		//		log.Debug("predicate verify", "tx", tx.Hash(), "address", address, "res", res)
		predicateResults[address] = res
	}
	return predicateResults, nil
}
