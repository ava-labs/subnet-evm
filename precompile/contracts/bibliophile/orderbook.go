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
	ORDER_INFO_SLOT           int64 = 53
	IS_VALIDATOR_SLOT         int64 = 54
	REDUCE_ONLY_AMOUNT_SLOT   int64 = 55
	IS_TRADING_AUTHORITY_SLOT int64 = 61
	LONG_OPEN_ORDERS_SLOT     int64 = 65
	SHORT_OPEN_ORDERS_SLOT    int64 = 66
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
func getReduceOnlyAmount(stateDB contract.StateDB, trader common.Address, ammIndex *big.Int) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(REDUCE_ONLY_AMOUNT_SLOT).Bytes(), 32)...))
	nestedMappingHash := crypto.Keccak256(append(common.LeftPadBytes(ammIndex.Bytes(), 32), baseMappingHash...))
	return fromTwosComplement(stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(nestedMappingHash)).Bytes())
}

func getLongOpenOrdersAmount(stateDB contract.StateDB, trader common.Address, ammIndex *big.Int) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(LONG_OPEN_ORDERS_SLOT).Bytes(), 32)...))
	nestedMappingHash := crypto.Keccak256(append(common.LeftPadBytes(ammIndex.Bytes(), 32), baseMappingHash...))
	return stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(nestedMappingHash)).Big()
}

func getShortOpenOrdersAmount(stateDB contract.StateDB, trader common.Address, ammIndex *big.Int) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(SHORT_OPEN_ORDERS_SLOT).Bytes(), 32)...))
	nestedMappingHash := crypto.Keccak256(append(common.LeftPadBytes(ammIndex.Bytes(), 32), baseMappingHash...))
	return stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BytesToHash(nestedMappingHash)).Big()
}

func getBlockPlaced(stateDB contract.StateDB, orderHash [32]byte) *big.Int {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	return new(big.Int).SetBytes(stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(orderInfo)).Bytes())
}

func getOrderFilledAmount(stateDB contract.StateDB, orderHash [32]byte) *big.Int {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	num := stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(new(big.Int).Add(orderInfo, big.NewInt(1)))).Bytes()
	return fromTwosComplement(num)
}

func getOrderStatus(stateDB contract.StateDB, orderHash [32]byte) int64 {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	return new(big.Int).SetBytes(stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(new(big.Int).Add(orderInfo, big.NewInt(3)))).Bytes()).Int64()
}

func orderInfoMappingStorageSlot(orderHash [32]byte) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(append(orderHash[:], common.LeftPadBytes(big.NewInt(ORDER_INFO_SLOT).Bytes(), 32)...)))
}

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
	spreadLimit := GetMaxOraclePriceSpread(stateDB, marketID)
	oraclePrice := getUnderlyingPriceForMarket(stateDB, marketID)
	return calculateBounds(spreadLimit, oraclePrice)
}

func GetAcceptableBoundsForLiquidation(stateDB contract.StateDB, marketID int64) (upperBound, lowerBound *big.Int) {
	spreadLimit := GetMaxLiquidationPriceSpread(stateDB, marketID)
	oraclePrice := getUnderlyingPriceForMarket(stateDB, marketID)
	return calculateBounds(spreadLimit, oraclePrice)
}

func calculateBounds(spreadLimit, oraclePrice *big.Int) (*big.Int, *big.Int) {
	upperbound := hu.Div1e6(hu.Mul(oraclePrice, hu.Add(hu.ONE_E_6, spreadLimit)))
	lowerbound := big.NewInt(0)
	if spreadLimit.Cmp(hu.ONE_E_6) == -1 {
		lowerbound = hu.Div1e6(hu.Mul(oraclePrice, hu.Sub(hu.ONE_E_6, spreadLimit)))
	}
	return upperbound, lowerbound
}
