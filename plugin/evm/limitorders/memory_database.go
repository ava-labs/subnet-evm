package limitorders

import (
	"math"
	"sort"
)

type LimitOrder struct {
	id                      uint64
	PositionType            string
	UserAddress             string
	BaseAssetQuantity       int
	FilledBaseAssetQuantity int
	Price                   float64
	Status                  string
	Salt                    int64
	Signature               []byte
	RawOrder                interface{}
	RawSignature            interface{}
	BlockNumber             uint64
}

type LimitOrderDatabase interface {
	GetAllOrders() []LimitOrder
	Add(order *LimitOrder)
	UpdateFilledBaseAssetQuantity(quantity uint, signature []byte)
	Delete(signature []byte)
	GetLongOrders() []LimitOrder
	GetShortOrders() []LimitOrder
}

type InMemoryDatabase struct {
	orderMap map[string]*LimitOrder
}

func NewInMemoryDatabase() *InMemoryDatabase {
	orderMap := map[string]*LimitOrder{}
	return &InMemoryDatabase{orderMap}
}

func (db *InMemoryDatabase) GetAllOrders() []LimitOrder {
	allOrders := []LimitOrder{}
	for _, order := range db.orderMap {
		allOrders = append(allOrders, *order)
	}
	return allOrders
}

func (db *InMemoryDatabase) Add(order *LimitOrder) {
	db.orderMap[string(order.Signature)] = order
}

func (db *InMemoryDatabase) UpdateFilledBaseAssetQuantity(quantity uint, signature []byte) {
	limitOrder := db.orderMap[string(signature)]
	if uint(math.Abs(float64(limitOrder.BaseAssetQuantity))) == quantity {
		deleteOrder(db, signature)
		return
	} else {
		if limitOrder.PositionType == "long" {
			limitOrder.FilledBaseAssetQuantity = int(quantity)
		}
		if limitOrder.PositionType == "short" {
			limitOrder.FilledBaseAssetQuantity = -int(quantity)
		}
	}
}

// Deletes silently
func (db *InMemoryDatabase) Delete(signature []byte) {
	deleteOrder(db, signature)
}

func (db *InMemoryDatabase) GetLongOrders() []LimitOrder {
	var longOrders []LimitOrder
	for _, order := range db.orderMap {
		if order.PositionType == "long" {
			longOrders = append(longOrders, *order)
		}
	}
	sortLongOrders(longOrders)
	return longOrders
}

func (db *InMemoryDatabase) GetShortOrders() []LimitOrder {
	var shortOrders []LimitOrder
	for _, order := range db.orderMap {
		if order.PositionType == "short" {
			shortOrders = append(shortOrders, *order)
		}
	}
	sortShortOrders(shortOrders)
	return shortOrders
}

func sortLongOrders(orders []LimitOrder) []LimitOrder {
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].Price > orders[j].Price {
			return true
		}
		if orders[i].Price == orders[j].Price {
			if orders[i].BlockNumber < orders[j].BlockNumber {
				return true
			}
		}
		return false
	})
	return orders
}

func sortShortOrders(orders []LimitOrder) []LimitOrder {
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].Price < orders[j].Price {
			return true
		}
		if orders[i].Price == orders[j].Price {
			if orders[i].BlockNumber < orders[j].BlockNumber {
				return true
			}
		}
		return false
	})
	return orders
}

func deleteOrder(db *InMemoryDatabase, signature []byte) {
	delete(db.orderMap, string(signature))
}
