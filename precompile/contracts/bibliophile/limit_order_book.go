package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	LIMIT_ORDERBOOK_GENESIS_ADDRESS       = "0x0300000000000000000000000000000000000005"
	ORDER_INFO_SLOT                 int64 = 1
	REDUCE_ONLY_AMOUNT_SLOT         int64 = 2
	LONG_OPEN_ORDERS_SLOT           int64 = 4
	SHORT_OPEN_ORDERS_SLOT          int64 = 5
)

func getOrderFilledAmount(stateDB contract.StateDB, orderHash [32]byte) *big.Int {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	num := stateDB.GetState(common.HexToAddress(LIMIT_ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(new(big.Int).Add(orderInfo, big.NewInt(1)))).Bytes()
	return fromTwosComplement(num)
}

func getOrderStatus(stateDB contract.StateDB, orderHash [32]byte) int64 {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	return new(big.Int).SetBytes(stateDB.GetState(common.HexToAddress(LIMIT_ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(new(big.Int).Add(orderInfo, big.NewInt(3)))).Bytes()).Int64()
}

func orderInfoMappingStorageSlot(orderHash [32]byte) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(append(orderHash[:], common.LeftPadBytes(big.NewInt(ORDER_INFO_SLOT).Bytes(), 32)...)))
}

func getReduceOnlyAmount(stateDB contract.StateDB, trader common.Address, ammIndex *big.Int) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(REDUCE_ONLY_AMOUNT_SLOT).Bytes(), 32)...))
	nestedMappingHash := crypto.Keccak256(append(common.LeftPadBytes(ammIndex.Bytes(), 32), baseMappingHash...))
	return fromTwosComplement(stateDB.GetState(common.HexToAddress(LIMIT_ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(nestedMappingHash)).Bytes())
}

func getLongOpenOrdersAmount(stateDB contract.StateDB, trader common.Address, ammIndex *big.Int) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(LONG_OPEN_ORDERS_SLOT).Bytes(), 32)...))
	nestedMappingHash := crypto.Keccak256(append(common.LeftPadBytes(ammIndex.Bytes(), 32), baseMappingHash...))
	return stateDB.GetState(common.HexToAddress(LIMIT_ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(nestedMappingHash)).Big()
}

func getShortOpenOrdersAmount(stateDB contract.StateDB, trader common.Address, ammIndex *big.Int) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(SHORT_OPEN_ORDERS_SLOT).Bytes(), 32)...))
	nestedMappingHash := crypto.Keccak256(append(common.LeftPadBytes(ammIndex.Bytes(), 32), baseMappingHash...))
	return stateDB.GetState(common.HexToAddress(LIMIT_ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(nestedMappingHash)).Big()
}

func getBlockPlaced(stateDB contract.StateDB, orderHash [32]byte) *big.Int {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	return new(big.Int).SetBytes(stateDB.GetState(common.HexToAddress(LIMIT_ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(orderInfo)).Bytes())
}
