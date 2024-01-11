package juror

import (
	"math/big"

	ob "github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	"github.com/ethereum/go-ethereum/common"
)

func ValidatePlaceLimitOrder(bibliophile b.BibliophileClient, inputStruct *ValidatePlaceLimitOrderInput) (response ValidatePlaceLimitOrderOutput) {
	order := inputStruct.Order
	sender := inputStruct.Sender

	response = ValidatePlaceLimitOrderOutput{Res: IOrderHandlerPlaceOrderRes{}}
	response.Res.ReserveAmount = big.NewInt(0)
	orderHash, err := GetLimitOrderHashFromContractStruct(&order)
	response.Orderhash = orderHash

	if err != nil {
		response.Err = err.Error()
		return
	}
	if order.Price.Sign() != 1 {
		response.Err = ErrInvalidPrice.Error()
		return
	}
	trader := order.Trader
	if trader != sender && !bibliophile.IsTradingAuthority(trader, sender) {
		response.Err = ErrNoTradingAuthority.Error()
		return
	}
	ammAddress := bibliophile.GetMarketAddressFromMarketID(order.AmmIndex.Int64())

	if (ammAddress == common.Address{}) {
		response.Err = ErrInvalidMarket.Error()
		return
	}
	response.Res.Amm = ammAddress
	if order.BaseAssetQuantity.Sign() == 0 {
		response.Err = ErrBaseAssetQuantityZero.Error()
		return
	}
	minSize := bibliophile.GetMinSizeRequirement(order.AmmIndex.Int64())
	if new(big.Int).Mod(order.BaseAssetQuantity, minSize).Sign() != 0 {
		response.Err = ErrNotMultiple.Error()
		return
	}
	status := OrderStatus(bibliophile.GetOrderStatus(orderHash))
	if status != Invalid {
		response.Err = ErrOrderAlreadyExists.Error()
		return
	}

	posSize := bibliophile.GetSize(ammAddress, &trader)
	reduceOnlyAmount := bibliophile.GetReduceOnlyAmount(trader, order.AmmIndex)
	// this should only happen when a trader with open reduce only orders was liquidated
	if (posSize.Sign() == 0 && reduceOnlyAmount.Sign() != 0) || (posSize.Sign() != 0 && new(big.Int).Mul(posSize, reduceOnlyAmount).Sign() == 1) {
		// if position is non-zero then reduceOnlyAmount should be zero or have the opposite sign as posSize
		response.Err = ErrStaleReduceOnlyOrders.Error()
		return
	}

	var orderSide Side = Side(Long)
	if order.BaseAssetQuantity.Sign() == -1 {
		orderSide = Side(Short)
	}
	if order.ReduceOnly {
		// a reduce only order should reduce position
		if !reducesPosition(posSize, order.BaseAssetQuantity) {
			response.Err = ErrReduceOnlyBaseAssetQuantityInvalid.Error()
			return
		}
		longOrdersAmount := bibliophile.GetLongOpenOrdersAmount(trader, order.AmmIndex)
		shortOrdersAmount := bibliophile.GetShortOpenOrdersAmount(trader, order.AmmIndex)

		// if the trader is placing a reduceOnly long that means they have a short position
		// we allow only 1 kind of order in the opposite direction of the position
		// otherwise we run the risk of having stale reduceOnly orders (orders that are not actually reducing the position)
		if (orderSide == Side(Long) && longOrdersAmount.Sign() != 0) ||
			(orderSide == Side(Short) && shortOrdersAmount.Sign() != 0) {
			response.Err = ErrOpenOrders.Error()
			return
		}
		if hu.Abs(hu.Add(reduceOnlyAmount, order.BaseAssetQuantity)).Cmp(hu.Abs(posSize)) == 1 {
			response.Err = ErrNetReduceOnlyAmountExceeded.Error()
			return
		}
	} else {
		// we allow only 1 kind of order in the opposite direction of the position
		if order.BaseAssetQuantity.Sign() != posSize.Sign() && reduceOnlyAmount.Sign() != 0 {
			response.Err = ErrOpenReduceOnlyOrders.Error()
			return
		}
		availableMargin := bibliophile.GetAvailableMargin(trader, hu.UpgradeVersionV0orV1(bibliophile.GetTimeStamp()))
		requiredMargin := getRequiredMargin(bibliophile, order)
		if availableMargin.Cmp(requiredMargin) == -1 {
			response.Err = ErrInsufficientMargin.Error()
			return
		}
		response.Res.ReserveAmount = requiredMargin
	}

	if order.PostOnly {
		asksHead := bibliophile.GetAsksHead(ammAddress)
		bidsHead := bibliophile.GetBidsHead(ammAddress)
		if (orderSide == Side(Short) && bidsHead.Sign() != 0 && order.Price.Cmp(bidsHead) != 1) || (orderSide == Side(Long) && asksHead.Sign() != 0 && order.Price.Cmp(asksHead) != -1) {
			response.Err = ErrCrossingMarket.Error()
			return
		}
	}

	if !bibliophile.HasReferrer(order.Trader) {
		response.Err = ErrNoReferrer.Error()
	}

	if hu.Mod(order.Price, bibliophile.GetPriceMultiplier(ammAddress)).Sign() != 0 {
		response.Err = ErrPricePrecision.Error()
		return
	}

	return response
}

func ValidateCancelLimitOrder(bibliophile b.BibliophileClient, inputStruct *ValidateCancelLimitOrderInput) (response ValidateCancelLimitOrderOutput) {
	order := inputStruct.Order
	sender := inputStruct.Sender
	assertLowMargin := inputStruct.AssertLowMargin

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
	if assertLowMargin && bibliophile.GetAvailableMargin(trader, hu.UpgradeVersionV0orV1(bibliophile.GetTimeStamp())).Sign() != -1 {
		response.Err = "Not Low Margin"
		return
	}
	response.Res.UnfilledAmount = big.NewInt(0).Sub(order.BaseAssetQuantity, bibliophile.GetOrderFilledAmount(orderHash))
	response.Res.Amm = bibliophile.GetMarketAddressFromMarketID(order.AmmIndex.Int64())

	return response
}

func ILimitOrderBookOrderToLimitOrder(o *ILimitOrderBookOrder) *ob.LimitOrder {
	return &ob.LimitOrder{
		BaseOrder: hu.BaseOrder{
			AmmIndex:          o.AmmIndex,
			Trader:            o.Trader,
			BaseAssetQuantity: o.BaseAssetQuantity,
			Price:             o.Price,
			Salt:              o.Salt,
			ReduceOnly:        o.ReduceOnly,
		},
		PostOnly: o.PostOnly,
	}
}

func GetLimitOrderHashFromContractStruct(o *ILimitOrderBookOrder) (common.Hash, error) {
	return ILimitOrderBookOrderToLimitOrder(o).Hash()
}
