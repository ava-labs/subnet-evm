package hubbleutils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type SignedOrder struct {
	LimitOrder
	OrderType uint8    `json:"orderType"`
	ExpireAt  *big.Int `json:"expireAt"`
	Sig       []byte   `json:"sig"`
}

var (
	ChainId           int64
	VerifyingContract string
)

func SetChainIdAndVerifyingSignedOrdersContract(chainId int64, verifyingContract string) {
	ChainId = chainId
	VerifyingContract = verifyingContract
}

func (order *SignedOrder) EncodeToABIWithoutType() ([]byte, error) {
	signedOrderType, err := getOrderType("signed")
	if err != nil {
		return nil, err
	}
	bytesTy, _ := abi.NewType("bytes", "bytes", nil)
	encodedOrder, err := abi.Arguments{{Type: signedOrderType}, {Type: bytesTy}}.Pack(order, order.Sig)
	if err != nil {
		return nil, err
	}
	return encodedOrder, nil
}

func (order *SignedOrder) EncodeToABI() ([]byte, error) {
	encodedSignedOrder, err := order.EncodeToABIWithoutType()
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}

	uint8Ty, _ := abi.NewType("uint8", "uint8", nil)
	bytesTy, _ := abi.NewType("bytes", "bytes", nil)

	encodedOrder, err := abi.Arguments{{Type: uint8Ty}, {Type: bytesTy}}.Pack(uint8(Signed), encodedSignedOrder)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}

	return encodedOrder, nil
}

func DecodeSignedOrder(encodedOrder []byte) (*SignedOrder, error) {
	signedOrderType, err := getOrderType("signed")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	bytesTy, _ := abi.NewType("bytes", "bytes", nil)
	decodedValues, err := abi.Arguments{{Type: signedOrderType}, {Type: bytesTy}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	signedOrder := &SignedOrder{
		Sig: decodedValues[1].([]byte),
	}
	signedOrder.DecodeFromRawOrder(decodedValues[0])
	return signedOrder, nil
}

func (order *SignedOrder) DecodeFromRawOrder(rawOrder interface{}) {
	marshalledOrder, _ := json.Marshal(rawOrder)
	// fmt.Println("marshalledOrder", string(marshalledOrder))
	err := json.Unmarshal(marshalledOrder, &order)
	if err != nil {
		log.Error("err in DecodeFromRawOrder", "err", err, "rawOrder", rawOrder)
	}
}

func (o *SignedOrder) String() string {
	return fmt.Sprintf(
		"Order: %s, OrderType: %d, ExpireAt: %d, Sig: %s",
		o.LimitOrder.String(),
		o.OrderType,
		o.ExpireAt,
		hex.EncodeToString(o.Sig),
	)
}

func (o *SignedOrder) Hash() (hash common.Hash, err error) {
	if VerifyingContract == "" || ChainId == 0 {
		return common.Hash{}, fmt.Errorf("ChainId or VerifyingContract not set")
	}
	message := map[string]interface{}{
		"orderType":         strconv.FormatUint(uint64(o.OrderType), 10),
		"expireAt":          o.ExpireAt.String(),
		"ammIndex":          o.AmmIndex.String(),
		"trader":            o.Trader.String(),
		"baseAssetQuantity": o.BaseAssetQuantity.String(),
		"price":             o.Price.String(),
		"salt":              o.Salt.String(),
		"reduceOnly":        o.ReduceOnly,
		"postOnly":          o.PostOnly,
	}
	domain := apitypes.TypedDataDomain{
		Name:              "Hubble",
		Version:           "2.0",
		ChainId:           math.NewHexOrDecimal256(ChainId),
		VerifyingContract: VerifyingContract,
	}
	typedData := apitypes.TypedData{
		Types:       Eip712OrderTypes,
		PrimaryType: "Order",
		Domain:      domain,
		Message:     message,
	}
	return EncodeForSigning(typedData)
}

// Trading API methods

// func (o *SignedOrder) UnmarshalJSON(data []byte) error {
// 	// Redefine the structs with simple types for JSON unmarshalling
// 	aux := &struct {
// 		AmmIndex          uint64         `json:"ammIndex"`
// 		Trader            common.Address `json:"trader"`
// 		BaseAssetQuantity string         `json:"baseAssetQuantity"`
// 		Price             string         `json:"price"`
// 		Salt              string         `json:"salt"`
// 		ReduceOnly        bool           `json:"reduceOnly"`
// 		PostOnly          bool           `json:"postOnly"`
// 		OrderType         uint8          `json:"orderType"`
// 		ExpireAt          uint64         `json:"expireAt"`
// 		Sig               string         `json:"sig"`
// 	}{}

// 	// Perform the unmarshalling
// 	if err := json.Unmarshal(data, aux); err != nil {
// 		return err
// 	}

// 	// Convert and assign the values to the original struct
// 	o.AmmIndex = new(big.Int).SetUint64(aux.AmmIndex)

// 	o.Trader = aux.Trader

// 	o.BaseAssetQuantity = new(big.Int)
// 	o.BaseAssetQuantity.SetString(aux.BaseAssetQuantity, 10)

// 	o.Price = new(big.Int)
// 	o.Price.SetString(aux.Price, 10)

// 	o.Salt = new(big.Int)
// 	o.Salt.SetBytes(common.FromHex(aux.Salt))

// 	o.ReduceOnly = aux.ReduceOnly
// 	o.PostOnly = aux.PostOnly
// 	o.OrderType = aux.OrderType

// 	o.ExpireAt = new(big.Int).SetUint64(aux.ExpireAt)
// 	o.Sig = common.FromHex(aux.Sig)
// 	return nil
// }

// func (order *SignedOrder) DecodeAPIOrder(rawOrder interface{}) error {
// 	order_, ok := rawOrder.(string)
// 	if !ok {
// 		fmt.Println("invalid data format")
// 	}

// 	orderJson := []byte(order_)
// 	err := json.Unmarshal(orderJson, &order)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
