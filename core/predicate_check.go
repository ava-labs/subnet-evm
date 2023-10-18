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
	if len(rules.Predicaters) == 0 {
		return predicateResults, nil
	}

	// Prepare the predicate storage slots from the transaction's access list
	for _, el := range tx.AccessList() {
		predicaterContract, exists := rules.Predicaters[el.Address]
		if !exists {
			continue
		}
		if predicateContext == nil || predicateContext.ProposerVMBlockCtx == nil {
			return nil, ErrMissingPredicateContext
		}
		verified := predicaterContract.VerifyPredicate(predicateContext, utils.HashSliceToBytes(el.StorageKeys))
		// Add bitset only if predicate is not verified
		if !verified {
			resultBitSet := set.NewBits()
			currentResult, ok := predicateResults[el.Address]
			if ok {
				resultBitSet = set.BitsFromBytes(currentResult)
			}
			currentIndex := resultBitSet.Len()
			resultBitSet.Add(currentIndex)
			predicateResults[el.Address] = resultBitSet.Bytes()
		}
		// fill with empty result
		if _, ok := predicateResults[el.Address]; !ok {
			predicateResults[el.Address] = []byte{}
		}
	}

	return predicateResults, nil
}
