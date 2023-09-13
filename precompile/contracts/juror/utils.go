package juror

import (
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
)

func GetLimitOrderHashFromContractStruct(o *ILimitOrderBookOrderV2) (common.Hash, error) {
	order := &orderbook.LimitOrder{
		BaseOrder: orderbook.BaseOrder{
			AmmIndex:          o.AmmIndex,
			BaseAssetQuantity: o.BaseAssetQuantity,
			Price:             o.Price,
			Salt:              o.Salt,
			ReduceOnly:        o.ReduceOnly,
			Trader:            o.Trader,
		},
		PostOnly: o.PostOnly,
	}
	return GetLimitOrderHash(order)
}

func GetLimitOrderHash(order *orderbook.LimitOrder) (common.Hash, error) {
	data, err := order.EncodeToABIWithoutType()
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(crypto.Keccak256(data)), nil
}

func GetIOCOrderHash(o *orderbook.IOCOrder) (hash common.Hash, err error) {
	data, err := o.EncodeToABIWithoutType()
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(crypto.Keccak256(data)), nil
}
