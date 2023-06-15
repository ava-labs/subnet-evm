package hubblebibliophile

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
func getBlockPlaced(stateDB contract.StateDB, orderHash [32]byte) *big.Int {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	return new(big.Int).SetBytes(stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(orderInfo)).Bytes())
}

func getOrderStatus(stateDB contract.StateDB, orderHash [32]byte) int64 {
	orderInfo := orderInfoMappingStorageSlot(orderHash)
	return new(big.Int).SetBytes(stateDB.GetState(common.HexToAddress(ORDERBOOK_GENESIS_ADDRESS), common.BigToHash(new(big.Int).Add(orderInfo, big.NewInt(3)))).Bytes()).Int64()
}

func orderInfoMappingStorageSlot(orderHash [32]byte) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(append(orderHash[:], common.LeftPadBytes(big.NewInt(ORDER_INFO_SLOT).Bytes(), 32)...)))
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

	market := getMarketAddressFromMarketID(longOrder.AmmIndex.Int64(), stateDB)
	minSize := GetMinSizeRequirement(stateDB, longOrder.AmmIndex.Int64())
	if new(big.Int).Mod(inputStruct.FillAmount, minSize).Cmp(big.NewInt(0)) != 0 {
		return nil, ErrNotMultiple
	}

	oraclePrice := getUnderlyingPrice(stateDB, market)
	spreadLimit := GetMaxOraclePriceSpread(stateDB, longOrder.AmmIndex.Int64())
	blockPlaced0 := getBlockPlaced(stateDB, inputStruct.OrderHashes[0])
	blockPlaced1 := getBlockPlaced(stateDB, inputStruct.OrderHashes[1])

	return determineFillPrice(oraclePrice, spreadLimit, longOrder.Price, shortOrder.Price, blockPlaced0, blockPlaced1)
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

	market := getMarketAddressFromMarketID(order.AmmIndex.Int64(), stateDB)
	minSize := GetMinSizeRequirement(stateDB, order.AmmIndex.Int64())
	if new(big.Int).Mod(inputStruct.FillAmount, minSize).Cmp(big.NewInt(0)) != 0 {
		return nil, ErrNotMultiple
	}

	oraclePrice := getUnderlyingPrice(stateDB, market)
	liquidationSpreadLimit := GetMaxLiquidationPriceSpread(stateDB, order.AmmIndex.Int64())
	liqUpperBound, liqLowerBound := calculateBounds(liquidationSpreadLimit, oraclePrice)

	oracleSpreadLimit := GetMaxOraclePriceSpread(stateDB, order.AmmIndex.Int64())
	upperbound, lowerbound := calculateBounds(oracleSpreadLimit, oraclePrice)
	return determineLiquidationFillPrice(order, liqUpperBound, liqLowerBound, upperbound, lowerbound)
}

func determineLiquidationFillPrice(order IHubbleBibliophileOrder, liqUpperBound, liqLowerBound, upperbound, lowerbound *big.Int) (*big.Int, error) {
	isLongOrder := true
	if order.BaseAssetQuantity.Cmp(big.NewInt(0)) == -1 {
		isLongOrder = false
	}

	if isLongOrder {
		// we are liquidating a long position
		// do not allow liquidation if order.Price < liqLowerBound, because that gives scope for malicious activity to a validator
		if order.Price.Cmp(liqLowerBound) == -1 {
			return nil, ErrTooLow
		}
		return utils.BigIntMin(order.Price, upperbound /* oracle spread upper bound */), nil
	}

	// short order
	if order.Price.Cmp(liqUpperBound) == 1 {
		return nil, ErrTooHigh
	}
	return utils.BigIntMax(order.Price, lowerbound /* oracle spread lower bound */), nil
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
