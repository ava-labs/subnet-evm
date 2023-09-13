package bibliophile

import (
	"errors"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/utils"

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

// Business Logic

func ValidateOrdersAndDetermineFillPrice(stateDB contract.StateDB, inputStruct *ValidateOrdersAndDetermineFillPriceInput) (*ValidateOrdersAndDetermineFillPriceOutput, error) {
	longOrder := inputStruct.Orders[0]
	shortOrder := inputStruct.Orders[1]

	if longOrder.BaseAssetQuantity.Cmp(big.NewInt(0)) <= 0 {
		return nil, ErrNotLongOrder
	}

	if shortOrder.BaseAssetQuantity.Cmp(big.NewInt(0)) >= 0 {
		return nil, ErrNotShortOrder
	}

	if longOrder.AmmIndex.Cmp(shortOrder.AmmIndex) != 0 {
		return nil, ErrNotSameAMM
	}

	if longOrder.Price.Cmp(shortOrder.Price) == -1 {
		return nil, ErrNoMatch
	}

	if getOrderStatus(stateDB, inputStruct.OrderHashes[0]) != 1 || getOrderStatus(stateDB, inputStruct.OrderHashes[1]) != 1 {
		return nil, ErrInvalidOrder
	}

	blockPlaced0 := getBlockPlaced(stateDB, inputStruct.OrderHashes[0])
	blockPlaced1 := getBlockPlaced(stateDB, inputStruct.OrderHashes[1])
	minSize := GetMinSizeRequirement(stateDB, longOrder.AmmIndex.Int64())
	if new(big.Int).Mod(inputStruct.FillAmount, minSize).Cmp(big.NewInt(0)) != 0 {
		return nil, ErrNotMultiple
	}
	return DetermineFillPrice(stateDB, longOrder.AmmIndex.Int64(), longOrder.Price, shortOrder.Price, blockPlaced0, blockPlaced1)
}

func DetermineFillPrice(stateDB contract.StateDB, marketId int64, longOrderPrice, shortOrderPrice, blockPlaced0, blockPlaced1 *big.Int) (*ValidateOrdersAndDetermineFillPriceOutput, error) {
	market := getMarketAddressFromMarketID(marketId, stateDB)
	oraclePrice := getUnderlyingPrice(stateDB, market)
	spreadLimit := GetMaxOraclePriceSpread(stateDB, marketId)
	return determineFillPrice(oraclePrice, spreadLimit, longOrderPrice, shortOrderPrice, blockPlaced0, blockPlaced1)
}

func determineFillPrice(oraclePrice, spreadLimit, longOrderPrice, shortOrderPrice, blockPlaced0, blockPlaced1 *big.Int) (*ValidateOrdersAndDetermineFillPriceOutput, error) {
	upperbound, lowerbound := calculateBounds(spreadLimit, oraclePrice)
	if longOrderPrice.Cmp(lowerbound) == -1 {
		return nil, ErrTooLow
	}
	if shortOrderPrice.Cmp(upperbound) == 1 {
		return nil, ErrTooHigh
	}

	output := ValidateOrdersAndDetermineFillPriceOutput{}
	if blockPlaced0.Cmp(blockPlaced1) == -1 {
		// long order is the maker order
		output.FillPrice = utils.BigIntMin(longOrderPrice, upperbound)
		output.Mode0 = 1 // Mode0 corresponds to the long order and `1` is maker
		output.Mode1 = 0 // Mode1 corresponds to the short order and `0` is taker
	} else { // if long order is placed after short order or in the same block as short
		// short order is the maker order
		output.FillPrice = utils.BigIntMax(shortOrderPrice, lowerbound)
		output.Mode0 = 0 // Mode0 corresponds to the long order and `0` is taker
		output.Mode1 = 1 // Mode1 corresponds to the short order and `1` is maker
	}
	return &output, nil
}

func ValidateLiquidationOrderAndDetermineFillPrice(stateDB contract.StateDB, inputStruct *ValidateLiquidationOrderAndDetermineFillPriceInput) (*big.Int, error) {
	order := inputStruct.Order
	minSize := GetMinSizeRequirement(stateDB, order.AmmIndex.Int64())
	if new(big.Int).Mod(inputStruct.FillAmount, minSize).Cmp(big.NewInt(0)) != 0 {
		return nil, ErrNotMultiple
	}
	return DetermineLiquidationFillPrice(stateDB, order.AmmIndex.Int64(), order.BaseAssetQuantity, order.Price)
}

func DetermineLiquidationFillPrice(stateDB contract.StateDB, marketId int64, baseAssetQuantity, price *big.Int) (*big.Int, error) {
	isLongOrder := true
	if baseAssetQuantity.Sign() < 0 {
		isLongOrder = false
	}
	market := getMarketAddressFromMarketID(marketId, stateDB)
	oraclePrice := getUnderlyingPrice(stateDB, market)
	liquidationSpreadLimit := GetMaxLiquidationPriceSpread(stateDB, marketId)
	liqUpperBound, liqLowerBound := calculateBounds(liquidationSpreadLimit, oraclePrice)

	oracleSpreadLimit := GetMaxOraclePriceSpread(stateDB, marketId)
	upperbound, lowerbound := calculateBounds(oracleSpreadLimit, oraclePrice)
	return determineLiquidationFillPrice(isLongOrder, price, liqUpperBound, liqLowerBound, upperbound, lowerbound)
}

func determineLiquidationFillPrice(isLongOrder bool, price, liqUpperBound, liqLowerBound, upperbound, lowerbound *big.Int) (*big.Int, error) {
	if isLongOrder {
		// we are liquidating a long position
		// do not allow liquidation if order.Price < liqLowerBound, because that gives scope for malicious activity to a validator
		if price.Cmp(liqLowerBound) == -1 {
			return nil, ErrTooLow
		}
		return utils.BigIntMin(price, upperbound /* oracle spread upper bound */), nil
	}

	// short order
	if price.Cmp(liqUpperBound) == 1 {
		return nil, ErrTooHigh
	}
	return utils.BigIntMax(price, lowerbound /* oracle spread lower bound */), nil
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
	upperbound := divide1e6(new(big.Int).Mul(oraclePrice, new(big.Int).Add(_1e6, spreadLimit)))
	lowerbound := big.NewInt(0)
	if spreadLimit.Cmp(_1e6) == -1 {
		lowerbound = divide1e6(new(big.Int).Mul(oraclePrice, new(big.Int).Sub(_1e6, spreadLimit)))
	}
	return upperbound, lowerbound
}
