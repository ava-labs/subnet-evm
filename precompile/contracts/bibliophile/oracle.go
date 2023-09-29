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

	AGGREGATOR_MAP_SLOT    int64 = 1
	RED_STONE_ADAPTER_SLOT int64 = 2
)

const (
	// this slot is from TestOracle.sol
	TEST_ORACLE_PRICES_MAPPING_SLOT int64 = 3
)

func getUnderlyingPrice(stateDB contract.StateDB, market common.Address) *big.Int {
	return getUnderlyingPrice_(stateDB, getUnderlyingAssetAddress(stateDB, market))
}

func getUnderlyingPrice_(stateDB contract.StateDB, underlying common.Address) *big.Int {
	oracle := getOracleAddress(stateDB) // this comes from margin account
	feedId := getRedStoneFeedId(stateDB, oracle, underlying)
	if feedId.Big().Sign() != 0 {
		// redstone oracle is configured for this market
		redStoneAdapter := getRedStoneAdapterAddress(stateDB, oracle)
		redstonePrice := getRedStonePrice(stateDB, redStoneAdapter, feedId)
		// log.Info("redstone-price", "amm", market, "price", redstonePrice)
		return redstonePrice
	}
	// red stone oracle is not enabled for this market, we use the default TestOracle
	slot := crypto.Keccak256(append(common.LeftPadBytes(underlying.Bytes(), 32), common.BigToHash(big.NewInt(TEST_ORACLE_PRICES_MAPPING_SLOT)).Bytes()...))
	return fromTwosComplement(stateDB.GetState(oracle, common.BytesToHash(slot)).Bytes())
}

func getRedStoneAdapterAddress(stateDB contract.StateDB, oracle common.Address) common.Address {
	return common.BytesToAddress(stateDB.GetState(oracle, common.BigToHash(big.NewInt(RED_STONE_ADAPTER_SLOT))).Bytes())
}

func getRedStonePrice(stateDB contract.StateDB, adapterAddress common.Address, redStoneFeedId common.Hash) *big.Int {
	latestRoundId := getlatestRoundId(stateDB, adapterAddress)
	slot := common.BytesToHash(crypto.Keccak256(append(append(redStoneFeedId.Bytes(), common.LeftPadBytes(latestRoundId.Bytes(), 32)...), RED_STONE_VALUES_MAPPING_STORAGE_LOCATION.Bytes()...)))
	return new(big.Int).Div(fromTwosComplement(stateDB.GetState(adapterAddress, slot).Bytes()), big.NewInt(100)) // we use 6 decimals precision everywhere
}

func getlatestRoundId(stateDB contract.StateDB, adapterAddress common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(adapterAddress, RED_STONE_LATEST_ROUND_ID_STORAGE_LOCATION).Bytes())
}

func getRedStoneFeedId(stateDB contract.StateDB, oracle, underlying common.Address) common.Hash {
	aggregatorMapSlot := crypto.Keccak256(append(common.LeftPadBytes(underlying.Bytes(), 32), common.BigToHash(big.NewInt(AGGREGATOR_MAP_SLOT)).Bytes()...))
	return stateDB.GetState(oracle, common.BytesToHash(aggregatorMapSlot))
}
