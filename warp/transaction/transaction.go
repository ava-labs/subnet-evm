// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package transaction

import (
	"math/big"

	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	byteUtils "github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

func NewWarpTx(
	chainID *big.Int,
	nonce uint64,
	to *common.Address,
	gas uint64,
	gasFeeCap *big.Int,
	gasTipCap *big.Int,
	value *big.Int,
	data []byte,
	accessList types.AccessList,
	signedMessage *avalancheWarp.Message,
) *types.Transaction {
	accessList = append(accessList, types.AccessTuple{
		Address:     warp.ContractAddress,
		StorageKeys: byteUtils.BytesToHashSlice(byteUtils.PackPredicate(signedMessage.Bytes())),
	})
	return types.NewTx(&types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		To:         to,
		Gas:        gas,
		GasFeeCap:  gasFeeCap,
		GasTipCap:  gasTipCap,
		Value:      value,
		Data:       data,
		AccessList: accessList,
	})
}
