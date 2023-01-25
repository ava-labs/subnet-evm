package limitorders

import (
	"math/big"
	"sort"
	"time"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

var maxLiquidationRatio *big.Int = big.NewInt(25 * 10e4)
var minSizeRequirement *big.Int = big.NewInt(0).Mul(big.NewInt(5), big.NewInt(1e18))

type Market int

const (
	AvaxPerp Market = iota
)

func GetActiveMarkets() []Market {
	return []Market{AvaxPerp}
}

type Collateral int

const (
	HUSD Collateral = iota
)

var collateralWeightMap map[Collateral]float64 = map[Collateral]float64{HUSD: 1}

type Status string

const (
	Placed    = "placed"
	Filled    = "filled"
	Cancelled = "cancelled"
)

type LimitOrder struct {
	id                      uint64
	Market                  Market
	PositionType            string
	UserAddress             string
	BaseAssetQuantity       *big.Int
	FilledBaseAssetQuantity *big.Int
	Price                   *big.Int
	Status                  Status
	Signature               []byte
	RawOrder                interface{}
	BlockNumber             *big.Int // block number order was placed on
}

func (order LimitOrder) GetUnFilledBaseAssetQuantity() *big.Int {
	return big.NewInt(0).Sub(order.BaseAssetQuantity, order.FilledBaseAssetQuantity)
}

type Position struct {
	OpenNotional         *big.Int
	Size                 *big.Int
	UnrealisedFunding    *big.Int
	LastPremiumFraction  *big.Int
	LiquidationThreshold *big.Int
}

type Trader struct {
	Positions   map[Market]*Position    // position for every market
	Margins     map[Collateral]*big.Int // available margin/balance for every market
	BlockNumber *big.Int
}

type LimitOrderDatabase interface {
	GetAllOrders() []LimitOrder
	Add(order *LimitOrder)
	Delete(signature []byte)
	UpdateFilledBaseAssetQuantity(quantity *big.Int, signature []byte)
	GetLongOrders(market Market) []LimitOrder
	GetShortOrders(market Market) []LimitOrder
	UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool)
	UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int)
	UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int)
	ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int)
	UpdateNextFundingTime(nextFundingTime uint64)
	GetNextFundingTime() uint64
	GetLiquidableTraders(market Market, oraclePrice *big.Int) []LiquidablePosition
	UpdateLastPrice(market Market, lastPrice *big.Int)
	GetLastPrice(market Market) *big.Int
}

type InMemoryDatabase struct {
	orderMap        map[string]*LimitOrder     // signature => order
	traderMap       map[common.Address]*Trader // address => trader info
	nextFundingTime uint64
	lastPrice       map[Market]*big.Int
}

func NewInMemoryDatabase() *InMemoryDatabase {
	orderMap := map[string]*LimitOrder{}
	lastPrice := map[Market]*big.Int{AvaxPerp: big.NewInt(0)}
	traderMap := map[common.Address]*Trader{}

	return &InMemoryDatabase{
		orderMap:        orderMap,
		traderMap:       traderMap,
		nextFundingTime: 0,
		lastPrice:       lastPrice,
	}
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

func (db *InMemoryDatabase) Delete(signature []byte) {
	deleteOrder(db, signature)
}

func (db *InMemoryDatabase) UpdateFilledBaseAssetQuantity(quantity *big.Int, signature []byte) {
	limitOrder := db.orderMap[string(signature)]
	if limitOrder.PositionType == "long" {
		limitOrder.FilledBaseAssetQuantity.Add(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled + quantity
	}
	if limitOrder.PositionType == "short" {
		limitOrder.FilledBaseAssetQuantity.Sub(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled - quantity
	}

	if limitOrder.BaseAssetQuantity.Cmp(limitOrder.FilledBaseAssetQuantity) == 0 {
		deleteOrder(db, signature)
	}
}

func (db *InMemoryDatabase) GetNextFundingTime() uint64 {
	return db.nextFundingTime
}

func (db *InMemoryDatabase) UpdateNextFundingTime(nextFundingTime uint64) {
	db.nextFundingTime = nextFundingTime
}

func (db *InMemoryDatabase) GetLongOrders(market Market) []LimitOrder {
	var longOrders []LimitOrder
	for _, order := range db.orderMap {
		if order.PositionType == "long" && order.Market == market {
			longOrders = append(longOrders, *order)
		}
	}
	sortLongOrders(longOrders)
	return longOrders
}

func (db *InMemoryDatabase) GetShortOrders(market Market) []LimitOrder {
	var shortOrders []LimitOrder
	for _, order := range db.orderMap {
		if order.PositionType == "short" && order.Market == market {
			shortOrders = append(shortOrders, *order)
		}
	}
	sortShortOrders(shortOrders)
	return shortOrders
}

func (db *InMemoryDatabase) UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int) {
	if _, ok := db.traderMap[trader]; !ok {
		db.traderMap[trader] = &Trader{
			Positions: map[Market]*Position{},
			Margins:   map[Collateral]*big.Int{},
		}
	}

	if _, ok := db.traderMap[trader].Margins[collateral]; !ok {
		db.traderMap[trader].Margins[collateral] = big.NewInt(0)
	}

	db.traderMap[trader].Margins[collateral].Add(db.traderMap[trader].Margins[collateral], addAmount)
}

func (db *InMemoryDatabase) UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool) {
	if _, ok := db.traderMap[trader]; !ok {
		db.traderMap[trader] = &Trader{
			Positions: map[Market]*Position{},
			Margins:   map[Collateral]*big.Int{},
		}
	}

	if _, ok := db.traderMap[trader].Positions[market]; !ok {
		db.traderMap[trader].Positions[market] = &Position{}
	}

	db.traderMap[trader].Positions[market].Size = size
	db.traderMap[trader].Positions[market].OpenNotional = openNotional

	if !isLiquidation {
		db.traderMap[trader].Positions[market].LiquidationThreshold = getLiquidationThreshold(size)
	}
}

func (db *InMemoryDatabase) UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int) {
	for _, trader := range db.traderMap {
		position := trader.Positions[market]
		if position != nil {
			position.UnrealisedFunding = dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size))
		}
	}
}

func (db *InMemoryDatabase) ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int) {
	if db.traderMap[trader] != nil {
		if _, ok := db.traderMap[trader].Positions[market]; ok {
			db.traderMap[trader].Positions[market].UnrealisedFunding = big.NewInt(0)
			db.traderMap[trader].Positions[market].LastPremiumFraction = cumulativePremiumFraction
		}
	}
}

func (db *InMemoryDatabase) UpdateLastPrice(market Market, lastPrice *big.Int) {
	db.lastPrice[market] = lastPrice
}

func (db *InMemoryDatabase) GetLastPrice(market Market) *big.Int {
	return db.lastPrice[market]
}

func sortLongOrders(orders []LimitOrder) []LimitOrder {
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].Price.Cmp(orders[j].Price) == 1 {
			return true
		}
		if orders[i].Price.Cmp(orders[j].Price) == 0 {
			if orders[i].BlockNumber.Cmp(orders[j].BlockNumber) == -1 {
				return true
			}
		}
		return false
	})
	return orders
}

func sortShortOrders(orders []LimitOrder) []LimitOrder {
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].Price.Cmp(orders[j].Price) == -1 {
			return true
		}
		if orders[i].Price.Cmp(orders[j].Price) == 0 {
			if orders[i].BlockNumber.Cmp(orders[j].BlockNumber) == -1 {
				return true
			}
		}
		return false
	})
	return orders
}

func getNextHour() time.Time {
	now := time.Now().UTC()
	nextHour := now.Round(time.Hour)
	if time.Since(nextHour) >= 0 {
		nextHour = nextHour.Add(time.Hour)
	}
	return nextHour
}

func deleteOrder(db *InMemoryDatabase, signature []byte) {
	delete(db.orderMap, string(signature))
}

func getLiquidationThreshold(size *big.Int) *big.Int {
	absSize := big.NewInt(0).Abs(size)
	maxLiquidationSize := divideByBasePrecision(big.NewInt(0).Mul(absSize, maxLiquidationRatio))
	threshold := big.NewInt(0).Add(maxLiquidationSize, big.NewInt(1))
	liquidationThreshold := utils.BigIntMax(threshold, minSizeRequirement)
	return big.NewInt(0).Mul(liquidationThreshold, big.NewInt(int64(size.Sign()))) // same sign as size
}
