// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package eth

import "github.com/ava-labs/subnet-evm/internal/ethapi"

type (
	TransactionArgs = ethapi.TransactionArgs
	Backend         = ethapi.Backend
	StateOverride   = ethapi.StateOverride
	BlockOverrides  = ethapi.BlockOverrides
	RPCTransaction  = ethapi.RPCTransaction
)

var (
	CreateAccessList = ethapi.CreateAccessList
	NewNetAPI        = ethapi.NewNetAPI
	RPCMarshalBlock  = ethapi.RPCMarshalBlock
)
