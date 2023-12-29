package juror

import (
	"errors"
	"math/big"
	"testing"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	"github.com/ethereum/go-ethereum/common"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type ValidatePlaceIOCOrderTestCase struct {
	Order  IImmediateOrCancelOrdersOrder
	Sender common.Address
	Error  error // response error
}

func testValidatePlaceIOCOrderTestCase(t *testing.T, mockBibliophile *b.MockBibliophileClient, c ValidatePlaceIOCOrderTestCase) {
	testInput := ValidatePlaceIOCOrderInput{
		Order:  c.Order,
		Sender: c.Sender,
	}

	// call precompile
	response := ValidatePlaceIOCorder(mockBibliophile, &testInput)

	// verify results
	if c.Error == nil && response.Err != "" {
		t.Fatalf("expected no error, got %v", response.Err)
	}
	if c.Error != nil && response.Err != c.Error.Error() {
		t.Fatalf("expected %v, got %v", c.Error, response.Err)
	}
}

func TestValidatePlaceIOCOrder(t *testing.T) {
	trader := common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("no trading authority", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(false)

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order: IImmediateOrCancelOrdersOrder{
				OrderType:         1,
				ExpireAt:          big.NewInt(0),
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(5),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(1),
				ReduceOnly:        false,
			},
			Sender: common.Address{1},
			Error:  ErrNoTradingAuthority,
		})
	})

	t.Run("invalid fill amount", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order: IImmediateOrCancelOrdersOrder{
				OrderType:         1,
				ExpireAt:          big.NewInt(0),
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(0),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(2),
				ReduceOnly:        false,
			},
			Sender: common.Address{1},
			Error:  ErrInvalidFillAmount,
		})
	})

	t.Run("not IOC order", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order: IImmediateOrCancelOrdersOrder{
				OrderType:         0,
				ExpireAt:          big.NewInt(0),
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(5),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(3),
				ReduceOnly:        false,
			},
			Sender: common.Address{1},
			Error:  ErrNotIOCOrder,
		})
	})

	t.Run("ioc expired", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order: IImmediateOrCancelOrdersOrder{
				OrderType:         1,
				ExpireAt:          big.NewInt(900),
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(5),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(4),
				ReduceOnly:        false,
			},
			Sender: common.Address{1},
			Error:  errors.New("ioc expired"),
		})
	})

	t.Run("ioc expiration too far", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order: IImmediateOrCancelOrdersOrder{
				OrderType:         1,
				ExpireAt:          big.NewInt(1006),
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(5),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(5),
				ReduceOnly:        false,
			},
			Sender: common.Address{1},
			Error:  errors.New("ioc expiration too far"),
		})
	})

	t.Run("not multiple", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order: IImmediateOrCancelOrdersOrder{
				OrderType:         1,
				ExpireAt:          big.NewInt(1004),
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(7),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(6),
				ReduceOnly:        false,
			},
			Sender: common.Address{1},
			Error:  ErrNotMultiple,
		})
	})

	t.Run("invalid order", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		order := IImmediateOrCancelOrdersOrder{
			OrderType:         1,
			ExpireAt:          big.NewInt(1004),
			AmmIndex:          big.NewInt(0),
			Trader:            trader,
			BaseAssetQuantity: big.NewInt(10),
			Price:             big.NewInt(100),
			Salt:              big.NewInt(7),
			ReduceOnly:        false,
		}
		hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(1))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order:  order,
			Sender: common.Address{1},
			Error:  ErrInvalidOrder,
		})
	})

	t.Run("no referrer", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		order := IImmediateOrCancelOrdersOrder{
			OrderType:         1,
			ExpireAt:          big.NewInt(1004),
			AmmIndex:          big.NewInt(0),
			Trader:            trader,
			BaseAssetQuantity: big.NewInt(10),
			Price:             big.NewInt(100),
			Salt:              big.NewInt(8),
			ReduceOnly:        false,
		}
		hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(0))
		mockBibliophile.EXPECT().HasReferrer(trader).Return(false)

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order:  order,
			Sender: common.Address{1},
			Error:  ErrNoReferrer,
		})
	})

	t.Run("invalid market", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		order := IImmediateOrCancelOrdersOrder{
			OrderType:         1,
			ExpireAt:          big.NewInt(1004),
			AmmIndex:          big.NewInt(0),
			Trader:            trader,
			BaseAssetQuantity: big.NewInt(10),
			Price:             big.NewInt(100),
			Salt:              big.NewInt(9),
			ReduceOnly:        false,
		}
		hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(0))
		mockBibliophile.EXPECT().HasReferrer(trader).Return(true)
		mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(common.Address{})

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order:  order,
			Sender: common.Address{1},
			Error:  ErrInvalidMarket,
		})
	})

	t.Run("reduce only - doesn't reduce position", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		order := IImmediateOrCancelOrdersOrder{
			OrderType:         1,
			ExpireAt:          big.NewInt(1004),
			AmmIndex:          big.NewInt(0),
			Trader:            trader,
			BaseAssetQuantity: big.NewInt(10),
			Price:             big.NewInt(100),
			Salt:              big.NewInt(10),
			ReduceOnly:        true,
		}
		hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(0))
		mockBibliophile.EXPECT().HasReferrer(trader).Return(true)
		mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(common.Address{101})
		mockBibliophile.EXPECT().GetSize(common.Address{101}, &trader).Return(big.NewInt(-5))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order:  order,
			Sender: common.Address{1},
			Error:  ErrReduceOnlyBaseAssetQuantityInvalid,
		})
	})

	t.Run("reduce only - reduce only amount exceeded", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		order := IImmediateOrCancelOrdersOrder{
			OrderType:         1,
			ExpireAt:          big.NewInt(1004),
			AmmIndex:          big.NewInt(0),
			Trader:            trader,
			BaseAssetQuantity: big.NewInt(10),
			Price:             big.NewInt(100),
			Salt:              big.NewInt(11),
			ReduceOnly:        true,
		}
		hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(0))
		mockBibliophile.EXPECT().HasReferrer(trader).Return(true)
		mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(common.Address{101})
		mockBibliophile.EXPECT().GetSize(common.Address{101}, &trader).Return(big.NewInt(-15))
		mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, order.AmmIndex).Return(big.NewInt(10))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order:  order,
			Sender: common.Address{1},
			Error:  ErrNetReduceOnlyAmountExceeded,
		})
	})

	t.Run("invalid price - negative price", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		order := IImmediateOrCancelOrdersOrder{
			OrderType:         1,
			ExpireAt:          big.NewInt(1004),
			AmmIndex:          big.NewInt(0),
			Trader:            trader,
			BaseAssetQuantity: big.NewInt(10),
			Price:             big.NewInt(-100),
			Salt:              big.NewInt(12),
			ReduceOnly:        true,
		}
		hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(0))
		mockBibliophile.EXPECT().HasReferrer(trader).Return(true)
		mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(common.Address{101})
		mockBibliophile.EXPECT().GetSize(common.Address{101}, &trader).Return(big.NewInt(-15))
		mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, order.AmmIndex).Return(big.NewInt(0))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order:  order,
			Sender: common.Address{1},
			Error:  ErrInvalidPrice,
		})
	})

	t.Run("invalid price - price not multiple of price multiplier", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)
		mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
		mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

		order := IImmediateOrCancelOrdersOrder{
			OrderType:         1,
			ExpireAt:          big.NewInt(1004),
			AmmIndex:          big.NewInt(0),
			Trader:            trader,
			BaseAssetQuantity: big.NewInt(10),
			Price:             big.NewInt(101),
			Salt:              big.NewInt(13),
			ReduceOnly:        true,
		}
		hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(0))
		mockBibliophile.EXPECT().HasReferrer(trader).Return(true)
		mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(common.Address{101})
		mockBibliophile.EXPECT().GetSize(common.Address{101}, &trader).Return(big.NewInt(-15))
		mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, order.AmmIndex).Return(big.NewInt(0))
		mockBibliophile.EXPECT().GetPriceMultiplier(common.Address{101}).Return(big.NewInt(10))

		testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
			Order:  order,
			Sender: common.Address{1},
			Error:  ErrPricePrecision,
		})

		t.Run("valid order", func(t *testing.T) {
			mockBibliophile := b.NewMockBibliophileClient(ctrl)
			mockBibliophile.EXPECT().IsTradingAuthority(trader, common.Address{1}).Return(true)
			mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
			mockBibliophile.EXPECT().IOC_GetExpirationCap().Return(big.NewInt(5))
			mockBibliophile.EXPECT().GetMinSizeRequirement(int64(0)).Return(big.NewInt(5))

			order := IImmediateOrCancelOrdersOrder{
				OrderType:         1,
				ExpireAt:          big.NewInt(1004),
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(10),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(13),
				ReduceOnly:        true,
			}
			hash, _ := IImmediateOrCancelOrdersOrderToIOCOrder(&order).Hash()

			mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(0))
			mockBibliophile.EXPECT().HasReferrer(trader).Return(true)
			mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(common.Address{101})
			mockBibliophile.EXPECT().GetSize(common.Address{101}, &trader).Return(big.NewInt(-15))
			mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, order.AmmIndex).Return(big.NewInt(0))
			mockBibliophile.EXPECT().GetPriceMultiplier(common.Address{101}).Return(big.NewInt(10))

			testValidatePlaceIOCOrderTestCase(t, mockBibliophile, ValidatePlaceIOCOrderTestCase{
				Order:  order,
				Sender: common.Address{1},
				Error:  nil,
			})
		})
	})
}

func TestValidateExecuteIOCOrder(t *testing.T) {
	trader := common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("not ioc order", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)

		order := hu.IOCOrder{
			OrderType: 0, // incoreect order type
			ExpireAt:  big.NewInt(1001),
			BaseOrder: hu.BaseOrder{
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(10),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(1),
				ReduceOnly:        false,
			},
		}
		m, err := validateExecuteIOCOrder(mockBibliophile, &order, Long, big.NewInt(10))
		assert.EqualError(t, err, "not ioc order")
		hash, _ := order.Hash()
		assert.Equal(t, m.OrderHash, hash)
	})

	t.Run("ioc expired", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)

		order := hu.IOCOrder{
			OrderType: 1,
			ExpireAt:  big.NewInt(990),
			BaseOrder: hu.BaseOrder{
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(10),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(1),
				ReduceOnly:        false,
			},
		}
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))

		m, err := validateExecuteIOCOrder(mockBibliophile, &order, Long, big.NewInt(10))
		assert.EqualError(t, err, "ioc expired")
		hash, _ := order.Hash()
		assert.Equal(t, m.OrderHash, hash)
	})

	t.Run("valid order", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)

		order := hu.IOCOrder{
			OrderType: 1,
			ExpireAt:  big.NewInt(1001),
			BaseOrder: hu.BaseOrder{
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(10),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(1),
				ReduceOnly:        false,
			},
		}
		hash, _ := order.Hash()
		ammAddress := common.Address{101}
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(ammAddress)
		mockBibliophile.EXPECT().IOC_GetOrderFilledAmount(hash).Return(big.NewInt(0))
		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(1))
		mockBibliophile.EXPECT().IOC_GetBlockPlaced(hash).Return(big.NewInt(21))

		m, err := validateExecuteIOCOrder(mockBibliophile, &order, Long, big.NewInt(10))
		assert.Nil(t, err)
		assertMetadataEquality(t, &Metadata{
			AmmIndex:          new(big.Int).Set(order.AmmIndex),
			Trader:            trader,
			BaseAssetQuantity: new(big.Int).Set(order.BaseAssetQuantity),
			BlockPlaced:       big.NewInt(21),
			Price:             new(big.Int).Set(order.Price),
			OrderHash:         hash,
		}, m)
	})

	t.Run("valid order - reduce only", func(t *testing.T) {
		mockBibliophile := b.NewMockBibliophileClient(ctrl)

		order := hu.IOCOrder{
			OrderType: 1,
			ExpireAt:  big.NewInt(1001),
			BaseOrder: hu.BaseOrder{
				AmmIndex:          big.NewInt(0),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(10),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(1),
				ReduceOnly:        true,
			},
		}
		hash, _ := order.Hash()
		ammAddress := common.Address{101}
		mockBibliophile.EXPECT().GetTimeStamp().Return(uint64(1000))
		mockBibliophile.EXPECT().GetMarketAddressFromMarketID(int64(0)).Return(ammAddress)
		mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(big.NewInt(-10))
		mockBibliophile.EXPECT().IOC_GetOrderFilledAmount(hash).Return(big.NewInt(0))
		mockBibliophile.EXPECT().IOC_GetOrderStatus(hash).Return(int64(1))
		mockBibliophile.EXPECT().IOC_GetBlockPlaced(hash).Return(big.NewInt(21))

		_, err := validateExecuteIOCOrder(mockBibliophile, &order, Long, big.NewInt(10))
		assert.Nil(t, err)
	})
}
