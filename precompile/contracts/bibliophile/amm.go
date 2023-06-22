package bibliophile

import (
	"math/big"
	"strings"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

const (
	MARK_PRICE_TWAP_DATA_SLOT       int64 = 1
	VAR_POSITIONS_SLOT              int64 = 5
	VAR_CUMULATIVE_PREMIUM_FRACTION int64 = 6
	MAX_ORACLE_SPREAD_RATIO_SLOT    int64 = 7
	MAX_LIQUIDATION_RATIO_SLOT      int64 = 8
	MIN_SIZE_REQUIREMENT_SLOT       int64 = 9
	ORACLE_SLOT                     int64 = 10
	UNDERLYING_ASSET_SLOT           int64 = 11
	MAX_LIQUIDATION_PRICE_SPREAD    int64 = 17
	RED_STONE_ADAPTER_SLOT          int64 = 21
	RED_STONE_FEED_ID_SLOT          int64 = 22
)

const (
	TEST_ORACLE_PRICES_MAPPING_SLOT int64 = 53
)

var (
	// Date and time (GMT): riday, 9 June 2023 14:40:00
	V2ActivationDate *big.Int = new(big.Int).SetInt64(1686321600)
)

// AMM State
func getLastPrice(stateDB contract.StateDB, market common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(MARK_PRICE_TWAP_DATA_SLOT))).Big()
}

func GetCumulativePremiumFraction(stateDB contract.StateDB, market common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(VAR_CUMULATIVE_PREMIUM_FRACTION))).Bytes())
}

// GetMaxOraclePriceSpread returns the maxOracleSpreadRatio for a given market
func GetMaxOraclePriceSpread(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_ORACLE_SPREAD_RATIO_SLOT))).Bytes())
}

// GetMaxLiquidationPriceSpread returns the maxOracleSpreadRatio for a given market
func GetMaxLiquidationPriceSpread(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_LIQUIDATION_PRICE_SPREAD))).Bytes())
}

// GetMaxLiquidationRatio returns the maxLiquidationPriceSpread for a given market
func GetMaxLiquidationRatio(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_LIQUIDATION_RATIO_SLOT))).Bytes())
}

// GetMinSizeRequirement returns the minSizeRequirement for a given market
func GetMinSizeRequirement(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MIN_SIZE_REQUIREMENT_SLOT))).Bytes())
}

func getOracleAddress(stateDB contract.StateDB, market common.Address) common.Address {
	return common.BytesToAddress(stateDB.GetState(market, common.BigToHash(big.NewInt(ORACLE_SLOT))).Bytes())
}

func getUnderlyingAssetAddress(stateDB contract.StateDB, market common.Address) common.Address {
	return common.BytesToAddress(stateDB.GetState(market, common.BigToHash(big.NewInt(UNDERLYING_ASSET_SLOT))).Bytes())
}

func getUnderlyingPriceForMarket(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return getUnderlyingPrice(stateDB, market)
}

func getRedStoneAdapterAddress(stateDB contract.StateDB, market common.Address) common.Address {
	return common.BytesToAddress(stateDB.GetState(market, common.BigToHash(big.NewInt(RED_STONE_ADAPTER_SLOT))).Bytes())
}

func getRedStoneFeedId(stateDB contract.StateDB, market common.Address) common.Hash {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(RED_STONE_FEED_ID_SLOT)))
}

func getUnderlyingPrice(stateDB contract.StateDB, market common.Address) *big.Int {
	redStoneAdapter := getRedStoneAdapterAddress(stateDB, market)
	if redStoneAdapter.Hash().Big().Sign() != 0 {
		feedId := getRedStoneFeedId(stateDB, market)
		// first we check the feedId, if it is set, it should imply we are using a redstone oracle
		// log.Info("red-stone-feed-id", "feedId", feedId.String())
		if feedId.Big().Sign() != 0 {
			// redstone oracle is configured for this market
			redstonePrice := getRedStonePrice(stateDB, redStoneAdapter, feedId)
			log.Info("redstone-price", "amm", market, "price", redstonePrice)
			return redstonePrice
		} else {
			// just log the red stone price, for testing before deployment
			var feedId common.Hash
			if strings.EqualFold(market.String(), "0xa72b463C21dA61cCc86069cFab82e9e8491152a0") { // eth amm
				feedId = common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000")
			} else if strings.EqualFold(market.String(), "0xd80e57dB448b0692C396B890eE9c791D7386dAdC") { // avax amm
				feedId = common.HexToHash("0x4156415800000000000000000000000000000000000000000000000000000000")
			}
			redstonePrice := getRedStonePrice(stateDB, redStoneAdapter, feedId)
			log.Info("log-only-redstone-price", "amm", market, "price", redstonePrice)
		}
	}
	// red stone oracle is not enabled for this market, we use the default TestOracle
	oracle := getOracleAddress(stateDB, market)
	underlying := getUnderlyingAssetAddress(stateDB, market)
	slot := crypto.Keccak256(append(common.LeftPadBytes(underlying.Bytes(), 32), common.LeftPadBytes(big.NewInt(TEST_ORACLE_PRICES_MAPPING_SLOT).Bytes(), 32)...))
	return fromTwosComplement(stateDB.GetState(oracle, common.BytesToHash(slot)).Bytes())
}

// Trader State

func positionsStorageSlot(trader *common.Address) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(VAR_POSITIONS_SLOT).Bytes(), 32)...)))
}

func getSize(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(positionsStorageSlot(trader))).Bytes())
}

func getOpenNotional(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(new(big.Int).Add(positionsStorageSlot(trader), big.NewInt(1)))).Big()
}

func GetLastPremiumFraction(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(new(big.Int).Add(positionsStorageSlot(trader), big.NewInt(2)))).Bytes())
}

// Utils

func getPendingFundingPayment(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	cumulativePremiumFraction := GetCumulativePremiumFraction(stateDB, market)
	return divide1e18(new(big.Int).Mul(new(big.Int).Sub(cumulativePremiumFraction, GetLastPremiumFraction(stateDB, market, trader)), getSize(stateDB, market, trader)))
}

func getOptimalPnl(stateDB contract.StateDB, market common.Address, oraclePrice *big.Int, lastPrice *big.Int, trader *common.Address, margin *big.Int, marginMode MarginMode, blockTimestamp *big.Int) (notionalPosition *big.Int, uPnL *big.Int) {
	size := getSize(stateDB, market, trader)
	if size.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0)
	}

	openNotional := getOpenNotional(stateDB, market, trader)
	// based on last price
	notionalPosition, unrealizedPnl, lastPriceBasedMF := getPositionMetadata(
		lastPrice,
		openNotional,
		size,
		margin,
		blockTimestamp,
	)

	// based on oracle price
	oracleBasedNotional, oracleBasedUnrealizedPnl, oracleBasedMF := getPositionMetadata(
		oraclePrice,
		openNotional,
		size,
		margin,
		blockTimestamp,
	)

	if (marginMode == Maintenance_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == 1) || // for liquidations
		(marginMode == Min_Allowable_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == -1) { // for increasing leverage
		return oracleBasedNotional, oracleBasedUnrealizedPnl
	}
	return notionalPosition, unrealizedPnl
}

func getPositionMetadata(price *big.Int, openNotional *big.Int, size *big.Int, margin *big.Int, blockTimestamp *big.Int) (notionalPos *big.Int, uPnl *big.Int, marginFraction *big.Int) {
	notionalPos = divide1e18(new(big.Int).Mul(price, new(big.Int).Abs(size)))
	if notionalPos.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}
	if size.Sign() == 1 {
		uPnl = new(big.Int).Sub(notionalPos, openNotional)
	} else {
		uPnl = new(big.Int).Sub(openNotional, notionalPos)
	}
	marginFraction = new(big.Int).Div(_multiply1e6(new(big.Int).Add(margin, uPnl), blockTimestamp), notionalPos)
	return notionalPos, uPnl, marginFraction
}

// Common Utils

var _1e18 = big.NewInt(1e18)
var _1e6 = big.NewInt(1e6)

func divide1e18(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, _1e18)
}

func divide1e6(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, _1e6)
}

func _multiply1e6(number *big.Int, blockTimestamp *big.Int) *big.Int {
	if blockTimestamp.Cmp(V2ActivationDate) == 1 {
		return multiply1e6(number)
	}
	return multiply1e6v1(number)
}

// multiple1e6 v1
func multiply1e6v1(number *big.Int) *big.Int {
	return new(big.Int).Div(number, big.NewInt(1e6))

}

func multiply1e6(number *big.Int) *big.Int {
	return new(big.Int).Mul(number, big.NewInt(1e6))
}

func fromTwosComplement(b []byte) *big.Int {
	t := new(big.Int).SetBytes(b)
	if b[0]&0x80 != 0 {
		t.Sub(t, new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8)))
	}
	return t
}
