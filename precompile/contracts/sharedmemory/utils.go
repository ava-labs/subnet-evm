// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sharedmemory

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CalculateANTAssetID calculates the hash of the caller address concatenated with the blockchainID, which
// serves as the assetID of the token that can be minted on this blockchain by the specified caller.
func CalculateANTAssetID(blockchainID common.Hash, caller common.Address) common.Hash {
	assetID := crypto.Keccak256Hash(blockchainID[:], caller[:])
	return assetID
}

func GetNamedUTXOs(tx *types.Transaction) ([]ids.ID, error) {
	namedUTXOs := make([]ids.ID, 0)
	for _, accessTuple := range tx.AccessList() {
		address := accessTuple.Address
		if address != ContractAddress {
			continue
		}
		predicateBytes := utils.HashSliceToBytes(accessTuple.StorageKeys)
		predicateBytes, err := utils.UnpackPredicate(predicateBytes)
		if err != nil {
			return nil, fmt.Errorf("predicate %s failed unpacking for tx %s: %w", address, tx.Hash(), err)
		}

		atomicPredicate := new(AtomicPredicate)
		version, err := codec.Codec.Unmarshal(predicateBytes, atomicPredicate)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal shared memory predicate: %w", err)
		}
		if version != codec.CodecVersion {
			return nil, fmt.Errorf("invalid version for shared memory predicate: %d", version)
		}

		for _, utxo := range atomicPredicate.ImportedUTXOs {
			namedUTXOs = append(namedUTXOs, utxo.ID)
		}
	}
	return namedUTXOs, nil
}
