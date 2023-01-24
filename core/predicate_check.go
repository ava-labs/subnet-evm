package core

import (
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/log"
)

// CheckPredicatesForSenderTxs checks the stateful precompile predicates of any of the given
// transactions that reference a stateful precompile address in their access list.
// The parameter [txs] represents a flattened, nonce-ordered list of transactions originating from the same sender
// Returns [len(txs), nil] if and only if all referenced predicates are met.
// Otherwise, returns the index of the first transaction that was invalid and the predicate error.
// In the failure case, all transactions after and including the index that failed the predicate should be considered invalid,
// since the input [txs] are nonce-ordered and from the same sender.
func CheckPredicatesForSenderTxs(rules params.Rules, snowCtx *snow.Context, txs types.Transactions, proposerVMBlockCtx *block.Context) (int, error) {
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

			if err := predicate(snowCtx, proposerVMBlockCtx, utils.HashSliceToBytes(accessTuple.StorageKeys)); err != nil {
				log.Debug("Transaction predicate verification failed.", "txId", tx.Hash(), "precompileAddress", accessTuple.Address.Hex())
				return i, err
			}
		}
	}

	return len(txs), nil
}
