package limitorders

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var positionType = "short"
var userAddress = "random-address"
var baseAssetQuantity = big.NewInt(-10)
var price = big.NewInt(20)
var status Status = Placed
var blockNumber = big.NewInt(2)

func TestNewInMemoryDatabase(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	assert.NotNil(t, inMemoryDatabase)
}

func TestAdd(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	signature := []byte("Here is a string....")
	id := uint64(123)
	salt := big.NewInt(time.Now().Unix())
	limitOrder := createLimitOrder(id, positionType, userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
	inMemoryDatabase.Add(&limitOrder)
	returnedOrder := inMemoryDatabase.GetAllOrders()[0]
	assert.Equal(t, id, returnedOrder.Id)
	assert.Equal(t, limitOrder.PositionType, returnedOrder.PositionType)
	assert.Equal(t, limitOrder.UserAddress, returnedOrder.UserAddress)
	assert.Equal(t, limitOrder.BaseAssetQuantity, returnedOrder.BaseAssetQuantity)
	assert.Equal(t, limitOrder.Price, returnedOrder.Price)
	assert.Equal(t, limitOrder.Status, returnedOrder.Status)
	assert.Equal(t, limitOrder.BlockNumber, returnedOrder.BlockNumber)
}

func TestGetAllOrders(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	totalOrders := uint64(5)
	for i := uint64(0); i < totalOrders; i++ {
		signature := []byte("signature")
		salt := big.NewInt(0).Add(big.NewInt(int64(i)), big.NewInt(time.Now().Unix()))
		limitOrder := createLimitOrder(i, positionType, userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
		inMemoryDatabase.Add(&limitOrder)
	}
	returnedOrders := inMemoryDatabase.GetAllOrders()
	assert.Equal(t, totalOrders, uint64(len(returnedOrders)))
	for _, returnedOrder := range returnedOrders {
		assert.Equal(t, positionType, returnedOrder.PositionType)
		assert.Equal(t, userAddress, returnedOrder.UserAddress)
		assert.Equal(t, baseAssetQuantity, returnedOrder.BaseAssetQuantity)
		assert.Equal(t, price, returnedOrder.Price)
		assert.Equal(t, status, returnedOrder.Status)
		assert.Equal(t, blockNumber, returnedOrder.BlockNumber)
	}
}

func TestGetShortOrders(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	totalLongOrders := uint64(2)
	longOrderPrice := big.NewInt(0).Add(price, big.NewInt(1))
	longOrderBaseAssetQuantity := big.NewInt(10)
	for i := uint64(0); i < totalLongOrders; i++ {
		signature := []byte("signature")
		salt := big.NewInt(0).Add(big.NewInt(int64(i)), big.NewInt(time.Now().Unix()))
		limitOrder := createLimitOrder(i, "long", userAddress, longOrderBaseAssetQuantity, longOrderPrice, status, signature, blockNumber, salt)
		inMemoryDatabase.Add(&limitOrder)
	}
	//Short order with price 10 and blockNumber 2
	id1 := uint64(1)
	signature1 := []byte(fmt.Sprintf("Signature short order is %d", id1))
	price1 := big.NewInt(10)
	blockNumber1 := big.NewInt(2)
	salt1 := big.NewInt(time.Now().Unix())
	shortOrder1 := createLimitOrder(id1, "short", userAddress, baseAssetQuantity, price1, status, signature1, blockNumber1, salt1)
	inMemoryDatabase.Add(&shortOrder1)

	//Short order with price 9 and blockNumber 2
	id2 := uint64(2)
	signature2 := []byte(fmt.Sprintf("Signature short order is %d", id2))
	price2 := big.NewInt(9)
	blockNumber2 := big.NewInt(2)
	salt2 := big.NewInt(0).Add(salt1, big.NewInt(1))
	shortOrder2 := createLimitOrder(id2, "short", userAddress, baseAssetQuantity, price2, status, signature2, blockNumber2, salt2)
	inMemoryDatabase.Add(&shortOrder2)

	//Short order with price 9.01 and blockNumber 3
	id3 := uint64(3)
	signature3 := []byte(fmt.Sprintf("Signature short order is %d", id3))
	price3 := big.NewInt(9)
	blockNumber3 := big.NewInt(3)
	salt3 := big.NewInt(0).Add(salt2, big.NewInt(1))
	shortOrder3 := createLimitOrder(id3, "short", userAddress, baseAssetQuantity, price3, status, signature3, blockNumber3, salt3)
	inMemoryDatabase.Add(&shortOrder3)

	returnedShortOrders := inMemoryDatabase.GetShortOrders(AvaxPerp)
	assert.Equal(t, 3, len(returnedShortOrders))

	for _, returnedOrder := range returnedShortOrders {
		assert.Equal(t, "short", returnedOrder.PositionType)
		assert.Equal(t, userAddress, returnedOrder.UserAddress)
		assert.Equal(t, baseAssetQuantity, returnedOrder.BaseAssetQuantity)
		assert.Equal(t, status, returnedOrder.Status)
	}

	//Test returnedShortOrders are sorted by price lowest to highest first and then block number from lowest to highest
	assert.Equal(t, id2, returnedShortOrders[0].Id)
	assert.Equal(t, price2, returnedShortOrders[0].Price)
	assert.Equal(t, blockNumber2, returnedShortOrders[0].BlockNumber)
	assert.Equal(t, id3, returnedShortOrders[1].Id)
	assert.Equal(t, price3, returnedShortOrders[1].Price)
	assert.Equal(t, blockNumber3, returnedShortOrders[1].BlockNumber)
	assert.Equal(t, id1, returnedShortOrders[2].Id)
	assert.Equal(t, price1, returnedShortOrders[2].Price)
	assert.Equal(t, blockNumber1, returnedShortOrders[2].BlockNumber)

}

func TestGetLongOrders(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	for i := uint64(0); i < 3; i++ {
		signature := []byte("signature")
		salt := big.NewInt(0).Add(big.NewInt(time.Now().Unix()), big.NewInt(int64(i)))
		limitOrder := createLimitOrder(i, "short", userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
		inMemoryDatabase.Add(&limitOrder)
	}

	//Long order with price 9 and blockNumber 2
	longOrderBaseAssetQuantity := big.NewInt(10)
	id1 := uint64(1)
	signature1 := []byte(fmt.Sprintf("Signature long order is %d", id1))
	price1 := big.NewInt(9)
	blockNumber1 := big.NewInt(2)
	salt1 := big.NewInt(time.Now().Unix())
	longOrder1 := createLimitOrder(id1, "long", userAddress, longOrderBaseAssetQuantity, price1, status, signature1, blockNumber1, salt1)
	inMemoryDatabase.Add(&longOrder1)

	//long order with price 9 and blockNumber 3
	id2 := uint64(2)
	signature2 := []byte(fmt.Sprintf("Signature long order is %d", id2))
	price2 := big.NewInt(9)
	blockNumber2 := big.NewInt(3)
	salt2 := big.NewInt(0).Add(salt1, big.NewInt(1))
	longOrder2 := createLimitOrder(id2, "long", userAddress, longOrderBaseAssetQuantity, price2, status, signature2, blockNumber2, salt2)
	inMemoryDatabase.Add(&longOrder2)

	//long order with price 10 and blockNumber 3
	id3 := uint64(3)
	signature3 := []byte(fmt.Sprintf("Signature long order is %d", id3))
	price3 := big.NewInt(10)
	blockNumber3 := big.NewInt(3)
	salt3 := big.NewInt(0).Add(salt2, big.NewInt(1))
	longOrder3 := createLimitOrder(id3, "long", userAddress, longOrderBaseAssetQuantity, price3, status, signature3, blockNumber3, salt3)
	inMemoryDatabase.Add(&longOrder3)

	returnedLongOrders := inMemoryDatabase.GetLongOrders(AvaxPerp)
	assert.Equal(t, 3, len(returnedLongOrders))

	//Test returnedLongOrders are sorted by price highest to lowest first and then block number from lowest to highest
	assert.Equal(t, id3, returnedLongOrders[0].Id)
	assert.Equal(t, price3, returnedLongOrders[0].Price)
	assert.Equal(t, blockNumber3, returnedLongOrders[0].BlockNumber)
	assert.Equal(t, id1, returnedLongOrders[1].Id)
	assert.Equal(t, price1, returnedLongOrders[1].Price)
	assert.Equal(t, blockNumber1, returnedLongOrders[1].BlockNumber)
	assert.Equal(t, id2, returnedLongOrders[2].Id)
	assert.Equal(t, price2, returnedLongOrders[2].Price)
	assert.Equal(t, blockNumber2, returnedLongOrders[2].BlockNumber)

	for _, returnedOrder := range returnedLongOrders {
		assert.Equal(t, "long", returnedOrder.PositionType)
		assert.Equal(t, userAddress, returnedOrder.UserAddress)
		assert.Equal(t, longOrderBaseAssetQuantity, returnedOrder.BaseAssetQuantity)
		assert.Equal(t, status, returnedOrder.Status)
	}
}

func TestUpdateFulfilledBaseAssetQuantityLimitOrder(t *testing.T) {
	t.Run("When order does not exists", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		signature := []byte("Here is a string....")
		id := uint64(123)
		salt := big.NewInt(time.Now().Unix())
		limitOrder := createLimitOrder(id, positionType, userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
		filledQuantity := big.NewInt(2)

		inMemoryDatabase.UpdateFilledBaseAssetQuantity(filledQuantity, getIdFromLimitOrder(limitOrder))
		updatedLimitOrder := inMemoryDatabase.OrderMap[getIdFromLimitOrder(limitOrder)]

		assert.Nil(t, updatedLimitOrder)

	})
	t.Run("when filled quantity is not equal to baseAssetQuantity", func(t *testing.T) {
		t.Run("When order type is short order", func(t *testing.T) {
			inMemoryDatabase := NewInMemoryDatabase()
			signature := []byte("Here is a string....")
			id := uint64(123)
			salt := big.NewInt(time.Now().Unix())
			limitOrder := createLimitOrder(id, positionType, userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
			inMemoryDatabase.Add(&limitOrder)

			filledQuantity := big.NewInt(2)

			inMemoryDatabase.UpdateFilledBaseAssetQuantity(filledQuantity, getIdFromLimitOrder(limitOrder))
			updatedLimitOrder := inMemoryDatabase.OrderMap[getIdFromLimitOrder(limitOrder)]

			assert.Equal(t, updatedLimitOrder.FilledBaseAssetQuantity, big.NewInt(0).Neg(filledQuantity))
			assert.Equal(t, updatedLimitOrder.FilledBaseAssetQuantity, filledQuantity.Mul(filledQuantity, big.NewInt(-1)))
		})
		t.Run("When order type is long order", func(t *testing.T) {
			inMemoryDatabase := NewInMemoryDatabase()
			signature := []byte("Here is a string....")
			id := uint64(123)
			positionType = "long"
			baseAssetQuantity = big.NewInt(10)
			salt := big.NewInt(time.Now().Unix())
			limitOrder := createLimitOrder(id, positionType, userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
			inMemoryDatabase.Add(&limitOrder)

			filledQuantity := big.NewInt(2)
			inMemoryDatabase.UpdateFilledBaseAssetQuantity(filledQuantity, getIdFromLimitOrder(limitOrder))
			updatedLimitOrder := inMemoryDatabase.OrderMap[getIdFromLimitOrder(limitOrder)]

			assert.Equal(t, updatedLimitOrder.FilledBaseAssetQuantity, filledQuantity)
		})
	})
	t.Run("when filled quantity is equal to baseAssetQuantity", func(t *testing.T) {
		t.Run("When order type is short order", func(t *testing.T) {
			inMemoryDatabase := NewInMemoryDatabase()
			signature := []byte("Here is a string....")
			id := uint64(123)
			salt := big.NewInt(time.Now().Unix())
			limitOrder := createLimitOrder(id, positionType, userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
			inMemoryDatabase.Add(&limitOrder)

			filledQuantity := big.NewInt(0).Abs(limitOrder.BaseAssetQuantity)
			inMemoryDatabase.UpdateFilledBaseAssetQuantity(filledQuantity, getIdFromLimitOrder(limitOrder))
			allOrders := inMemoryDatabase.GetAllOrders()

			assert.Equal(t, 0, len(allOrders))
		})
		t.Run("When order type is long order", func(t *testing.T) {
			inMemoryDatabase := NewInMemoryDatabase()
			signature := []byte("Here is a string....")
			id := uint64(123)
			positionType = "long"
			baseAssetQuantity = big.NewInt(10)
			salt := big.NewInt(time.Now().Unix())
			limitOrder := createLimitOrder(id, positionType, userAddress, baseAssetQuantity, price, status, signature, blockNumber, salt)
			inMemoryDatabase.Add(&limitOrder)

			filledQuantity := big.NewInt(0).Abs(limitOrder.BaseAssetQuantity)
			inMemoryDatabase.UpdateFilledBaseAssetQuantity(filledQuantity, getIdFromLimitOrder(limitOrder))
			allOrders := inMemoryDatabase.GetAllOrders()

			assert.Equal(t, 0, len(allOrders))
		})
	})
}

func TestUpdatePosition(t *testing.T) {
	t.Run("When no positions exists for trader, it updates trader map with new positions", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var market Market = 1
		size := big.NewInt(20.00)
		openNotional := big.NewInt(200.00)
		inMemoryDatabase.UpdatePosition(address, market, size, openNotional, false)
		position := inMemoryDatabase.TraderMap[address].Positions[market]
		assert.Equal(t, size, position.Size)
		assert.Equal(t, openNotional, position.OpenNotional)
	})
	t.Run("When positions exists for trader, it overwrites old positions with new data", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var market Market = 1
		size := big.NewInt(20.00)
		openNotional := big.NewInt(200.00)
		inMemoryDatabase.UpdatePosition(address, market, size, openNotional, false)

		newSize := big.NewInt(25.00)
		newOpenNotional := big.NewInt(250.00)
		inMemoryDatabase.UpdatePosition(address, market, newSize, newOpenNotional, false)
		position := inMemoryDatabase.TraderMap[address].Positions[market]
		assert.Equal(t, newSize, position.Size)
		assert.Equal(t, newOpenNotional, position.OpenNotional)
	})
}

func TestUpdateMargin(t *testing.T) {
	t.Run("when adding margin for first time it updates margin in tradermap", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var collateral Collateral = 1
		amount := big.NewInt(20.00)
		inMemoryDatabase.UpdateMargin(address, collateral, amount)
		margin := inMemoryDatabase.TraderMap[address].Margins[collateral]
		assert.Equal(t, amount, margin)
	})
	t.Run("When more margin is added, it updates margin in tradermap", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var collateral Collateral = 1
		amount := big.NewInt(20.00)
		inMemoryDatabase.UpdateMargin(address, collateral, amount)

		removedMargin := big.NewInt(15.00)
		inMemoryDatabase.UpdateMargin(address, collateral, removedMargin)
		margin := inMemoryDatabase.TraderMap[address].Margins[collateral]
		assert.Equal(t, big.NewInt(0).Add(amount, removedMargin), margin)
	})
	t.Run("When margin is removed, it updates margin in tradermap", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var collateral Collateral = 1
		amount := big.NewInt(20.00)
		inMemoryDatabase.UpdateMargin(address, collateral, amount)

		removedMargin := big.NewInt(-15.00)
		inMemoryDatabase.UpdateMargin(address, collateral, removedMargin)
		margin := inMemoryDatabase.TraderMap[address].Margins[collateral]
		assert.Equal(t, big.NewInt(0).Add(amount, removedMargin), margin)
	})
}

func TestUpdateUnrealizedFunding(t *testing.T) {
	t.Run("When trader has no positions, it does not update anything", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var market Market = 1
		cumulativePremiumFraction := big.NewInt(2)
		trader := inMemoryDatabase.TraderMap[address]
		inMemoryDatabase.UpdateUnrealisedFunding(market, cumulativePremiumFraction)
		updatedTrader := inMemoryDatabase.TraderMap[address]
		assert.Equal(t, trader, updatedTrader)
	})
	t.Run("When trader has positions", func(t *testing.T) {
		t.Run("when unrealized funding is zero, it updates unrealized funding in trader's positions", func(t *testing.T) {
			inMemoryDatabase := NewInMemoryDatabase()
			addresses := [2]common.Address{common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"), common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")}
			var market Market = 1
			openNotional := big.NewInt(200.00)
			cumulativePremiumFraction := big.NewInt(0)
			for i, address := range addresses {
				iterator := i + 1
				size := big.NewInt(int64(20 * iterator))
				inMemoryDatabase.UpdatePosition(address, market, size, openNotional, false)
				inMemoryDatabase.ResetUnrealisedFunding(market, address, cumulativePremiumFraction)
			}
			newCumulativePremiumFraction := big.NewInt(5)
			inMemoryDatabase.UpdateUnrealisedFunding(market, newCumulativePremiumFraction)
			for _, address := range addresses {
				unrealizedFunding := inMemoryDatabase.TraderMap[address].Positions[market].UnrealisedFunding
				size := inMemoryDatabase.TraderMap[address].Positions[market].Size
				expectedUnrealizedFunding := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Sub(newCumulativePremiumFraction, cumulativePremiumFraction), size), SIZE_BASE_PRECISION)
				assert.Equal(t, expectedUnrealizedFunding, unrealizedFunding)
			}
		})
		t.Run("when unrealized funding is not zero, it adds new funding to old unrealized funding in trader's positions", func(t *testing.T) {
			inMemoryDatabase := NewInMemoryDatabase()
			address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
			var market Market = 1
			openNotional := big.NewInt(200.00)
			size := big.NewInt(20.00)
			inMemoryDatabase.UpdatePosition(address, market, size, openNotional, false)
			cumulativePremiumFraction := big.NewInt(2)
			inMemoryDatabase.ResetUnrealisedFunding(market, address, cumulativePremiumFraction)

			newCumulativePremiumFraction := big.NewInt(-1)
			inMemoryDatabase.UpdateUnrealisedFunding(market, newCumulativePremiumFraction)
			newUnrealizedFunding := inMemoryDatabase.TraderMap[address].Positions[market].UnrealisedFunding
			expectedUnrealizedFunding := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Sub(newCumulativePremiumFraction, cumulativePremiumFraction), size), SIZE_BASE_PRECISION)
			assert.Equal(t, expectedUnrealizedFunding, newUnrealizedFunding)
		})
	})
}

func TestResetUnrealisedFunding(t *testing.T) {
	t.Run("When trader has no positions, it does not update anything", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var market Market = 1
		trader := inMemoryDatabase.TraderMap[address]
		cumulativePremiumFraction := big.NewInt(5)
		inMemoryDatabase.ResetUnrealisedFunding(market, address, cumulativePremiumFraction)
		updatedTrader := inMemoryDatabase.TraderMap[address]
		assert.Equal(t, trader, updatedTrader)
	})
	t.Run("When trader has positions, it resets unrealized funding to zero", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		address := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		var market Market = 1
		openNotional := big.NewInt(200)
		size := big.NewInt(20)
		inMemoryDatabase.UpdatePosition(address, market, size, openNotional, false)
		cumulativePremiumFraction := big.NewInt(1)
		inMemoryDatabase.ResetUnrealisedFunding(market, address, cumulativePremiumFraction)
		unrealizedFundingFee := inMemoryDatabase.TraderMap[address].Positions[market].UnrealisedFunding
		assert.Equal(t, big.NewInt(0), unrealizedFundingFee)
	})
}

func TestUpdateNextFundingTime(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	nextFundingTime := uint64(time.Now().Unix())
	inMemoryDatabase.UpdateNextFundingTime(nextFundingTime)
	assert.Equal(t, nextFundingTime, inMemoryDatabase.NextFundingTime)
}

func TestGetNextFundingTime(t *testing.T) {
	t.Run("when funding time is not set", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		assert.Equal(t, uint64(0), inMemoryDatabase.GetNextFundingTime())
	})
	t.Run("when funding time is set", func(t *testing.T) {
		inMemoryDatabase := NewInMemoryDatabase()
		nextFundingTime := uint64(time.Now().Unix())
		inMemoryDatabase.UpdateNextFundingTime(nextFundingTime)
		assert.Equal(t, nextFundingTime, inMemoryDatabase.GetNextFundingTime())
	})
}

func TestUpdateLastPrice(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	var market Market = 1
	lastPrice := big.NewInt(20)
	inMemoryDatabase.UpdateLastPrice(market, lastPrice)
	assert.Equal(t, lastPrice, inMemoryDatabase.LastPrice[market])
}
func TestGetLastPrice(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	var market Market = 1
	lastPrice := big.NewInt(20)
	inMemoryDatabase.UpdateLastPrice(market, lastPrice)
	assert.Equal(t, lastPrice, inMemoryDatabase.GetLastPrice(market))
}

func createLimitOrder(id uint64, positionType string, userAddress string, baseAssetQuantity *big.Int, price *big.Int, status Status, signature []byte, blockNumber *big.Int, salt *big.Int) LimitOrder {
	return LimitOrder{
		Id:                      id,
		PositionType:            positionType,
		UserAddress:             userAddress,
		FilledBaseAssetQuantity: big.NewInt(0),
		BaseAssetQuantity:       baseAssetQuantity,
		Price:                   price,
		Status:                  Status(status),
		Salt:                    salt,
		Signature:               signature,
		BlockNumber:             blockNumber,
	}
}

func TestGetUnfilledBaseAssetQuantity(t *testing.T) {
	t.Run("When limit FilledBaseAssetQuantity is zero, it returns BaseAssetQuantity", func(t *testing.T) {
		baseAssetQuantityLongOrder := big.NewInt(10)
		signature := []byte("Here is a long order")
		salt1 := big.NewInt(time.Now().Unix())
		longOrder := createLimitOrder(uint64(1), "long", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", baseAssetQuantityLongOrder, big.NewInt(21), "unfulfilled", signature, big.NewInt(2), salt1)
		longOrder.FilledBaseAssetQuantity = big.NewInt(0)
		//baseAssetQuantityLongOrder - filledBaseAssetQuantity
		expectedUnFilledForLongOrder := big.NewInt(10)
		assert.Equal(t, expectedUnFilledForLongOrder, longOrder.GetUnFilledBaseAssetQuantity())

		signature = []byte("Here is a short order")
		baseAssetQuantityShortOrder := big.NewInt(-10)
		salt2 := big.NewInt(0).Add(salt1, big.NewInt(1))
		shortOrder := createLimitOrder(uint64(1), "short", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", baseAssetQuantityShortOrder, big.NewInt(21), "unfulfilled", signature, big.NewInt(2), salt2)
		shortOrder.FilledBaseAssetQuantity = big.NewInt(0)
		//baseAssetQuantityLongOrder - filledBaseAssetQuantity
		expectedUnFilledForShortOrder := big.NewInt(-10)
		assert.Equal(t, expectedUnFilledForShortOrder, shortOrder.GetUnFilledBaseAssetQuantity())
	})
	t.Run("When limit FilledBaseAssetQuantity is not zero, it returns BaseAssetQuantity - FilledBaseAssetQuantity", func(t *testing.T) {
		baseAssetQuantityLongOrder := big.NewInt(10)
		signature := []byte("Here is a long order")
		salt1 := big.NewInt(time.Now().Unix())
		longOrder := createLimitOrder(uint64(1), "long", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", baseAssetQuantityLongOrder, big.NewInt(21), "unfulfilled", signature, big.NewInt(2), salt1)
		longOrder.FilledBaseAssetQuantity = big.NewInt(5)
		//baseAssetQuantityLongOrder - filledBaseAssetQuantity
		expectedUnFilledForLongOrder := big.NewInt(5)
		assert.Equal(t, expectedUnFilledForLongOrder, longOrder.GetUnFilledBaseAssetQuantity())

		signature = []byte("Here is a short order")
		baseAssetQuantityShortOrder := big.NewInt(-10)
		salt2 := big.NewInt(0).Add(salt1, big.NewInt(1))
		shortOrder := createLimitOrder(uint64(1), "short", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", baseAssetQuantityShortOrder, big.NewInt(21), "unfulfilled", signature, big.NewInt(2), salt2)
		shortOrder.FilledBaseAssetQuantity = big.NewInt(-5)
		//baseAssetQuantityLongOrder - filledBaseAssetQuantity
		expectedUnFilledForShortOrder := big.NewInt(-5)
		assert.Equal(t, expectedUnFilledForShortOrder, shortOrder.GetUnFilledBaseAssetQuantity())
	})
}

func getOrderFromLimitOrder(limitOrder LimitOrder) Order {
	return Order{
		Trader:            common.HexToAddress(limitOrder.UserAddress),
		AmmIndex:          big.NewInt(0),
		BaseAssetQuantity: limitOrder.BaseAssetQuantity,
		Price:             limitOrder.Price,
		Salt:              limitOrder.Salt,
	}
}
