// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/log"
)

// CheckPredicatesForSenderTxs checks the stateful precompile predicates of any of the given
// transactions that reference a stateful precompile address in their access list.
// The parameter [txs] represents a flattened list of transactions. It is up to the caller to decide the ordering of
// the transactions in [txs] in relation to the return value of this function.
// Returns [len(txs), nil] if and only if all referenced predicates are met.
// Otherwise, returns the index of the first transaction that was invalid and the predicate error.
func CheckPredicatesForSenderTxs(rules params.Rules, predicateContext *precompile.PredicateContext, txs types.Transactions) (int, error) {
	precompileConfigs := rules.Precompiles
	for i, tx := range txs {
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
				log.Debug("Transaction predicate verification failed.", "txHash", tx.Hash(), "precompileAddress", accessTuple.Address.Hex())
				return i, err
			}
		}
	}

	return len(txs), nil
}
