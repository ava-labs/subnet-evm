// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package predicate

import (
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

// GetPredicatesFromAccessList returns the predicates associated with the address in the access list
// Note: if an address is specified multiple times in the access list, each storage slot for that address is
// appended to a slice of byte slices. Each byte slice represents a predicate, making it a slice of predicates
// for each access list address, and every predicate in the slice goes through verification.
func GetPredicatesFromAccessList(list types.AccessList, address common.Address) [][]byte {
	var predicates [][]byte
	for _, el := range list {
		if el.Address == address {
			predicates = append(predicates, utils.HashSliceToBytes(el.StorageKeys))
		}
	}
	return predicates
}
