package hubblebibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	VAR_LAST_PRICE_SLOT             int64 = 1
	VAR_POSITIONS_SLOT              int64 = 2
	VAR_CUMULATIVE_PREMIUM_FRACTION int64 = 3
	MAX_ORACLE_SPREAD_RATIO_SLOT    int64 = 4
	MAX_LIQUIDATION_RATIO_SLOT      int64 = 5
	MIN_SIZE_REQUIREMENT_SLOT       int64 = 6
)

// Reader

// AMM State
func getLastPrice(stateDB contract.StateDB, market common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(big.NewInt(VAR_LAST_PRICE_SLOT))).Big()
}

func getCumulativePremiumFraction(stateDB contract.StateDB, market common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(VAR_CUMULATIVE_PREMIUM_FRACTION))).Bytes())
}

// GetMaxOracleSpreadRatioForMarket returns the maxOracleSpreadRatio for a given market
func GetMaxOracleSpreadRatioForMarket(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_ORACLE_SPREAD_RATIO_SLOT))).Bytes())
}

// GetMaxLiquidationRatioForMarket returns the maxLiquidationRatio for a given market
func GetMaxLiquidationRatioForMarket(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MAX_LIQUIDATION_RATIO_SLOT))).Bytes())
}

// GetMinSizeRequirementForMarket returns the minSizeRequirement for a given market
func GetMinSizeRequirementForMarket(stateDB contract.StateDB, marketID int64) *big.Int {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(big.NewInt(MIN_SIZE_REQUIREMENT_SLOT))).Bytes())
}

// Trader State

func positionsStorageSlot(trader *common.Address) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(VAR_POSITIONS_SLOT).Bytes(), 32)...)))
}

func getSize(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(positionsStorageSlot(trader))).Bytes())
}

func fromTwosComplement(b []byte) *big.Int {
	t := new(big.Int).SetBytes(b)
	if b[0]&0x80 != 0 {
		t.Sub(t, new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8)))
	}
	return t
}

func getOpenNotional(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	return stateDB.GetState(market, common.BigToHash(new(big.Int).Add(positionsStorageSlot(trader), big.NewInt(1)))).Big()
}

func getLastPremiumFraction(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	return fromTwosComplement(stateDB.GetState(market, common.BigToHash(new(big.Int).Add(positionsStorageSlot(trader), big.NewInt(2)))).Bytes())
}

// utilities

func getPendingFundingPayment(stateDB contract.StateDB, market common.Address, trader *common.Address) *big.Int {
	cumulativePremiumFraction := getCumulativePremiumFraction(stateDB, market)
	return divide1e18(new(big.Int).Mul(new(big.Int).Sub(cumulativePremiumFraction, getLastPremiumFraction(stateDB, market, trader)), getSize(stateDB, market, trader)))
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
	notionalPos = divide1e18(new(big.Int).Mul(price, new(big.Int).Abs(size)))
	if notionalPos.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}
	if size.Sign() == 1 {
		uPnl = new(big.Int).Sub(notionalPos, openNotional)
	} else {
		uPnl = new(big.Int).Sub(openNotional, notionalPos)
	}
	marginFraction = new(big.Int).Div(multiply1e6(new(big.Int).Add(margin, uPnl)), notionalPos)
	return notionalPos, uPnl, marginFraction
}

func divide1e18(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, big.NewInt(1e18))
}

func divide1e6(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, big.NewInt(1e6))
}

func multiply1e6(number *big.Int) *big.Int {
	return new(big.Int).Div(number, big.NewInt(1e6))
}

// getMarketAddressFromMarketID returns the market address for a given marketID
func getMarketAddressFromMarketID(marketID int64, stateDB contract.StateDB) common.Address {
	baseStorageSlot := marketsStorageSlot()
	amm := stateDB.GetState(common.HexToAddress(CLEARING_HOUSE_GENESIS_ADDRESS), common.BigToHash(new(big.Int).Add(baseStorageSlot, big.NewInt(marketID))))
	return common.BytesToAddress(amm.Bytes())
}
