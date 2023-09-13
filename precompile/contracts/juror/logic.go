package juror

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

type OrderType uint8

// has to be exact same as expected in contracts
const (
	Limit OrderType = iota
	IOC
	PostOnly
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
	OrderType         OrderType
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
	ErrTooLow                             = errors.New("long price below lower bound")
	ErrTooHigh                            = errors.New("short price above upper bound")
	ErrOverFill                           = errors.New("overfill")
	ErrReduceOnlyAmountExceeded           = errors.New("not reducing pos")
	ErrBaseAssetQuantityZero              = errors.New("baseAssetQuantity is zero")
	ErrReduceOnlyBaseAssetQuantityInvalid = errors.New("reduce only order must reduce position")
	ErrNetReduceOnlyAmountExceeded        = errors.New("net reduce only amount exceeded")
	ErrStaleReduceOnlyOrders              = errors.New("cancel stale reduce only orders")
	ErrInsufficientMargin                 = errors.New("insufficient margin")
	ErrCrossingMarket                     = errors.New("crossing market")
	ErrIOCOrderExpired                    = errors.New("IOC order expired")
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
	if err != nil {
		return nil, err
	}
	m0, err := validateOrder(bibliophile, decodeStep0.OrderType, decodeStep0.EncodedOrder, Long, inputStruct.FillAmount)
	if err != nil {
		return nil, err
	}

	decodeStep1, err := decodeTypeAndEncodedOrder(inputStruct.Data[1])
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

	fillPriceAndModes, err := determineFillPrice(bibliophile, m0, m1)
	if err != nil {
		return nil, err
	}

	output := &ValidateOrdersAndDetermineFillPriceOutput{
		Instructions: [2]IClearingHouseInstruction{
			IClearingHouseInstruction{
				AmmIndex:  m0.AmmIndex,
				Trader:    m0.Trader,
				OrderHash: m0.OrderHash,
				Mode:      uint8(fillPriceAndModes.Mode0),
			},
			IClearingHouseInstruction{
				AmmIndex:  m1.AmmIndex,
				Trader:    m1.Trader,
				OrderHash: m1.OrderHash,
				Mode:      uint8(fillPriceAndModes.Mode1),
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

type executionMode uint8

// DO NOT change this ordering because it is critical for the clearing house to determine the correct fill mode
const (
	Taker executionMode = iota
	Maker
)

type FillPriceAndModes struct {
	FillPrice *big.Int
	Mode0     executionMode
	Mode1     executionMode
}

func determineFillPrice(bibliophile b.BibliophileClient, m0, m1 *Metadata) (*FillPriceAndModes, error) {
	output := FillPriceAndModes{}
	upperBound, lowerBound := bibliophile.GetUpperAndLowerBoundForMarket(m0.AmmIndex.Int64())
	if m0.Price.Cmp(lowerBound) == -1 {
		return nil, ErrTooLow
	}
	if m1.Price.Cmp(upperBound) == 1 {
		return nil, ErrTooHigh
	}

	blockDiff := m0.BlockPlaced.Cmp(m1.BlockPlaced)
	if blockDiff == -1 {
		// order0 came first, can't be IOC order
		if m0.OrderType == IOC {
			return nil, ErrIOCOrderExpired
		}
		// order1 came second, can't be post only order
		if m1.OrderType == PostOnly {
			return nil, ErrCrossingMarket
		}
		output.Mode0 = Maker
		output.Mode1 = Taker
	} else if blockDiff == 1 {
		// order1 came first, can't be IOC order
		if m1.OrderType == IOC {
			return nil, ErrIOCOrderExpired
		}
		// order0 came second, can't be post only order
		if m0.OrderType == PostOnly {
			return nil, ErrCrossingMarket
		}
		output.Mode0 = Taker
		output.Mode1 = Maker
	} else {
		// both orders were placed in same block
		if m1.OrderType == IOC {
			// order1 is IOC, order0 is Limit or post only
			output.Mode0 = Maker
			output.Mode1 = Taker
		} else {
			// scenarios:
			// 1. order0 is IOC, order1 is Limit or post only
			// 2. both order0 and order1 are Limit or post only (in that scenario we default to long being the taker order, which can sometimes result in a better execution price for them)
			output.Mode0 = Taker
			output.Mode1 = Maker
		}
	}

	if output.Mode0 == Maker {
		output.FillPrice = utils.BigIntMin(m0.Price, upperBound)
	} else {
		output.FillPrice = utils.BigIntMax(m1.Price, lowerBound)
	}
	return &output, nil
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

	fillPrice, err := determineLiquidationFillPrice(bibliophile, m0)
	if err != nil {
		return nil, err
	}

	output := &ValidateLiquidationOrderAndDetermineFillPriceOutput{
		Instruction: IClearingHouseInstruction{
			AmmIndex:  m0.AmmIndex,
			Trader:    m0.Trader,
			OrderHash: m0.OrderHash,
			Mode:      uint8(Maker),
		},
		OrderType:    uint8(decodeStep0.OrderType),
		EncodedOrder: decodeStep0.EncodedOrder,
		FillPrice:    fillPrice,
		FillAmount:   fillAmount,
	}
	return output, nil
}

func determineLiquidationFillPrice(bibliophile b.BibliophileClient, m0 *Metadata) (*big.Int, error) {
	liqUpperBound, liqLowerBound := bibliophile.GetAcceptableBoundsForLiquidation(m0.AmmIndex.Int64())
	upperBound, lowerBound := bibliophile.GetUpperAndLowerBoundForMarket(m0.AmmIndex.Int64())
	if m0.BaseAssetQuantity.Sign() > 0 {
		// we are liquidating a long position
		// do not allow liquidation if order.Price < liqLowerBound, because that gives scope for malicious activity to a validator
		if m0.Price.Cmp(liqLowerBound) == -1 {
			return nil, ErrTooLow
		}
		return utils.BigIntMin(m0.Price, upperBound /* oracle spread upper bound */), nil
	}

	// we are liquidating a short position
	if m0.Price.Cmp(liqUpperBound) == 1 {
		return nil, ErrTooHigh
	}
	return utils.BigIntMax(m0.Price, lowerBound /* oracle spread lower bound */), nil
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
		metadata, err := validateExecuteLimitOrder(bibliophile, order, side, fillAmount, orderHash)
		if order.PostOnly {
			metadata.OrderType = PostOnly
		}
		if err != nil {
			return nil, err
		}
	}
	if orderType == IOC {
		order, err := orderbook.DecodeIOCOrder(encodedOrder)
		if err != nil {
			return nil, err
		}
		return validateExecuteIOCOrder(bibliophile, order, side, fillAmount)
	}
	return nil, errors.New("invalid order type")
}

// Limit Orders

func validateExecuteLimitOrder(bibliophile b.BibliophileClient, order *orderbook.LimitOrder, side Side, fillAmount *big.Int, orderHash [32]byte) (metadata *Metadata, err error) {
	if err := validateLimitOrderLike(bibliophile, &order.BaseOrder, bibliophile.GetOrderFilledAmount(orderHash), OrderStatus(bibliophile.GetOrderStatus(orderHash)), side, fillAmount); err != nil {
		return nil, err
	}
	return &Metadata{
		AmmIndex:          order.AmmIndex,
		Trader:            order.Trader,
		BaseAssetQuantity: order.BaseAssetQuantity,
		BlockPlaced:       bibliophile.GetBlockPlaced(orderHash),
		Price:             order.Price,
		OrderHash:         orderHash,
		OrderType:         Limit,
	}, nil
}

func validateLimitOrderLike(bibliophile b.BibliophileClient, order *orderbook.BaseOrder, filledAmount *big.Int, status OrderStatus, side Side, fillAmount *big.Int) error {
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
		orderHashes[i], err = GetIOCOrderHash(&orderbook.IOCOrder{
			OrderType: order.OrderType,
			ExpireAt:  order.ExpireAt,
			BaseOrder: orderbook.BaseOrder{
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
	orderHash, err := GetIOCOrderHash(order)
	if err != nil {
		return nil, err
	}
	if err := validateLimitOrderLike(bibliophile, &order.BaseOrder, bibliophile.IOC_GetOrderFilledAmount(orderHash), OrderStatus(bibliophile.IOC_GetOrderStatus(orderHash)), side, fillAmount); err != nil {
		return nil, err
	}
	return &Metadata{
		AmmIndex:          order.AmmIndex,
		Trader:            order.Trader,
		BaseAssetQuantity: order.BaseAssetQuantity,
		BlockPlaced:       bibliophile.IOC_GetBlockPlaced(orderHash),
		Price:             order.Price,
		OrderHash:         orderHash,
		OrderType:         IOC,
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

// Limit Orders
func ValidateCancelLimitOrder(bibliophile b.BibliophileClient, inputStruct *ValidateCancelLimitOrderInput) (response *ValidateCancelLimitOrderOutput) {
	order := inputStruct.Order
	sender := inputStruct.Trader
	assertLowMargin := inputStruct.AssertLowMargin

	response = &ValidateCancelLimitOrderOutput{Res: IOrderHandlerCancelOrderRes{}}
	response.Res.UnfilledAmount = big.NewInt(0)

	trader := order.Trader
	if (!assertLowMargin && trader != sender && !bibliophile.IsTradingAuthority(trader, sender)) ||
		(assertLowMargin && !bibliophile.IsValidator(sender)) {
		response.Err = ErrNoTradingAuthority.Error()
		return
	}
	orderHash, err := GetLimitOrderHashFromContractStruct(&order)
	response.OrderHash = orderHash
	if err != nil {
		response.Err = err.Error()
		return
	}
	switch status := OrderStatus(bibliophile.GetOrderStatus(orderHash)); status {
	case Invalid:
		response.Err = "Invalid"
		return
	case Filled:
		response.Err = "Filled"
		return
	case Cancelled:
		response.Err = "Cancelled"
		return
	default:
	}
	if assertLowMargin && bibliophile.GetAvailableMargin(trader).Sign() != -1 {
		response.Err = "Not Low Margin"
		return
	}
	response.Res.UnfilledAmount = big.NewInt(0).Sub(order.BaseAssetQuantity, bibliophile.GetOrderFilledAmount(orderHash))
	response.Res.Amm = bibliophile.GetMarketAddressFromMarketID(order.AmmIndex.Int64())

	return response
}

func ValidatePlaceLimitOrder(bibliophile b.BibliophileClient, order ILimitOrderBookOrderV2, sender common.Address) (response *ValidatePlaceLimitOrderOutput) {
	response = &ValidatePlaceLimitOrderOutput{Res: IOrderHandlerPlaceOrderRes{}}
	response.Res.ReserveAmount = big.NewInt(0)
	orderHash, err := GetLimitOrderHashFromContractStruct(&order)
	response.Orderhash = orderHash

	if err != nil {
		response.Errs = err.Error()
		return
	}
	if order.Price.Sign() != 1 {
		response.Errs = ErrInvalidPrice.Error()
		return
	}
	trader := order.Trader
	if trader != sender && !bibliophile.IsTradingAuthority(trader, sender) {
		response.Errs = ErrNoTradingAuthority.Error()
		return
	}
	ammAddress := bibliophile.GetMarketAddressFromMarketID(order.AmmIndex.Int64())
	response.Res.Amm = ammAddress
	if order.BaseAssetQuantity.Sign() == 0 {
		response.Errs = ErrBaseAssetQuantityZero.Error()
		return
	}
	minSize := bibliophile.GetMinSizeRequirement(order.AmmIndex.Int64())
	if new(big.Int).Mod(order.BaseAssetQuantity, minSize).Sign() != 0 {
		response.Errs = ErrNotMultiple.Error()
		return
	}
	status := OrderStatus(bibliophile.GetOrderStatus(orderHash))
	if status != Invalid {
		response.Errs = ErrOrderAlreadyExists.Error()
		return
	}

	posSize := bibliophile.GetSize(ammAddress, &trader)
	reduceOnlyAmount := bibliophile.GetReduceOnlyAmount(trader, order.AmmIndex)
	// this should only happen when a trader with open reduce only orders was liquidated
	if (posSize.Sign() == 0 && reduceOnlyAmount.Sign() != 0) || (posSize.Sign() != 0 && new(big.Int).Mul(posSize, reduceOnlyAmount).Sign() == 1) {
		// if position is non-zero then reduceOnlyAmount should be zero or have the opposite sign as posSize
		response.Errs = ErrStaleReduceOnlyOrders.Error()
		return
	}

	var orderSide Side = Side(Long)
	if order.BaseAssetQuantity.Sign() == -1 {
		orderSide = Side(Short)
	}
	if order.ReduceOnly {
		if !reducesPosition(posSize, order.BaseAssetQuantity) {
			response.Errs = ErrReduceOnlyBaseAssetQuantityInvalid.Error()
			return
		}
		longOrdersAmount := bibliophile.GetLongOpenOrdersAmount(trader, order.AmmIndex)
		shortOrdersAmount := bibliophile.GetShortOpenOrdersAmount(trader, order.AmmIndex)
		if (orderSide == Side(Long) && longOrdersAmount.Sign() != 0) || (orderSide == Side(Short) && shortOrdersAmount.Sign() != 0) {
			response.Errs = ErrOpenOrders.Error()
			return
		}
		if big.NewInt(0).Abs(big.NewInt(0).Add(reduceOnlyAmount, order.BaseAssetQuantity)).Cmp(big.NewInt(0).Abs(posSize)) == 1 {
			response.Errs = ErrNetReduceOnlyAmountExceeded.Error()
			return
		}
	} else {
		if reduceOnlyAmount.Sign() != 0 && order.BaseAssetQuantity.Sign() != posSize.Sign() {
			response.Errs = ErrOpenReduceOnlyOrders.Error()
			return
		}
		availableMargin := bibliophile.GetAvailableMargin(trader)
		requiredMargin := getRequiredMargin(bibliophile, order)
		if availableMargin.Cmp(requiredMargin) == -1 {
			response.Errs = ErrInsufficientMargin.Error()
			return
		}
		response.Res.ReserveAmount = requiredMargin
	}
	if order.PostOnly {
		asksHead := bibliophile.GetAsksHead(ammAddress)
		bidsHead := bibliophile.GetBidsHead(ammAddress)
		if (orderSide == Side(Short) && bidsHead.Sign() != 0 && order.Price.Cmp(bidsHead) != 1) || (orderSide == Side(Long) && asksHead.Sign() != 0 && order.Price.Cmp(asksHead) != -1) {
			response.Errs = ErrCrossingMarket.Error()
			return
		}
	}
	return response
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
		orderHash, err := GetIOCOrderHash(order)
		if err != nil {
			return orderJson
		}
		orderJson["hash"] = orderHash.String()
		return orderJson
	}
	return nil
}
