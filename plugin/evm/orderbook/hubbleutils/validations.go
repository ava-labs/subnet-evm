package hubbleutils

import (
	"errors"
	// "fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SignedOrderValidationFields struct {
	Now                uint64
	ActiveMarketsCount int64
	MinSize            *big.Int
	PriceMultiplier    *big.Int
	Status             int64
}

var (
	ErrNotSignedOrder        = errors.New("not signed order")
	ErrInvalidPrice          = errors.New("invalid price")
	ErrOrderExpired          = errors.New("order expired")
	ErrBaseAssetQuantityZero = errors.New("baseAssetQuantity is zero")
	ErrNotPostOnly           = errors.New("not post only")
	ErrInvalidMarket         = errors.New("invalid market")
	ErrNotMultiple           = errors.New("not multiple")
	ErrPricePrecision        = errors.New("invalid price precision")
	ErrOrderAlreadyExists    = errors.New("order already exists")
	ErrCrossingMarket        = errors.New("crossing market")
	ErrNoTradingAuthority    = errors.New("no trading authority")
)

// Common Checks
// 1. orderType == Signed
// 2. Not expired
// 3. order should be post only
// 4. baseAssetQuantity is not 0 and multiple of minSize
// 5. price > 0 and price precision check
// 6. signer is valid trading authority
// 7. market is valid
// 8. order is not already filled or cancelled

// Place Order Checks
// P1. Order is not already in memdb (placed)
// P2. Margin is available for non-reduce only orders
// P3. Sum of all reduce only orders should not exceed the total position size (not in state, simply compared to other active orders) and/or opposite direction validations
// P4. Post only order shouldn't cross the market
// P5. HasReferrer

// Matching Order Checks
// M1. order is not being overfilled
// M2. reduce only order should reduce the position size
// M3. HasReferrer
// M4. Not both post only orders are being matched

func ValidateSignedOrder(order *SignedOrder, fields SignedOrderValidationFields) (trader, signer common.Address, err error) {
	if OrderType(order.OrderType) != Signed { // 1.
		err = ErrNotSignedOrder
		return trader, signer, err
	}

	if order.ExpireAt.Uint64() < fields.Now { // 2.
		err = ErrOrderExpired
		return trader, signer, err
	}

	if !order.PostOnly { // 3.
		err = ErrNotPostOnly
		return trader, signer, err
	}

	// 4.
	if order.BaseAssetQuantity.Sign() == 0 {
		err = ErrBaseAssetQuantityZero
		return trader, signer, err
	}
	if new(big.Int).Mod(order.BaseAssetQuantity, fields.MinSize).Sign() != 0 {
		err = ErrNotMultiple
		return trader, signer, err
	}

	if order.Price.Sign() != 1 { // 5.
		err = ErrInvalidPrice
		return trader, signer, err
	}
	if Mod(order.Price, fields.PriceMultiplier).Sign() != 0 {
		err = ErrPricePrecision
		return trader, signer, err
	}

	// 6. caller will perform the check
	orderHash, err := order.Hash()
	if err != nil {
		return trader, signer, err
	}
	signer, err = ECRecover(orderHash.Bytes(), order.Sig[:])
	// fmt.Println("signer", signer)
	if err != nil {
		return trader, signer, err
	}
	trader = order.Trader

	// assumes all markets are active and in sequential order
	if order.AmmIndex.Int64() >= fields.ActiveMarketsCount { // 7.
		err = ErrInvalidMarket
		return trader, signer, err
	}

	if OrderStatus(fields.Status) != Invalid { // 8.
		err = ErrOrderAlreadyExists
		return trader, signer, err
	}
	return trader, signer, nil
}
