// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

// CheckPredicates checks that all precompile predicates are satisfied within the current [predicateContext] for [tx]
func CheckPredicates(rules params.Rules, predicateContext *contract.PredicateContext, tx *types.Transaction) error {
	precompileConfigs := rules.ActivePrecompiles
	// Track addresses that we've performed a predicate check for
	precompileAddressChecks := make(map[common.Address]struct{})
	for _, accessTuple := range tx.AccessList() {
		address := accessTuple.Address
		_, isPrecompile := precompileConfigs[address]
		if !isPrecompile {
			continue
		}

		module, ok := modules.GetPrecompileModuleByAddress(address)
		if !ok {
			return fmt.Errorf("predicate accessed precompile config under address %s with no registered module for tx %s", address, tx.Hash())
		}
		predicater, ok := module.Contract.(contract.Predicater)
		if !ok {
			continue
		}
		// Return an error if we've already checked a predicate for this address
		if _, ok := precompileAddressChecks[address]; ok {
			return fmt.Errorf("predicate %s failed verification for tx %s: specified %s in access list multiple times", address, tx.Hash(), address)
		}
		precompileAddressChecks[address] = struct{}{}
		
		if err := predicater.VerifyPredicate(predicateContext, utils.HashSliceToBytes(accessTuple.StorageKeys)); err != nil {
			return fmt.Errorf("predicate %s failed verification for tx %s: %w", address, tx.Hash(), err)
		}
	}

	return nil
}
