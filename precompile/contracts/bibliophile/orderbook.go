package bibliophile

import (
	"errors"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	ORDERBOOK_GENESIS_ADDRESS       = "0x0300000000000000000000000000000000000000"
	IS_VALIDATOR_SLOT         int64 = 1
	IS_TRADING_AUTHORITY_SLOT int64 = 2
)

var (
	ErrNotLongOrder  = errors.New("OB_order_0_is_not_long")
	ErrNotShortOrder = errors.New("OB_order_1_is_not_short")
	ErrNotSameAMM    = errors.New("OB_orders_for_different_amms")
	ErrNoMatch       = errors.New("OB_orders_do_not_match")
	ErrInvalidOrder  = errors.New("OB_invalid_order")
	ErrNotMultiple   = errors.New("OB.not_multiple")
	ErrTooLow        = errors.New("OB_long_order_price_too_low")
	ErrTooHigh       = errors.New("OB_short_order_price_too_high")
)

// State Reader

func IsTradingAuthority(stateDB contract.StateDB, trader, senderOrSigner common.Address) bool {
	tradingAuthorityMappingSlot := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(IS_TRADING_AUTHORITY_SLOT).Bytes(), 32)...))
	tradingAuthorityMappingSlot = crypto.Keccak256(append(common.LeftPadBytes(senderOrSigner.Bytes(), 32), tradingAuthorityMappingSlot...))
	return stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(tradingAuthorityMappingSlot)).Big().Cmp(big.NewInt(1)) == 0
}

func IsValidator(stateDB contract.StateDB, senderOrSigner common.Address) bool {
	isValidatorMappingSlot := crypto.Keccak256(append(common.LeftPadBytes(senderOrSigner.Bytes(), 32), common.LeftPadBytes(big.NewInt(IS_VALIDATOR_SLOT).Bytes(), 32)...))
	return stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(isValidatorMappingSlot)).Big().Cmp(big.NewInt(1)) == 0
}

// Helper functions

func GetAcceptableBounds(stateDB contract.StateDB, marketID int64) (upperBound, lowerBound *big.Int) {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return calculateBounds(getMaxOraclePriceSpread(stateDB, market), getUnderlyingPrice(stateDB, market), getMultiplier(stateDB, market))
}

func GetAcceptableBoundsForLiquidation(stateDB contract.StateDB, marketID int64) (upperBound, lowerBound *big.Int) {
	market := getMarketAddressFromMarketID(marketID, stateDB)
	return calculateBounds(getMaxLiquidationPriceSpread(stateDB, market), getUnderlyingPrice(stateDB, market), getMultiplier(stateDB, market))
}

func calculateBounds(spreadLimit, oraclePrice, multiplier *big.Int) (*big.Int, *big.Int) {
	upperbound := hu.RoundOff(hu.Div1e6(hu.Mul(oraclePrice, hu.Add1e6(spreadLimit))), multiplier)
	lowerbound := big.NewInt(0)
	if spreadLimit.Cmp(hu.ONE_E_6) == -1 {
		lowerbound = hu.RoundOff(hu.Div1e6(hu.Mul(oraclePrice, hu.Sub(hu.ONE_E_6, spreadLimit))), multiplier)
	}
	return upperbound, lowerbound
}
