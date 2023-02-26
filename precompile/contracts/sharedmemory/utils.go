// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sharedmemory

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CalculateANTAssetID calculates the hash of the caller address concatenated with the blockchainID, which
// serves as the assetID of the token that can be minted on this blockchain by the specified caller.
func CalculateANTAssetID(blockchainID common.Hash, caller common.Address) common.Hash {
	assetID := crypto.Keccak256Hash(blockchainID[:], caller[:])
	return assetID
}
