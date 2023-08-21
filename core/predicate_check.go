// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// CheckPredicates verifies the predicates of [tx] and returns the result. Returning an error invalidates the block.
func CheckPredicates(rules params.Rules, predicateContext *precompileconfig.PredicateContext, tx *types.Transaction) (map[common.Address][]byte, error) {
	// Check that the transaction can cover its IntrinsicGas (including the gas required by the predicate) before
	// verifying the predicate.
	intrinsicGas, err := IntrinsicGas(tx.Data(), tx.AccessList(), tx.To() == nil, rules)
	if err != nil {
		return nil, err
	}
	if tx.Gas() < intrinsicGas {
		return nil, fmt.Errorf("insufficient gas for predicate verification (%d) < intrinsic gas (%d)", tx.Gas(), intrinsicGas)
	}
	return checkPrecompilePredicates(rules, predicateContext, tx)
}

func checkPrecompilePredicates(rules params.Rules, predicateContext *precompileconfig.PredicateContext, tx *types.Transaction) (map[common.Address][]byte, error) {
	predicateResults := make(map[common.Address][]byte)
	// Short circuit early if there are no precompile predicates to verify
	if len(rules.PredicatePrecompiles) == 0 {
		return predicateResults, nil
	}
	precompilePredicates := rules.PredicatePrecompiles
	predicateArguments := make(map[common.Address][][]byte)
	for _, accessTuple := range tx.AccessList() {
		address := accessTuple.Address
		_, ok := precompilePredicates[address]
		if !ok {
			continue
		}

		predicateArguments[address] = append(predicateArguments[address], predicateutils.HashSliceToBytes(accessTuple.StorageKeys))
	}

	for address, predicates := range predicateArguments {
		predicate := precompilePredicates[address]
		res := predicate.VerifyPredicate(predicateContext, predicates)
		log.Debug("predicate verify", "tx", tx.Hash(), "address", address, "res", res)
		predicateResults[address] = res
	}

	return predicateResults, nil
}
