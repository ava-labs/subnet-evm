package juror

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	"github.com/ethereum/go-ethereum/common"
)

type OrderType uint8

// has to be exact same as expected in contracts
const (
	Limit OrderType = iota
	IOC
)

type DecodeStep struct {
	OrderType    OrderType
	EncodedOrder []byte
}

type LimitOrder orderbook.LimitOrder

type IOCOrder struct {
	LimitOrder
	OrderType OrderType
	expireAt  *big.Int
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
	ErrInvalidFillAmount = errors.New("invalid fillAmount")
	ErrNotLongOrder      = errors.New("not long")
	ErrNotShortOrder     = errors.New("not short")
	ErrNotSameAMM        = errors.New("OB_orders_for_different_amms")
	ErrNoMatch           = errors.New("OB_orders_do_not_match")
	ErrNotMultiple       = errors.New("not multiple")

	ErrInvalidOrder             = errors.New("invalid order")
	ErrTooLow                   = errors.New("OB_long_order_price_too_low")
	ErrTooHigh                  = errors.New("OB_short_order_price_too_high")
	ErrOverFill                 = errors.New("overfill")
	ErrReduceOnlyAmountExceeded = errors.New("not reducing pos")
)

// Business Logic
func ValidateOrdersAndDetermineFillPrice(bibliophile b.BibliophileClient, inputStruct *ValidateOrdersAndDetermineFillPriceInput) (*ValidateOrdersAndDetermineFillPriceOutput, error) {
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
		order, err := decodeLimitOrder(encodedOrder)
		if err != nil {
			return nil, err
		}
		return validateExecuteLimitOrder(bibliophile, order, side, fillAmount)
	}
	if orderType == IOC {
		order, err := decodeIOCOrder(encodedOrder)
		if err != nil {
			return nil, err
		}
		return validateExecuteIOCOrder(bibliophile, order, side, fillAmount)
	}
	return nil, errors.New("invalid order type")
}

// Limit Orders

func decodeLimitOrder(encodedOrder []byte) (*LimitOrder, error) {
	limitOrderType, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "ammIndex", Type: "uint256"},
		{Name: "trader", Type: "address"},
		{Name: "baseAssetQuantity", Type: "int256"},
		{Name: "price", Type: "uint256"},
		{Name: "salt", Type: "uint256"},
		{Name: "reduceOnly", Type: "bool"},
	})
	order, err := abi.Arguments{{Type: limitOrderType}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	source, ok := order[0].(struct {
		AmmIndex          *big.Int       `json:"ammIndex"`
		Trader            common.Address `json:"trader"`
		BaseAssetQuantity *big.Int       `json:"baseAssetQuantity"`
		Price             *big.Int       `json:"price"`
		Salt              *big.Int       `json:"salt"`
		ReduceOnly        bool           `json:"reduceOnly"`
	})
	if !ok {
		return nil, errors.New("couldnt decode limit order")
	}
	fmt.Println(source)
	return &LimitOrder{
		AmmIndex:          source.AmmIndex,
		Trader:            source.Trader,
		BaseAssetQuantity: source.BaseAssetQuantity,
		Price:             source.Price,
		Salt:              source.Salt,
		ReduceOnly:        source.ReduceOnly,
	}, nil
}

func validateExecuteLimitOrder(bibliophile b.BibliophileClient, order *LimitOrder, side Side, fillAmount *big.Int) (metadata *Metadata, err error) {
	orderHash, err := GetLimitOrderHash(order)
	if err != nil {
		return nil, err
	}
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

func validateLimitOrderLike(bibliophile b.BibliophileClient, order *LimitOrder, filledAmount *big.Int, status OrderStatus, side Side, fillAmount *big.Int) error {
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
func decodeIOCOrder(encodedOrder []byte) (*IOCOrder, error) {
	// @todo
	return nil, nil
}

func validateExecuteIOCOrder(bibliophile b.BibliophileClient, order *IOCOrder, side Side, fillAmount *big.Int) (metadata *Metadata, err error) {
	if order.OrderType != IOC {
		return nil, errors.New("not ioc order")
	}
	if order.expireAt.Cmp(bibliophile.GetAccessibleState().GetBlockContext().Timestamp()) < 0 {
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
