package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	RED_STONE_VALUES_MAPPING_STORAGE_LOCATION  = common.HexToHash("0x4dd0c77efa6f6d590c97573d8c70b714546e7311202ff7c11c484cc841d91bfc") // keccak256("RedStone.oracleValuesMapping");
	RED_STONE_LATEST_ROUND_ID_STORAGE_LOCATION = common.HexToHash("0xc68d7f1ee07d8668991a8951e720010c9d44c2f11c06b5cac61fbc4083263938") // keccak256("RedStone.latestRoundId");
)

func getRedStonePrice(stateDB contract.StateDB, adapterAddress common.Address, redStoneFeedId common.Hash) *big.Int {
	latestRoundId := getlatestRoundId(stateDB, adapterAddress)
	slot := common.BytesToHash(crypto.Keccak256(append(append(redStoneFeedId.Bytes(), common.LeftPadBytes(latestRoundId.Bytes(), 32)...), RED_STONE_VALUES_MAPPING_STORAGE_LOCATION.Bytes()...)))
	return new(big.Int).Div(fromTwosComplement(stateDB.GetState(adapterAddress, slot).Bytes()), big.NewInt(100)) // we use 6 decimals precision everywhere
}

func getlatestRoundId(stateDB contract.StateDB, adapterAddress common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(adapterAddress, RED_STONE_LATEST_ROUND_ID_STORAGE_LOCATION).Bytes())
}
