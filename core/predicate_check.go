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
)

// CheckPredicates checks that all precompile predicates are satisfied within the current [predicateContext] for [tx]
func CheckPredicates(rules params.Rules, predicateContext *contract.PredicateContext, tx *types.Transaction) error {
	precompileConfigs := rules.ActivePrecompiles
	for _, accessTuple := range tx.AccessList() {
		_, isPrecompile := precompileConfigs[accessTuple.Address]
		if !isPrecompile {
			continue
		}

		module, ok := modules.GetPrecompileModuleByAddress(accessTuple.Address)
		if !ok {
			return fmt.Errorf("accessed precompile config under address %s with no registered module", accessTuple.Address)
		}
		predicater, ok := module.Contract.(contract.Predicater)
		if !ok {
			continue
		}

		if err := predicater.VerifyPredicate(predicateContext, utils.HashSliceToBytes(accessTuple.StorageKeys)); err != nil {
			return fmt.Errorf("predicate %s failed verification for tx %s: %w", accessTuple.Address, tx.Hash(), err)
		}
	}

	return nil
}
