package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	SIGNED_ORDER_INFO_SLOT int64 = 53
)

// State Reader
func GetSignedOrderFilledAmount(stateDB contract.StateDB, orderHash [32]byte) *big.Int {
	orderInfo := signedOrderInfoMappingStorageSlot(orderHash)
	num := stateDB.GetState(GetSignedOrderBookAddress(stateDB), common.BigToHash(orderInfo)).Bytes()
	return fromTwosComplement(num)
}

func GetSignedOrderStatus(stateDB contract.StateDB, orderHash [32]byte) int64 {
	a := GetSignedOrderBookAddress(stateDB)
	orderInfo := signedOrderInfoMappingStorageSlot(orderHash)
	return new(big.Int).SetBytes(stateDB.GetState(a, common.BigToHash(new(big.Int).Add(orderInfo, big.NewInt(1)))).Bytes()).Int64()
}

func signedOrderInfoMappingStorageSlot(orderHash [32]byte) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(append(orderHash[:], common.LeftPadBytes(big.NewInt(SIGNED_ORDER_INFO_SLOT).Bytes(), 32)...)))
}

func GetSignedOrderBookAddress(stateDB contract.StateDB) common.Address {
	slot := crypto.Keccak256(append(common.LeftPadBytes(big.NewInt(2).Bytes() /* orderType */, 32), common.LeftPadBytes(big.NewInt(ORDER_HANDLER_STORAGE_SLOT).Bytes(), 32)...))
	return common.BytesToAddress(stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(slot)).Bytes())
}
