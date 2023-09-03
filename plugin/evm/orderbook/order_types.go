package orderbook

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

type ContractOrder interface {
	EncodeToABI() ([]byte, error)
	DecodeFromRawOrder(rawOrder interface{})
	Map() map[string]interface{}
}

// LimitOrder type is copy of Order struct defined in LimitOrderbook contract
type LimitOrder struct {
	AmmIndex          *big.Int       `json:"ammIndex"`
	Trader            common.Address `json:"trader"`
	BaseAssetQuantity *big.Int       `json:"baseAssetQuantity"`
	Price             *big.Int       `json:"price"`
	Salt              *big.Int       `json:"salt"`
	ReduceOnly        bool           `json:"reduceOnly"`
}

// LimitOrderV2 type is copy of OrderV2 struct defined in LimitOrderbook contract
type LimitOrderV2 struct {
	LimitOrder
	PostOnly bool `json:"postOnly"`
}

// IOCOrder type is copy of IOCOrder struct defined in Orderbook contract
type IOCOrder struct {
	LimitOrder
	OrderType uint8    `json:"orderType"`
	ExpireAt  *big.Int `json:"expireAt"`
}

// LimitOrder
func (order *LimitOrder) EncodeToABI() ([]byte, error) {
	limitOrderType, err := getOrderType("limit")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	encodedLimitOrder, err := abi.Arguments{{Type: limitOrderType}}.Pack(order)
	if err != nil {
		return nil, fmt.Errorf("limit order packing failed: %w", err)
	}

	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	// 0 means ordertype = limit order
	encodedOrder, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Pack(uint8(0), encodedLimitOrder)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}

	return encodedOrder, nil
}

func (order *LimitOrder) DecodeFromRawOrder(rawOrder interface{}) {
	marshalledOrder, _ := json.Marshal(rawOrder)
	json.Unmarshal(marshalledOrder, &order)
}

func (order *LimitOrder) Map() map[string]interface{} {
	return map[string]interface{}{
		"ammIndex":          order.AmmIndex,
		"trader":            order.Trader,
		"baseAssetQuantity": utils.BigIntToFloat(order.BaseAssetQuantity, 18),
		"price":             utils.BigIntToFloat(order.Price, 6),
		"reduceOnly":        order.ReduceOnly,
		"salt":              order.Salt,
	}
}

func DecodeLimitOrder(encodedOrder []byte) (*LimitOrder, error) {
	limitOrderType, err := getOrderType("limit")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	order, err := abi.Arguments{{Type: limitOrderType}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	limitOrder := &LimitOrder{}
	limitOrder.DecodeFromRawOrder(order[0])
	return limitOrder, nil
}

// LimitOrderV2
func (order *LimitOrderV2) EncodeToABIWithoutType() ([]byte, error) {
	limitOrderV2Type, err := getOrderType("limit_v2")
	if err != nil {
		return nil, err
	}
	encodedLimitOrderV2, err := abi.Arguments{{Type: limitOrderV2Type}}.Pack(order)
	if err != nil {
		return nil, err
	}
	return encodedLimitOrderV2, nil
}

func (order *LimitOrderV2) EncodeToABI() ([]byte, error) {
	encodedLimitOrderV2, err := order.EncodeToABIWithoutType()
	if err != nil {
		return nil, fmt.Errorf("limit order packing failed: %w", err)
	}
	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	// 2 means ordertype = limit order V2
	encodedOrder, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Pack(uint8(2) /* Limit Order v2 */, encodedLimitOrderV2)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}
	return encodedOrder, nil
}

func (order *LimitOrderV2) DecodeFromRawOrder(rawOrder interface{}) {
	marshalledOrder, _ := json.Marshal(rawOrder)
	json.Unmarshal(marshalledOrder, &order)
}

func (order *LimitOrderV2) Map() map[string]interface{} {
	return map[string]interface{}{
		"ammIndex":          order.AmmIndex,
		"trader":            order.Trader,
		"baseAssetQuantity": utils.BigIntToFloat(order.BaseAssetQuantity, 18),
		"price":             utils.BigIntToFloat(order.Price, 6),
		"reduceOnly":        order.ReduceOnly,
		"postOnly":          order.PostOnly,
		"salt":              order.Salt,
	}
}

func DecodeLimitOrderV2(encodedOrder []byte) (*LimitOrderV2, error) {
	limitOrderV2Type, err := getOrderType("limit_v2")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	order, err := abi.Arguments{{Type: limitOrderV2Type}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	limitOrderV2 := &LimitOrderV2{}
	limitOrderV2.DecodeFromRawOrder(order[0])
	return limitOrderV2, nil
}

// ----------------------------------------------------------------------------
// IOCOrder

func (order *IOCOrder) EncodeToABI() ([]byte, error) {
	iocOrderType, err := getOrderType("ioc")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	encodedIOCOrder, err := abi.Arguments{{Type: iocOrderType}}.Pack(order)
	if err != nil {
		return nil, fmt.Errorf("limit order packing failed: %w", err)
	}

	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	// 1 means ordertype = IOC/market order
	encodedOrder, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Pack(uint8(1), encodedIOCOrder)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}

	return encodedOrder, nil
}

func (order *IOCOrder) DecodeFromRawOrder(rawOrder interface{}) {
	marshalledOrder, _ := json.Marshal(rawOrder)
	json.Unmarshal(marshalledOrder, &order)
}

func (order *IOCOrder) Map() map[string]interface{} {
	return map[string]interface{}{
		"ammIndex":          order.AmmIndex,
		"trader":            order.Trader,
		"baseAssetQuantity": utils.BigIntToFloat(order.BaseAssetQuantity, 18),
		"price":             utils.BigIntToFloat(order.Price, 6),
		"reduceOnly":        order.ReduceOnly,
		"salt":              order.Salt,
		"orderType":         order.OrderType,
		"expireAt":          order.ExpireAt,
	}
}

func DecodeIOCOrder(encodedOrder []byte) (*IOCOrder, error) {
	iocOrderType, err := getOrderType("ioc")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	order, err := abi.Arguments{{Type: iocOrderType}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	iocOrder := &IOCOrder{}
	iocOrder.DecodeFromRawOrder(order[0])
	return iocOrder, nil
}

// ----------------------------------------------------------------------------
// Helper functions
func getOrderType(orderType string) (abi.Type, error) {
	if orderType == "limit" {
		return abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "ammIndex", Type: "uint256"},
			{Name: "trader", Type: "address"},
			{Name: "baseAssetQuantity", Type: "int256"},
			{Name: "price", Type: "uint256"},
			{Name: "salt", Type: "uint256"},
			{Name: "reduceOnly", Type: "bool"},
		})
	}
	if orderType == "limit_v2" {
		return abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "ammIndex", Type: "uint256"},
			{Name: "trader", Type: "address"},
			{Name: "baseAssetQuantity", Type: "int256"},
			{Name: "price", Type: "uint256"},
			{Name: "salt", Type: "uint256"},
			{Name: "reduceOnly", Type: "bool"},
			{Name: "postOnly", Type: "bool"},
		})
	}
	if orderType == "ioc" {
		return abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "orderType", Type: "uint8"},
			{Name: "expireAt", Type: "uint256"},
			{Name: "ammIndex", Type: "uint256"},
			{Name: "trader", Type: "address"},
			{Name: "baseAssetQuantity", Type: "int256"},
			{Name: "price", Type: "uint256"},
			{Name: "salt", Type: "uint256"},
			{Name: "reduceOnly", Type: "bool"},
		})
	}
	return abi.Type{}, fmt.Errorf("invalid order type")
}
