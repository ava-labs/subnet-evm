package bibliophile

import (
	"math/big"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	VAR_POSITIONS_SLOT              int64 = 1
	VAR_CUMULATIVE_PREMIUM_FRACTION int64 = 2
	MAX_ORACLE_SPREAD_RATIO_SLOT    int64 = 3
	MAX_LIQUIDATION_RATIO_SLOT      int64 = 4
	MIN_SIZE_REQUIREMENT_SLOT       int64 = 5
	ORACLE_SLOT                     int64 = 6
	UNDERLYING_ASSET_SLOT           int64 = 7
	MAX_LIQUIDATION_PRICE_SPREAD    int64 = 12
	MULTIPLIER_SLOT                 int64 = 13
	RED_STONE_ADAPTER_SLOT          int64 = 16
	RED_STONE_FEED_ID_SLOT          int64 = 17
	IMPACT_MARGIN_NOTIONAL_SLOT     int64 = 22
	LAST_TRADE_PRICE_SLOT           int64 = 23
	BIDS_SLOT                       int64 = 24
	ASKS_SLOT                       int64 = 25
	BIDS_HEAD_SLOT                  int64 = 26
	ASKS_HEAD_SLOT                  int64 = 27
)

const (
	// this slot is from TestOracle.sol
	TEST_ORACLE_PRICES_MAPPING_SLOT int64 = 3
)

// AMM State
func getBidsHead(stateDB contract.StateDB, market common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(BIDS_HEAD_SLOT))).Big()
}

func getAsksHead(stateDB contract.StateDB, market common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(ASKS_HEAD_SLOT))).Big()
}

func getLastPrice(stateDB contract.StateDB, market common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(LAST_TRADE_PRICE_SLOT))).Big()
}

func GetCumulativePremiumFraction(stateDB contract.StateDB, market common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(VAR_CUMULATIVE_PREMIUM_FRACTION))).Bytes())
}

// GetMaxOraclePriceSpread returns the maxOracleSpreadRatio for a given market
func GetMaxOraclePriceSpread(stateDB contract.StateDB, marketID int64) *big.Int {
	return getMaxOraclePriceSpread(stateDB, getMarketAddressFromMarketID(marketID, stateDB))
}

func getMaxOraclePriceSpread(stateDB contract.StateDB, market common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_ORACLE_SPREAD_RATIO_SLOT))).Bytes())
}

// GetMaxLiquidationPriceSpread returns the maxOracleSpreadRatio for a given market
func GetMaxLiquidationPriceSpread(stateDB contract.StateDB, marketID int64) *big.Int {
	return getMaxLiquidationPriceSpread(stateDB, getMarketAddressFromMarketID(marketID, stateDB))
}

func getMaxLiquidationPriceSpread(stateDB contract.StateDB, market common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_LIQUIDATION_PRICE_SPREAD))).Bytes())
}

// GetMaxLiquidationRatio returns the maxLiquidationPriceSpread for a given market
func GetMaxLiquidationRatio(stateDB contract.StateDB, marketID int64) *big.Int {
	return getMaxLiquidationRatio(stateDB, getMarketAddressFromMarketID(marketID, stateDB))
}

func getMaxLiquidationRatio(stateDB contract.StateDB, market common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_LIQUIDATION_RATIO_SLOT))).Bytes())
}

// GetMinSizeRequirement returns the minSizeRequirement for a given market
func GetMinSizeRequirement(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MIN_SIZE_REQUIREMENT_SLOT))).Bytes())
}

func getMultiplier(stateDB contract.StateDB, market common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(MULTIPLIER_SLOT))).Big()
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
			// log.Info("redstone-price", "amm", market, "price", redstonePrice)
			return redstonePrice
		}
		/* else {
			// just log the red stone price, for testing before deployment
			var feedId common.Hash
			if strings.EqualFold(market.String(), "0xa72b463C21dA61cCc86069cFab82e9e8491152a0") { // eth amm
				feedId = common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000")
			} else if strings.EqualFold(market.String(), "0xd80e57dB448b0692C396B890eE9c791D7386dAdC") { // avax amm
				feedId = common.HexToHash("0x4156415800000000000000000000000000000000000000000000000000000000")
			}
			// redstonePrice := getRedStonePrice(stateDB, redStoneAdapter, feedId)
			// log.Info("log-only-redstone-price", "amm", market, "price", redstonePrice)
		} */
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

func bidsStorageSlot(price *big.Int) common.Hash {
	return common.BytesToHash(crypto.Keccak256(append(common.LeftPadBytes(price.Bytes(), 32), common.LeftPadBytes(big.NewInt(BIDS_SLOT).Bytes(), 32)...)))
}

func asksStorageSlot(price *big.Int) common.Hash {
	return common.BytesToHash(crypto.Keccak256(append(common.LeftPadBytes(price.Bytes(), 32), common.LeftPadBytes(big.NewInt(ASKS_SLOT).Bytes(), 32)...)))
}

func getBidSize(stateDB contract.StateDB, market common.Address, price *big.Int) *big.Int {
	return stateDB.GetState(market, common.BigToHash(new(big.Int).Add(bidsStorageSlot(price).Big(), big.NewInt(1)))).Big()
}

func getAskSize(stateDB contract.StateDB, market common.Address, price *big.Int) *big.Int {
	return stateDB.GetState(market, common.BigToHash(new(big.Int).Add(asksStorageSlot(price).Big(), big.NewInt(1)))).Big()
}

func getNextBid(stateDB contract.StateDB, market common.Address, price *big.Int) *big.Int {
	return stateDB.GetState(market, bidsStorageSlot(price)).Big()
}

func getNextAsk(stateDB contract.StateDB, market common.Address, price *big.Int) *big.Int {
	return stateDB.GetState(market, asksStorageSlot(price)).Big()
}

func getImpactMarginNotional(stateDB contract.StateDB, market common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(IMPACT_MARGIN_NOTIONAL_SLOT))).Big()
}

// Utils

func getPendingFundingPayment(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	cumulativePremiumFraction := GetCumulativePremiumFraction(stateDB, market)
	return hu.Div1e18(new(big.Int).Mul(new(big.Int).Sub(cumulativePremiumFraction, GetLastPremiumFraction(stateDB, market, trader)), getSize(stateDB, market, trader)))
}

func getOptimalPnl(stateDB contract.StateDB, market common.Address, oraclePrice *big.Int, lastPrice *big.Int, trader *common.Address, margin *big.Int, marginMode MarginMode) (notionalPosition *big.Int, uPnL *big.Int) {
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
	)

	// based on oracle price
	oracleBasedNotional, oracleBasedUnrealizedPnl, oracleBasedMF := getPositionMetadata(
		oraclePrice,
		openNotional,
		size,
		margin,
	)

	if (marginMode == Maintenance_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == 1) || // for liquidations
		(marginMode == Min_Allowable_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == -1) { // for increasing leverage
		return oracleBasedNotional, oracleBasedUnrealizedPnl
	}
	return notionalPosition, unrealizedPnl
}

func getPositionMetadata(price *big.Int, openNotional *big.Int, size *big.Int, margin *big.Int) (notionalPos *big.Int, uPnl *big.Int, marginFraction *big.Int) {
	notionalPos = hu.Div1e18(new(big.Int).Mul(price, new(big.Int).Abs(size)))
	if notionalPos.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}
	if size.Sign() == 1 {
		uPnl = new(big.Int).Sub(notionalPos, openNotional)
	} else {
		uPnl = new(big.Int).Sub(openNotional, notionalPos)
	}
	marginFraction = new(big.Int).Div(hu.Mul1e6(new(big.Int).Add(margin, uPnl)), notionalPos)
	return notionalPos, uPnl, marginFraction
}

// Common Utils
func fromTwosComplement(b []byte) *big.Int {
	t := new(big.Int).SetBytes(b)
	if b[0]&0x80 != 0 {
		t.Sub(t, new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8)))
	}
	return t
}
