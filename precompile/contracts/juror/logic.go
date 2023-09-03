package juror

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type OrderType uint8

// has to be exact same as expected in contracts
const (
	Limit OrderType = iota
	IOC
	LimitV2
)

type DecodeStep struct {
	OrderType    OrderType
	EncodedOrder []byte
}

type Metadata struct {
	AmmIndex          *big.Int
	Trader            common.Address
	BaseAssetQuantity *big.Int
	Price             *big.Int
	BlockPlaced       *big.Int
	OrderHash         common.Hash
}

type Side uint8

const (
	Long Side = iota
	Short
	Liquidation
)

type OrderStatus uint8

// has to be exact same as IOrderHandler
const (
	Invalid OrderStatus = iota
	Placed
	Filled
	Cancelled
)

var (
	ErrTwoOrders         = errors.New("need 2 orders")
	ErrInvalidFillAmount = errors.New("invalid fillAmount")
	ErrNotLongOrder      = errors.New("not long")
	ErrNotShortOrder     = errors.New("not short")
	ErrNotSameAMM        = errors.New("OB_orders_for_different_amms")
	ErrNoMatch           = errors.New("OB_orders_do_not_match")
	ErrNotMultiple       = errors.New("not multiple")

	ErrInvalidOrder                       = errors.New("invalid order")
	ErrInvalidPrice                       = errors.New("invalid price")
	ErrCancelledOrder                     = errors.New("cancelled order")
	ErrFilledOrder                        = errors.New("filled order")
	ErrOrderAlreadyExists                 = errors.New("order already exists")
	ErrTooLow                             = errors.New("OB_long_order_price_too_low")
	ErrTooHigh                            = errors.New("OB_short_order_price_too_high")
	ErrOverFill                           = errors.New("overfill")
	ErrReduceOnlyAmountExceeded           = errors.New("not reducing pos")
	ErrBaseAssetQuantityZero              = errors.New("baseAssetQuantity is zero")
	ErrReduceOnlyBaseAssetQuantityInvalid = errors.New("reduce only order must reduce position")
	ErrNetReduceOnlyAmountExceeded        = errors.New("net reduce only amount exceeded")
	ErrStaleReduceOnlyOrders              = errors.New("cancel stale reduce only orders")
	ErrInsufficientMargin                 = errors.New("insufficient margin")
	ErrCrossingMarket                     = errors.New("crossing market")
	ErrOpenOrders                         = errors.New("open orders")
	ErrOpenReduceOnlyOrders               = errors.New("open reduce only orders")
	ErrNoTradingAuthority                 = errors.New("no trading authority")
)

// Business Logic
func ValidateOrdersAndDetermineFillPrice(bibliophile b.BibliophileClient, inputStruct *ValidateOrdersAndDetermineFillPriceInput) (*ValidateOrdersAndDetermineFillPriceOutput, error) {
	if len(inputStruct.Data) != 2 {
		return nil, ErrTwoOrders
	}

	if inputStruct.FillAmount.Sign() <= 0 {
		return nil, ErrInvalidFillAmount
	}

	decodeStep0, err := decodeTypeAndEncodedOrder(inputStruct.Data[0])
	log.Info("decodeStep0", "decodeStep0", decodeStep0, "err", err)
	if err != nil {
		return nil, err
	}
	m0, err := validateOrder(bibliophile, decodeStep0.OrderType, decodeStep0.EncodedOrder, Long, inputStruct.FillAmount)
	if err != nil {
		return nil, err
	}

	decodeStep1, err := decodeTypeAndEncodedOrder(inputStruct.Data[1])
	log.Info("decodeStep1", "decodeStep1", decodeStep1, "err", err)
	if err != nil {
		return nil, err
	}
	m1, err := validateOrder(bibliophile, decodeStep1.OrderType, decodeStep1.EncodedOrder, Short, new(big.Int).Neg(inputStruct.FillAmount))
	if err != nil {
		return nil, err
	}

	if m0.AmmIndex.Cmp(m1.AmmIndex) != 0 {
		return nil, ErrNotSameAMM
	}

	if m0.Price.Cmp(m1.Price) < 0 {
		return nil, ErrNoMatch
	}

	minSize := bibliophile.GetMinSizeRequirement(m0.AmmIndex.Int64())
	if new(big.Int).Mod(inputStruct.FillAmount, minSize).Cmp(big.NewInt(0)) != 0 {
		return nil, ErrNotMultiple
	}

	fillPriceAndModes, err := bibliophile.DetermineFillPrice(m0.AmmIndex.Int64(), m0.Price, m1.Price, m0.BlockPlaced, m1.BlockPlaced)
	if err != nil {
		return nil, err
	}

	output := &ValidateOrdersAndDetermineFillPriceOutput{
		Instructions: [2]IClearingHouseInstruction{
			IClearingHouseInstruction{
				AmmIndex:  m0.AmmIndex,
				Trader:    m0.Trader,
				OrderHash: m0.OrderHash,
				Mode:      fillPriceAndModes.Mode0,
			},
			IClearingHouseInstruction{
				AmmIndex:  m1.AmmIndex,
				Trader:    m1.Trader,
				OrderHash: m1.OrderHash,
				Mode:      fillPriceAndModes.Mode1,
			},
		},
		OrderTypes: [2]uint8{uint8(decodeStep0.OrderType), uint8(decodeStep1.OrderType)},
		EncodedOrders: [2][]byte{
			decodeStep0.EncodedOrder,
			decodeStep1.EncodedOrder,
		},
		FillPrice: fillPriceAndModes.FillPrice,
	}
	return output, nil
}

func ValidateLiquidationOrderAndDetermineFillPrice(bibliophile b.BibliophileClient, inputStruct *ValidateLiquidationOrderAndDetermineFillPriceInput) (*ValidateLiquidationOrderAndDetermineFillPriceOutput, error) {
	fillAmount := new(big.Int).Set(inputStruct.LiquidationAmount)
	if fillAmount.Sign() <= 0 {
		return nil, ErrInvalidFillAmount
	}

	decodeStep0, err := decodeTypeAndEncodedOrder(inputStruct.Data)
	if err != nil {
		return nil, err
	}
	m0, err := validateOrder(bibliophile, decodeStep0.OrderType, decodeStep0.EncodedOrder, Liquidation, fillAmount)
	if err != nil {
		return nil, err
	}

	if m0.BaseAssetQuantity.Sign() < 0 {
		fillAmount = new(big.Int).Neg(fillAmount)
	}

	minSize := bibliophile.GetMinSizeRequirement(m0.AmmIndex.Int64())
	if new(big.Int).Mod(fillAmount, minSize).Cmp(big.NewInt(0)) != 0 {
		return nil, ErrNotMultiple
	}

	fillPrice, err := bibliophile.DetermineLiquidationFillPrice(m0.AmmIndex.Int64(), m0.BaseAssetQuantity, m0.Price)
	if err != nil {
		return nil, err
	}

	output := &ValidateLiquidationOrderAndDetermineFillPriceOutput{
		Instruction: IClearingHouseInstruction{
			AmmIndex:  m0.AmmIndex,
			Trader:    m0.Trader,
			OrderHash: m0.OrderHash,
			Mode:      1, // Maker
		},
		OrderType:    uint8(decodeStep0.OrderType),
		EncodedOrder: decodeStep0.EncodedOrder,
		FillPrice:    fillPrice,
		FillAmount:   fillAmount,
	}
	return output, nil
}

func decodeTypeAndEncodedOrder(data []byte) (*DecodeStep, error) {
	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	decodedValues, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return &DecodeStep{
		OrderType:    OrderType(decodedValues[0].(uint8)),
		EncodedOrder: decodedValues[1].([]byte),
	}, nil
}

func validateOrder(bibliophile b.BibliophileClient, orderType OrderType, encodedOrder []byte, side Side, fillAmount *big.Int) (metadata *Metadata, err error) {
	if orderType == Limit {
		order, err := orderbook.DecodeLimitOrder(encodedOrder)
		if err != nil {
			return nil, err
		}
		orderHash, err := GetLimitOrderHash(order)
		if err != nil {
			return nil, err
		}
		return validateExecuteLimitOrder(bibliophile, order, side, fillAmount, orderHash)
	}
	if orderType == IOC {
		order, err := orderbook.DecodeIOCOrder(encodedOrder)
		if err != nil {
			return nil, err
		}
		return validateExecuteIOCOrder(bibliophile, order, side, fillAmount)
	}
	if orderType == LimitV2 {
		order, err := orderbook.DecodeLimitOrderV2(encodedOrder)
		if err != nil {
			return nil, err
		}
		orderHash, err := GetLimitOrderV2Hash_2(order)
		if err != nil {
			return nil, err
		}
		// order.postOnly field is not required to be validated while matching
		return validateExecuteLimitOrder(bibliophile, &order.LimitOrder, side, fillAmount, orderHash)
	}
	return nil, errors.New("invalid order type")
}

// Limit Orders

func validateExecuteLimitOrder(bibliophile b.BibliophileClient, order *orderbook.LimitOrder, side Side, fillAmount *big.Int, orderHash [32]byte) (metadata *Metadata, err error) {
	if err := validateLimitOrderLike(bibliophile, order, bibliophile.GetOrderFilledAmount(orderHash), OrderStatus(bibliophile.GetOrderStatus(orderHash)), side, fillAmount); err != nil {
		return nil, err
	}
	return &Metadata{
		AmmIndex:          order.AmmIndex,
		Trader:            order.Trader,
		BaseAssetQuantity: order.BaseAssetQuantity,
		BlockPlaced:       bibliophile.GetBlockPlaced(orderHash),
		Price:             order.Price,
		OrderHash:         orderHash,
	}, nil
}

func validateLimitOrderLike(bibliophile b.BibliophileClient, order *orderbook.LimitOrder, filledAmount *big.Int, status OrderStatus, side Side, fillAmount *big.Int) error {
	if status != Placed {
		return ErrInvalidOrder
	}

	// in case of liquidations, side of the order is determined by the sign of the base asset quantity, so basically base asset quantity check is redundant
	if side == Liquidation {
		if order.BaseAssetQuantity.Sign() > 0 {
			side = Long
		} else if order.BaseAssetQuantity.Sign() < 0 {
			side = Short
			fillAmount = new(big.Int).Neg(fillAmount)
		}
	}

	market := bibliophile.GetMarketAddressFromMarketID(order.AmmIndex.Int64())
	if side == Long {
		if order.BaseAssetQuantity.Sign() <= 0 {
			return ErrNotLongOrder
		}
		if fillAmount.Sign() <= 0 {
			return ErrInvalidFillAmount
		}
		if new(big.Int).Add(filledAmount, fillAmount).Cmp(order.BaseAssetQuantity) > 0 {
			return ErrOverFill
		}
		if order.ReduceOnly {
			posSize := bibliophile.GetSize(market, &order.Trader)
			// posSize should be closed to continue to be Short
			// this also returns err if posSize >= 0, which should not happen because we are executing a long reduceOnly order on this account
			if new(big.Int).Add(posSize, fillAmount).Sign() > 0 {
				return ErrReduceOnlyAmountExceeded
			}
		}
	} else if side == Short {
		if order.BaseAssetQuantity.Sign() >= 0 {
			return ErrNotShortOrder
		}
		if fillAmount.Sign() >= 0 {
			return ErrInvalidFillAmount
		}
		if new(big.Int).Add(filledAmount, fillAmount).Cmp(order.BaseAssetQuantity) < 0 { // all quantities are -ve
			return ErrOverFill
		}
		if order.ReduceOnly {
			posSize := bibliophile.GetSize(market, &order.Trader)
			// posSize should continue to be Long
			// this also returns is posSize <= 0, which should not happen because we are executing a short reduceOnly order on this account
			if new(big.Int).Add(posSize, fillAmount).Sign() < 0 {
				return ErrReduceOnlyAmountExceeded
			}
		}
	} else {
		return errors.New("invalid side")
	}
	return nil
}

// IOC Orders
func ValidatePlaceIOCOrders(bibliophile b.BibliophileClient, inputStruct *ValidatePlaceIOCOrdersInput) (orderHashes [][32]byte, err error) {
	orders := inputStruct.Orders
	if len(orders) == 0 {
		return nil, errors.New("no orders")
	}
	trader := orders[0].Trader
	if !strings.EqualFold(trader.String(), inputStruct.Sender.String()) && !bibliophile.IsTradingAuthority(trader, inputStruct.Sender) {
		return nil, errors.New("no trading authority")
	}
	blockTimestamp := big.NewInt(int64(bibliophile.GetAccessibleState().GetBlockContext().Timestamp()))
	expireWithin := new(big.Int).Add(blockTimestamp, bibliophile.IOC_GetExpirationCap())
	orderHashes = make([][32]byte, len(orders))
	for i, order := range orders {
		if order.BaseAssetQuantity.Sign() == 0 {
			return nil, ErrInvalidFillAmount
		}
		if !strings.EqualFold(order.Trader.String(), trader.String()) {
			return nil, errors.New("OB_trader_mismatch")
		}
		if OrderType(order.OrderType) != IOC {
			return nil, errors.New("not_ioc_order")
		}
		if order.ExpireAt.Cmp(blockTimestamp) < 0 {
			return nil, errors.New("ioc expired")
		}
		if order.ExpireAt.Cmp(expireWithin) > 0 {
			return nil, errors.New("ioc expiration too far")
		}
		minSize := bibliophile.GetMinSizeRequirement(order.AmmIndex.Int64())
		if new(big.Int).Mod(order.BaseAssetQuantity, minSize).Sign() != 0 {
			return nil, ErrNotMultiple
		}
		// this check is as such not required, because even if this order is not reducing the position, it will be rejected by the matching engine and expire away
		// this check is sort of also redundant because either ways user can circumvent this by placing several reduceOnly orders
		// if order.ReduceOnly {}
		orderHashes[i], err = getIOCOrderHash(&orderbook.IOCOrder{
			OrderType: order.OrderType,
			ExpireAt:  order.ExpireAt,
			LimitOrder: orderbook.LimitOrder{
				AmmIndex:          order.AmmIndex,
				Trader:            order.Trader,
				BaseAssetQuantity: order.BaseAssetQuantity,
				Price:             order.Price,
				Salt:              order.Salt,
				ReduceOnly:        order.ReduceOnly,
			},
		})
		if err != nil {
			return
		}
		if OrderStatus(bibliophile.IOC_GetOrderStatus(orderHashes[i])) != Invalid {
			return nil, ErrInvalidOrder
		}
	}
	return
}

func validateExecuteIOCOrder(bibliophile b.BibliophileClient, order *orderbook.IOCOrder, side Side, fillAmount *big.Int) (metadata *Metadata, err error) {
	if OrderType(order.OrderType) != IOC {
		return nil, errors.New("not ioc order")
	}
	blockTimestamp := big.NewInt(int64(bibliophile.GetAccessibleState().GetBlockContext().Timestamp()))
	if order.ExpireAt.Cmp(blockTimestamp) < 0 {
		return nil, errors.New("ioc expired")
	}
	orderHash, err := getIOCOrderHash(order)
	if err != nil {
		return nil, err
	}
	if err := validateLimitOrderLike(bibliophile, &order.LimitOrder, bibliophile.IOC_GetOrderFilledAmount(orderHash), OrderStatus(bibliophile.IOC_GetOrderStatus(orderHash)), side, fillAmount); err != nil {
		return nil, err
	}
	return &Metadata{
		AmmIndex:          order.AmmIndex,
		Trader:            order.Trader,
		BaseAssetQuantity: order.BaseAssetQuantity,
		BlockPlaced:       bibliophile.IOC_GetBlockPlaced(orderHash),
		Price:             order.Price,
		OrderHash:         orderHash,
	}, nil
}

// Liquidity Probing Methods

func GetPrevTick(bibliophile b.BibliophileClient, input GetPrevTickInput) (*big.Int, error) {
	if input.Tick.Sign() == 0 {
		return nil, errors.New("tick price cannot be zero")
	}
	if input.IsBid {
		currentTick := bibliophile.GetBidsHead(input.Amm)
		if input.Tick.Cmp(currentTick) >= 0 {
			return nil, fmt.Errorf("tick %d is greater than or equal to bidsHead %d", input.Tick, currentTick)
		}
		for {
			nextTick := bibliophile.GetNextBidPrice(input.Amm, currentTick)
			if nextTick.Cmp(input.Tick) <= 0 {
				return currentTick, nil
			}
			currentTick = nextTick
		}
	}
	currentTick := bibliophile.GetAsksHead(input.Amm)
	if currentTick.Sign() == 0 {
		return nil, errors.New("asksHead is zero")
	}
	if input.Tick.Cmp(currentTick) <= 0 {
		return nil, fmt.Errorf("tick %d is less than or equal to asksHead %d", input.Tick, currentTick)
	}
	for {
		nextTick := bibliophile.GetNextAskPrice(input.Amm, currentTick)
		if nextTick.Cmp(input.Tick) >= 0 || nextTick.Sign() == 0 {
			return currentTick, nil
		}
		currentTick = nextTick
	}
}

func SampleImpactBid(bibliophile b.BibliophileClient, ammAddress common.Address) *big.Int {
	impactMarginNotional := bibliophile.GetImpactMarginNotional(ammAddress)
	if impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	return _sampleImpactBid(bibliophile, ammAddress, impactMarginNotional)
}

func _sampleImpactBid(bibliophile b.BibliophileClient, ammAddress common.Address, _impactMarginNotional *big.Int) *big.Int {
	if _impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	impactMarginNotional := new(big.Int).Mul(_impactMarginNotional, big.NewInt(1e12))
	accNotional := big.NewInt(0)
	accBaseQ := big.NewInt(0)
	tick := bibliophile.GetBidsHead(ammAddress)
	for tick.Sign() != 0 {
		amount := bibliophile.GetBidSize(ammAddress, tick)
		accumulator := new(big.Int).Add(accNotional, divide1e6(big.NewInt(0).Mul(amount, tick)))
		if accumulator.Cmp(impactMarginNotional) >= 0 {
			break
		}
		accNotional = accumulator
		accBaseQ.Add(accBaseQ, amount)
		tick = bibliophile.GetNextBidPrice(ammAddress, tick)
	}
	if tick.Sign() == 0 {
		return big.NewInt(0)
	}
	baseQAtTick := new(big.Int).Div(multiply1e6(new(big.Int).Sub(impactMarginNotional, accNotional)), tick)
	return new(big.Int).Div(multiply1e6(impactMarginNotional), new(big.Int).Add(baseQAtTick, accBaseQ)) // return value is in 6 decimals
}

func SampleImpactAsk(bibliophile b.BibliophileClient, ammAddress common.Address) *big.Int {
	impactMarginNotional := bibliophile.GetImpactMarginNotional(ammAddress)
	if impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	return _sampleImpactAsk(bibliophile, ammAddress, impactMarginNotional)
}

func _sampleImpactAsk(bibliophile b.BibliophileClient, ammAddress common.Address, _impactMarginNotional *big.Int) *big.Int {
	if _impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	impactMarginNotional := new(big.Int).Mul(_impactMarginNotional, big.NewInt(1e12))
	tick := bibliophile.GetAsksHead(ammAddress)
	accNotional := big.NewInt(0)
	accBaseQ := big.NewInt(0)
	for tick.Sign() != 0 {
		amount := bibliophile.GetAskSize(ammAddress, tick)
		accumulator := new(big.Int).Add(accNotional, divide1e6(big.NewInt(0).Mul(amount, tick)))
		if accumulator.Cmp(impactMarginNotional) >= 0 {
			break
		}
		accNotional = accumulator
		accBaseQ.Add(accBaseQ, amount)
		tick = bibliophile.GetNextAskPrice(ammAddress, tick)
	}
	if tick.Sign() == 0 {
		return big.NewInt(0)
	}
	baseQAtTick := new(big.Int).Div(multiply1e6(new(big.Int).Sub(impactMarginNotional, accNotional)), tick)
	return new(big.Int).Div(multiply1e6(impactMarginNotional), new(big.Int).Add(baseQAtTick, accBaseQ)) // return value is in 6 decimals
}

func GetQuote(bibliophile b.BibliophileClient, ammAddress common.Address, baseAssetQuantity *big.Int) *big.Int {
	return big.NewInt(0)
}

func GetBaseQuote(bibliophile b.BibliophileClient, ammAddress common.Address, quoteAssetQuantity *big.Int) *big.Int {
	if quoteAssetQuantity.Sign() > 0 { // get the qoute to long quoteQuantity dollars
		return _sampleImpactAsk(bibliophile, ammAddress, quoteAssetQuantity)
	}
	// get the qoute to short quoteQuantity dollars
	return _sampleImpactBid(bibliophile, ammAddress, new(big.Int).Neg(quoteAssetQuantity))
}

// Limit Orders V2
func ValidateCancelLimitOrderV2(bibliophile b.BibliophileClient, inputStruct *ValidateCancelLimitOrderInput) *ValidateCancelLimitOrderOutput {
	errorString, orderHash, ammAddress, unfilledAmount := validateCancelLimitOrderV2(bibliophile, inputStruct.Order, inputStruct.Trader, inputStruct.AssertLowMargin)
	return &ValidateCancelLimitOrderOutput{
		Err:       errorString,
		OrderHash: orderHash,
		Res: IOrderHandlerCancelOrderRes{
			Amm:            ammAddress,
			UnfilledAmount: unfilledAmount,
		},
	}
}

func validateCancelLimitOrderV2(bibliophile b.BibliophileClient, order ILimitOrderBookOrderV2, sender common.Address, assertLowMargin bool) (errorString string, orderHash [32]byte, ammAddress common.Address, unfilledAmount *big.Int) {
	unfilledAmount = big.NewInt(0)
	trader := order.Trader
	if trader != sender && !bibliophile.IsTradingAuthority(trader, sender) {
		errorString = ErrNoTradingAuthority.Error()
		return
	}
	orderHash, err := GetLimitOrderV2Hash(&order)
	if err != nil {
		errorString = err.Error()
		return
	}
	switch status := OrderStatus(bibliophile.GetOrderStatus(orderHash)); status {
	case Invalid:
		errorString = "Invalid"
		return
	case Filled:
		errorString = "Filled"
		return
	case Cancelled:
		errorString = "Cancelled"
		return
	default:
	}
	if assertLowMargin && bibliophile.GetAvailableMargin(trader).Sign() != -1 {
		errorString = "Not Low Margin"
		return
	}
	unfilledAmount = big.NewInt(0).Sub(order.BaseAssetQuantity, bibliophile.GetOrderFilledAmount(orderHash))
	ammAddress = bibliophile.GetMarketAddressFromMarketID(order.AmmIndex.Int64())
	return
}

func ValidatePlaceLimitOrderV2(bibliophile b.BibliophileClient, order ILimitOrderBookOrderV2, trader common.Address) *ValidatePlaceLimitOrderOutput {
	errorString, orderHash, ammAddress, reserveAmount := validatePlaceLimitOrderV2(bibliophile, order, trader)
	return &ValidatePlaceLimitOrderOutput{
		Errs:      errorString,
		Orderhash: orderHash,
		Res: IOrderHandlerPlaceOrderRes{
			Amm:           ammAddress,
			ReserveAmount: reserveAmount,
		},
	}
}

func validatePlaceLimitOrderV2(bibliophile b.BibliophileClient, order ILimitOrderBookOrderV2, sender common.Address) (errorString string, orderHash [32]byte, ammAddress common.Address, reserveAmount *big.Int) {
	reserveAmount = big.NewInt(0)
	orderHash, err := GetLimitOrderV2Hash(&order)
	if err != nil {
		errorString = err.Error()
		return
	}
	if order.Price.Sign() != 1 {
		errorString = ErrInvalidPrice.Error()
		return
	}
	trader := order.Trader
	if trader != sender && !bibliophile.IsTradingAuthority(trader, sender) {
		errorString = ErrNoTradingAuthority.Error()
		return
	}
	ammAddress = bibliophile.GetMarketAddressFromMarketID(order.AmmIndex.Int64())
	if order.BaseAssetQuantity.Sign() == 0 {
		errorString = ErrBaseAssetQuantityZero.Error()
		return
	}
	minSize := bibliophile.GetMinSizeRequirement(order.AmmIndex.Int64())
	if new(big.Int).Mod(order.BaseAssetQuantity, minSize).Sign() != 0 {
		errorString = ErrNotMultiple.Error()
		return
	}
	status := OrderStatus(bibliophile.GetOrderStatus(orderHash))
	if status != Invalid {
		errorString = ErrOrderAlreadyExists.Error()
		return
	}

	posSize := bibliophile.GetSize(ammAddress, &trader)
	reduceOnlyAmount := bibliophile.GetReduceOnlyAmount(trader, order.AmmIndex)
	// this should only happen when a trader with open reduce only orders was liquidated
	if (posSize.Sign() == 0 && reduceOnlyAmount.Sign() != 0) || (posSize.Sign() != 0 && new(big.Int).Mul(posSize, reduceOnlyAmount).Sign() == 1) {
		// if position is non-zero then reduceOnlyAmount should be zero or have the opposite sign as posSize
		errorString = ErrStaleReduceOnlyOrders.Error()
		return
	}

	var orderSide Side = Side(Long)
	if order.BaseAssetQuantity.Sign() == -1 {
		orderSide = Side(Short)
	}
	if order.ReduceOnly {
		if !reducesPosition(posSize, order.BaseAssetQuantity) {
			errorString = ErrReduceOnlyBaseAssetQuantityInvalid.Error()
			return
		}
		longOrdersAmount := bibliophile.GetLongOpenOrdersAmount(trader, order.AmmIndex)
		shortOrdersAmount := bibliophile.GetShortOpenOrdersAmount(trader, order.AmmIndex)
		if (orderSide == Side(Long) && longOrdersAmount.Sign() != 0) || (orderSide == Side(Short) && shortOrdersAmount.Sign() != 0) {
			errorString = ErrOpenOrders.Error()
			return
		}
		if big.NewInt(0).Abs(big.NewInt(0).Add(reduceOnlyAmount, order.BaseAssetQuantity)).Cmp(big.NewInt(0).Abs(posSize)) == 1 {
			errorString = ErrNetReduceOnlyAmountExceeded.Error()
			return
		}
	} else {
		if reduceOnlyAmount.Sign() != 0 && order.BaseAssetQuantity.Sign() != posSize.Sign() {
			errorString = ErrOpenReduceOnlyOrders.Error()
			return
		}
		availableMargin := bibliophile.GetAvailableMargin(trader)
		requiredMargin := getRequiredMargin(bibliophile, order)
		if availableMargin.Cmp(requiredMargin) == -1 {
			errorString = ErrInsufficientMargin.Error()
			return
		}
		reserveAmount = requiredMargin
	}
	if order.PostOnly {
		asksHead := bibliophile.GetAsksHead(ammAddress)
		bidsHead := bibliophile.GetBidsHead(ammAddress)
		if (orderSide == Side(Short) && bidsHead.Sign() != 0 && order.Price.Cmp(bidsHead) != 1) || (orderSide == Side(Long) && asksHead.Sign() != 0 && order.Price.Cmp(asksHead) != -1) {
			errorString = ErrCrossingMarket.Error()
			return
		}
	}
	return
}

func reducesPosition(positionSize *big.Int, baseAssetQuantity *big.Int) bool {
	if positionSize.Sign() == 1 && baseAssetQuantity.Sign() == -1 && big.NewInt(0).Add(positionSize, baseAssetQuantity).Sign() != -1 {
		return true
	}
	if positionSize.Sign() == -1 && baseAssetQuantity.Sign() == 1 && big.NewInt(0).Add(positionSize, baseAssetQuantity).Sign() != 1 {
		return true
	}
	return false
}

func getRequiredMargin(bibliophile b.BibliophileClient, order ILimitOrderBookOrderV2) *big.Int {
	price := order.Price
	upperBound, _ := bibliophile.GetUpperAndLowerBoundForMarket(order.AmmIndex.Int64())
	if order.BaseAssetQuantity.Sign() == -1 && order.Price.Cmp(upperBound) == -1 {
		price = upperBound
	}
	quoteAsset := big.NewInt(0).Abs(big.NewInt(0).Div(new(big.Int).Mul(order.BaseAssetQuantity, price), big.NewInt(1e18)))
	requiredMargin := big.NewInt(0).Div(big.NewInt(0).Mul(bibliophile.GetMinAllowableMargin(), quoteAsset), big.NewInt(1e6))
	takerFee := big.NewInt(0).Div(big.NewInt(0).Mul(quoteAsset, bibliophile.GetTakerFee()), big.NewInt(1e6))
	requiredMargin.Add(requiredMargin, takerFee)
	return requiredMargin
}

func divide1e6(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, big.NewInt(1e6))
}

func multiply1e6(number *big.Int) *big.Int {
	return new(big.Int).Mul(number, big.NewInt(1e6))
}

func formatOrder(orderBytes []byte) interface{} {
	decodeStep0, err := decodeTypeAndEncodedOrder(orderBytes)
	if err != nil {
		return orderBytes
	}

	if decodeStep0.OrderType == Limit {
		order, err := orderbook.DecodeLimitOrder(decodeStep0.EncodedOrder)
		if err != nil {
			return decodeStep0
		}
		orderJson := order.Map()
		orderHash, err := GetLimitOrderHash(order)
		if err != nil {
			return orderJson
		}
		orderJson["hash"] = orderHash.String()
		return orderJson
	}
	if decodeStep0.OrderType == IOC {
		order, err := orderbook.DecodeIOCOrder(decodeStep0.EncodedOrder)
		if err != nil {
			return decodeStep0
		}
		orderJson := order.Map()
		orderHash, err := getIOCOrderHash(order)
		if err != nil {
			return orderJson
		}
		orderJson["hash"] = orderHash.String()
		return orderJson
	}
	if decodeStep0.OrderType == LimitV2 {
		order, err := orderbook.DecodeLimitOrderV2(decodeStep0.EncodedOrder)
		if err != nil {
			return decodeStep0
		}
		orderJson := order.Map()
		orderHash, err := GetLimitOrderV2Hash_2(order)
		if err != nil {
			return orderJson
		}
		orderJson["hash"] = orderHash.String()
		return orderJson
	}
	return nil
}
