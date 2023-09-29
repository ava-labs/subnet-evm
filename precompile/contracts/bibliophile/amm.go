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
	UNDERLYING_ASSET_SLOT           int64 = 6
	MAX_LIQUIDATION_PRICE_SPREAD    int64 = 11
	MULTIPLIER_SLOT                 int64 = 12
	IMPACT_MARGIN_NOTIONAL_SLOT     int64 = 19
	LAST_TRADE_PRICE_SLOT           int64 = 20
	BIDS_SLOT                       int64 = 21
	ASKS_SLOT                       int64 = 22
	BIDS_HEAD_SLOT                  int64 = 23
	ASKS_HEAD_SLOT                  int64 = 24
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

func getUnderlyingAssetAddress(stateDB contract.StateDB, market common.Address) common.Address {
	return common.BytesToAddress(stateDB.GetState(market, common.BigToHash(big.NewInt(UNDERLYING_ASSET_SLOT))).Bytes())
}

func getUnderlyingPriceForMarket(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return getUnderlyingPrice(stateDB, market)
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

func getPosition(stateDB contract.StateDB, market common.Address, trader *common.Address) *hu.Position {
	return &hu.Position{
		Size:         getSize(stateDB, market, trader),
		OpenNotional: getOpenNotional(stateDB, market, trader),
	}
}

// Utils

func getPendingFundingPayment(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	cumulativePremiumFraction := GetCumulativePremiumFraction(stateDB, market)
	return hu.Div1e18(new(big.Int).Mul(new(big.Int).Sub(cumulativePremiumFraction, GetLastPremiumFraction(stateDB, market, trader)), getSize(stateDB, market, trader)))
}

// Common Utils
func fromTwosComplement(b []byte) *big.Int {
	t := new(big.Int).SetBytes(b)
	if b[0]&0x80 != 0 {
		t.Sub(t, new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8)))
	}
	return t
}
