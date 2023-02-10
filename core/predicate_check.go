// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/utils"
)

// CheckPredicates checks that all precompile predicates are satisifed within the current [predicateContext] for [tx]
func CheckPredicates(rules params.Rules, predicateContext *precompile.PredicateContext, tx *types.Transaction) error {
	precompileConfigs := rules.Precompiles
	for _, accessTuple := range tx.AccessList() {
		precompileConfig, isPrecompileAccess := precompileConfigs[accessTuple.Address]
		if !isPrecompileAccess {
			continue
		}

		predicate := precompileConfig.Predicate()
		if predicate == nil {
			continue
		}

		if err := predicate(predicateContext, utils.HashSliceToBytes(accessTuple.StorageKeys)); err != nil {
			return fmt.Errorf("predicate %s failed verification for tx %s: %w", accessTuple.Address, tx.Hash(), err)
		}
	}

	return nil
}
