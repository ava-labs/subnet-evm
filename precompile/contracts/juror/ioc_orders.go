package juror

import (
	"errors"
	"math/big"

	ob "github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
)

func ValidatePlaceIOCorder(bibliophile b.BibliophileClient, inputStruct *ValidatePlaceIOCOrderInput) (response ValidatePlaceIOCOrderOutput) {
	order := inputStruct.Order
	trader := order.Trader

	var err error
	response.OrderHash, err = IImmediateOrCancelOrdersOrderToIOCOrder(&inputStruct.Order).Hash()
	if err != nil {
		response.Err = err.Error()
		return
	}

	if trader != inputStruct.Sender && !bibliophile.IsTradingAuthority(trader, inputStruct.Sender) {
		response.Err = ErrNoTradingAuthority.Error()
		return
	}
	blockTimestamp := bibliophile.GetAccessibleState().GetBlockContext().Timestamp()
	expireWithin := blockTimestamp + bibliophile.IOC_GetExpirationCap().Uint64()
	if order.BaseAssetQuantity.Sign() == 0 {
		response.Err = ErrInvalidFillAmount.Error()
		return
	}
	if ob.OrderType(order.OrderType) != ob.IOC {
		response.Err = errors.New("not_ioc_order").Error()
		return
	}
	if order.ExpireAt.Uint64() < blockTimestamp {
		response.Err = errors.New("ioc expired").Error()
		return
	}
	if order.ExpireAt.Uint64() > expireWithin {
		response.Err = errors.New("ioc expiration too far").Error()
		return
	}
	minSize := bibliophile.GetMinSizeRequirement(order.AmmIndex.Int64())
	if new(big.Int).Mod(order.BaseAssetQuantity, minSize).Sign() != 0 {
		response.Err = ErrNotMultiple.Error()
		return
	}

	if OrderStatus(bibliophile.IOC_GetOrderStatus(response.OrderHash)) != Invalid {
		response.Err = ErrInvalidOrder.Error()
		return
	}
	// this check is sort of redundant because either ways user can circumvent this by placing several reduceOnly order in a single tx/block
	// if order.ReduceOnly {
	// @todo
	// }
	return response
}

func IImmediateOrCancelOrdersOrderToIOCOrder(order *IImmediateOrCancelOrdersOrder) *ob.IOCOrder {
	return &ob.IOCOrder{
		BaseOrder: ob.BaseOrder{
			AmmIndex:          order.AmmIndex,
			Trader:            order.Trader,
			BaseAssetQuantity: order.BaseAssetQuantity,
			Price:             order.Price,
			Salt:              order.Salt,
			ReduceOnly:        order.ReduceOnly,
		},
		OrderType: order.OrderType,
		ExpireAt:  order.ExpireAt,
	}
}
