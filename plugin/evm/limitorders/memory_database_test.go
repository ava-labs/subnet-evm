package limitorders

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var positionType = "short"
var userAddress = "random-address"
var baseAssetQuantity = -10
var price float64 = 20.01
var status = "unfulfilled"
var salt = "salt"
var blockNumber uint64 = 2

func TestNewInMemoryDatabase(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	assert.NotNil(t, inMemoryDatabase)
}

func TestAdd(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	signature := []byte("Here is a string....")
	id := uint64(123)
	limitOrder := createLimitOrder(id, positionType, userAddress, baseAssetQuantity, price, status, salt, signature, blockNumber)
	inMemoryDatabase.Add(&limitOrder)
	returnedOrder := inMemoryDatabase.GetAllOrders()[0]
	assert.Equal(t, id, returnedOrder.id)
	assert.Equal(t, limitOrder.PositionType, returnedOrder.PositionType)
	assert.Equal(t, limitOrder.UserAddress, returnedOrder.UserAddress)
	assert.Equal(t, limitOrder.BaseAssetQuantity, returnedOrder.BaseAssetQuantity)
	assert.Equal(t, limitOrder.Price, returnedOrder.Price)
	assert.Equal(t, limitOrder.Status, returnedOrder.Status)
	assert.Equal(t, limitOrder.Salt, returnedOrder.Salt)
	assert.Equal(t, limitOrder.BlockNumber, returnedOrder.BlockNumber)
}

func TestGetAllOrders(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	totalOrders := uint64(5)
	for i := uint64(0); i < totalOrders; i++ {
		signature := []byte(fmt.Sprintf("Signature is %d", i))
		limitOrder := createLimitOrder(i, positionType, userAddress, baseAssetQuantity, price, status, salt, signature, blockNumber)
		inMemoryDatabase.Add(&limitOrder)
	}
	returnedOrders := inMemoryDatabase.GetAllOrders()
	assert.Equal(t, totalOrders, uint64(len(returnedOrders)))
	fmt.Println(returnedOrders)
	for _, returedOrder := range returnedOrders {
		assert.Equal(t, positionType, returedOrder.PositionType)
		assert.Equal(t, userAddress, returedOrder.UserAddress)
		assert.Equal(t, baseAssetQuantity, returedOrder.BaseAssetQuantity)
		assert.Equal(t, price, returedOrder.Price)
		assert.Equal(t, status, returedOrder.Status)
		assert.Equal(t, salt, returedOrder.Salt)
		assert.Equal(t, blockNumber, returedOrder.BlockNumber)
	}
}

func TestDelete(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	totalOrders := uint64(5)
	for i := uint64(0); i < totalOrders; i++ {
		signature := []byte(fmt.Sprintf("Signature is %d", i))
		limitOrder := createLimitOrder(i, positionType, userAddress, baseAssetQuantity, price, status, salt, signature, blockNumber)
		inMemoryDatabase.Add(&limitOrder)
	}

	deletedOrderId := 3
	inMemoryDatabase.Delete([]byte(fmt.Sprintf("Signature is %d", deletedOrderId)))
	expectedReturnedOrdersIds := []int{0, 1, 2, 4}

	returnedOrders := inMemoryDatabase.GetAllOrders()
	assert.Equal(t, totalOrders-1, uint64(len(returnedOrders)))
	var returnedOrderIds []int
	for _, returnedOrder := range returnedOrders {
		returnedOrderIds = append(returnedOrderIds, int(returnedOrder.id))
	}
	sort.Ints(returnedOrderIds)
	assert.Equal(t, expectedReturnedOrdersIds, returnedOrderIds)
}

func TestGetShortOrders(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	totalLongOrders := uint64(2)
	longOrderPrice := price + 1
	longOrderBaseAssetQuantity := 10
	for i := uint64(0); i < totalLongOrders; i++ {
		signature := []byte(fmt.Sprintf("Signature long order is %d", i))
		limitOrder := createLimitOrder(i, "long", userAddress, longOrderBaseAssetQuantity, longOrderPrice, status, salt, signature, blockNumber)
		inMemoryDatabase.Add(&limitOrder)
	}
	//Short order with price 10.01 and blockNumber 2
	id1 := uint64(1)
	signature1 := []byte(fmt.Sprintf("Signature short order is %d", id1))
	price1 := 10.01
	var blockNumber1 uint64 = 2
	shortOrder1 := createLimitOrder(id1, "short", userAddress, baseAssetQuantity, price1, status, salt, signature1, blockNumber1)
	inMemoryDatabase.Add(&shortOrder1)

	//Short order with price 9.01 and blockNumber 2
	id2 := uint64(2)
	signature2 := []byte(fmt.Sprintf("Signature short order is %d", id2))
	price2 := 9.01
	var blockNumber2 uint64 = 2
	shortOrder2 := createLimitOrder(id2, "short", userAddress, baseAssetQuantity, price2, status, salt, signature2, blockNumber2)
	inMemoryDatabase.Add(&shortOrder2)

	//Short order with price 9.01 and blockNumber 3
	id3 := uint64(3)
	signature3 := []byte(fmt.Sprintf("Signature short order is %d", id3))
	price3 := 9.01
	var blockNumber3 uint64 = 3
	shortOrder3 := createLimitOrder(id3, "short", userAddress, baseAssetQuantity, price3, status, salt, signature3, blockNumber3)
	inMemoryDatabase.Add(&shortOrder3)

	returnedShortOrders := inMemoryDatabase.GetShortOrders()
	assert.Equal(t, 3, len(returnedShortOrders))

	for _, returedOrder := range returnedShortOrders {
		assert.Equal(t, "short", returedOrder.PositionType)
		assert.Equal(t, userAddress, returedOrder.UserAddress)
		assert.Equal(t, baseAssetQuantity, returedOrder.BaseAssetQuantity)
		assert.Equal(t, status, returedOrder.Status)
		assert.Equal(t, salt, returedOrder.Salt)
	}

	//Test returnedShortOrders are sorted by price lowest to highest first and then block number from lowest to highest
	assert.Equal(t, id2, returnedShortOrders[0].id)
	assert.Equal(t, price2, returnedShortOrders[0].Price)
	assert.Equal(t, blockNumber2, returnedShortOrders[0].BlockNumber)
	assert.Equal(t, id3, returnedShortOrders[1].id)
	assert.Equal(t, price3, returnedShortOrders[1].Price)
	assert.Equal(t, blockNumber3, returnedShortOrders[1].BlockNumber)
	assert.Equal(t, id1, returnedShortOrders[2].id)
	assert.Equal(t, price1, returnedShortOrders[2].Price)
	assert.Equal(t, blockNumber1, returnedShortOrders[2].BlockNumber)

}

func TestGetLongOrders(t *testing.T) {
	inMemoryDatabase := NewInMemoryDatabase()
	for i := uint64(0); i < 3; i++ {
		signature := []byte(fmt.Sprintf("Signature short order is %d", i))
		limitOrder := createLimitOrder(i, "short", userAddress, baseAssetQuantity, price, status, salt, signature, blockNumber)
		inMemoryDatabase.Add(&limitOrder)
	}

	//Long order with price 9.01 and blockNumber 2
	longOrderBaseAssetQuantity := 10
	id1 := uint64(1)
	signature1 := []byte(fmt.Sprintf("Signature long order is %d", id1))
	price1 := 9.01
	var blockNumber1 uint64 = 2
	longOrder1 := createLimitOrder(id1, "long", userAddress, longOrderBaseAssetQuantity, price1, status, salt, signature1, blockNumber1)
	inMemoryDatabase.Add(&longOrder1)

	//long order with price 9.01 and blockNumber 3
	id2 := uint64(2)
	signature2 := []byte(fmt.Sprintf("Signature long order is %d", id2))
	price2 := 9.01
	var blockNumber2 uint64 = 3
	longOrder2 := createLimitOrder(id2, "long", userAddress, longOrderBaseAssetQuantity, price2, status, salt, signature2, blockNumber2)
	inMemoryDatabase.Add(&longOrder2)

	//long order with price 10.01 and blockNumber 3
	id3 := uint64(3)
	signature3 := []byte(fmt.Sprintf("Signature long order is %d", id3))
	price3 := 10.01
	var blockNumber3 uint64 = 3
	longOrder3 := createLimitOrder(id3, "long", userAddress, longOrderBaseAssetQuantity, price3, status, salt, signature3, blockNumber3)
	inMemoryDatabase.Add(&longOrder3)

	returnedLongOrders := inMemoryDatabase.GetLongOrders()
	assert.Equal(t, 3, len(returnedLongOrders))

	//Test returnedLongOrders are sorted by price highest to lowest first and then block number from lowest to highest
	assert.Equal(t, id3, returnedLongOrders[0].id)
	assert.Equal(t, price3, returnedLongOrders[0].Price)
	assert.Equal(t, blockNumber3, returnedLongOrders[0].BlockNumber)
	assert.Equal(t, id1, returnedLongOrders[1].id)
	assert.Equal(t, price1, returnedLongOrders[1].Price)
	assert.Equal(t, blockNumber1, returnedLongOrders[1].BlockNumber)
	assert.Equal(t, id2, returnedLongOrders[2].id)
	assert.Equal(t, price2, returnedLongOrders[2].Price)
	assert.Equal(t, blockNumber2, returnedLongOrders[2].BlockNumber)

	for _, returedOrder := range returnedLongOrders {
		assert.Equal(t, "long", returedOrder.PositionType)
		assert.Equal(t, userAddress, returedOrder.UserAddress)
		assert.Equal(t, longOrderBaseAssetQuantity, returedOrder.BaseAssetQuantity)
		assert.Equal(t, status, returedOrder.Status)
		assert.Equal(t, salt, returedOrder.Salt)
	}
}

func createLimitOrder(id uint64, positionType string, userAddress string, baseAssetQuantity int, price float64, status string, salt string, signature []byte, blockNumber uint64) LimitOrder {
	return LimitOrder{
		id:                id,
		PositionType:      positionType,
		UserAddress:       userAddress,
		BaseAssetQuantity: baseAssetQuantity,
		Price:             price,
		Status:            status,
		Salt:              salt,
		Signature:         signature,
		BlockNumber:       blockNumber,
	}
}
