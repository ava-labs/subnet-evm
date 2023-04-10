package limitorders

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var _1e18 = big.NewInt(1e18)
var _1e6 = big.NewInt(1e6)

var maxLiquidationRatio *big.Int = big.NewInt(25 * 1e4) // 25%
var minSizeRequirement *big.Int = big.NewInt(0).Mul(big.NewInt(5), _1e18)

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

type Status uint8

const (
	Placed Status = iota
	FulFilled
	Cancelled
	Execution_Failed
)

type Lifecycle struct {
	BlockNumber uint64
	Status      Status
}

type LimitOrder struct {
	Market Market
	// @todo make this an enum
	PositionType            string
	UserAddress             string
	BaseAssetQuantity       *big.Int
	FilledBaseAssetQuantity *big.Int
	Salt                    *big.Int
	Price                   *big.Int
	LifecycleList           []Lifecycle
	Signature               []byte
	BlockNumber             *big.Int    // block number order was placed on
	RawOrder                interface{} `json:"-"`
}

type LimitOrderJson struct {
	Market                  Market      `json:"market"`
	PositionType            string      `json:"position_type"`
	UserAddress             string      `json:"user_address"`
	BaseAssetQuantity       *big.Int    `json:"base_asset_quantity"`
	FilledBaseAssetQuantity *big.Int    `json:"filled_base_asset_quantity"`
	Salt                    *big.Int    `json:"salt"`
	Price                   *big.Int    `json:"price"`
	LifecycleList           []Lifecycle `json:"lifecycle_list"`
	Signature               string      `json:"signature"`
	BlockNumber             *big.Int    `json:"block_number"` // block number order was placed on
}

func (order *LimitOrder) MarshalJSON() ([]byte, error) {
	limitOrderJson := LimitOrderJson{
		Market:                  order.Market,
		PositionType:            order.PositionType,
		UserAddress:             strings.ToLower(order.UserAddress),
		BaseAssetQuantity:       order.BaseAssetQuantity,
		FilledBaseAssetQuantity: order.FilledBaseAssetQuantity,
		Salt:                    order.Salt,
		Price:                   order.Price,
		LifecycleList:           order.LifecycleList,
		Signature:               hex.EncodeToString(order.Signature),
		BlockNumber:             order.BlockNumber,
	}
	return json.Marshal(limitOrderJson)
}

func (order LimitOrder) GetUnFilledBaseAssetQuantity() *big.Int {
	return big.NewInt(0).Sub(order.BaseAssetQuantity, order.FilledBaseAssetQuantity)
}

func (order LimitOrder) getOrderStatus() Lifecycle {
	lifecycle := order.LifecycleList
	return lifecycle[len(lifecycle)-1]
}

func (order LimitOrder) String() string {
	return fmt.Sprintf("LimitOrder: Market: %v, PositionType: %v, UserAddress: %v, BaseAssetQuantity: %s, FilledBaseAssetQuantity: %s, Salt: %v, Price: %s, Signature: %v, BlockNumber: %s", order.Market, order.PositionType, order.UserAddress, prettifyScaledBigInt(order.BaseAssetQuantity, 18), prettifyScaledBigInt(order.FilledBaseAssetQuantity, 18), order.Salt, prettifyScaledBigInt(order.Price, 6), hex.EncodeToString(order.Signature), order.BlockNumber)
}

type Position struct {
	OpenNotional         *big.Int `json:"open_notional"`
	Size                 *big.Int `json:"size"`
	UnrealisedFunding    *big.Int `json:"unrealised_funding"`
	LastPremiumFraction  *big.Int `json:"last_premium_fraction"`
	LiquidationThreshold *big.Int `json:"liquidation_threshold"`
}

type Trader struct {
	Positions map[Market]*Position    `json:"positions"` // position for every market
	Margins   map[Collateral]*big.Int `json:"margins"`   // available margin/balance for every market
}

type LimitOrderDatabase interface {
	GetAllOrders() []LimitOrder
	Add(orderId common.Hash, order *LimitOrder)
	Delete(orderId common.Hash)
	UpdateFilledBaseAssetQuantity(quantity *big.Int, orderId common.Hash, blockNumber uint64)
	GetLongOrders(market Market) []LimitOrder
	GetShortOrders(market Market) []LimitOrder
	UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool)
	UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int)
	UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int)
	ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int)
	UpdateNextFundingTime(nextFundingTime uint64)
	GetNextFundingTime() uint64
	UpdateLastPrice(market Market, lastPrice *big.Int)
	GetLastPrice(market Market) *big.Int
	GetAllTraders() map[common.Address]Trader
	GetOrderBookData() InMemoryDatabase
	Accept(blockNumber uint64)
	SetOrderStatus(orderId common.Hash, status Status, blockNumber uint64) error
	RevertLastStatus(orderId common.Hash) error
}

type InMemoryDatabase struct {
	mu              sync.Mutex                  `json:"-"`
	OrderMap        map[common.Hash]*LimitOrder `json:"order_map"`  // ID => order
	TraderMap       map[common.Address]*Trader  `json:"trader_map"` // address => trader info
	NextFundingTime uint64                      `json:"next_funding_time"`
	LastPrice       map[Market]*big.Int         `json:"last_price"`
}

func NewInMemoryDatabase() *InMemoryDatabase {
	orderMap := map[common.Hash]*LimitOrder{}
	lastPrice := map[Market]*big.Int{AvaxPerp: big.NewInt(0)}
	traderMap := map[common.Address]*Trader{}

	return &InMemoryDatabase{
		OrderMap:        orderMap,
		TraderMap:       traderMap,
		NextFundingTime: 0,
		LastPrice:       lastPrice,
	}
}

func (db *InMemoryDatabase) Accept(blockNumber uint64) {
	for orderId, order := range db.OrderMap {
		lifecycle := order.getOrderStatus()
		if lifecycle.Status != Placed && lifecycle.BlockNumber <= blockNumber {
			delete(db.OrderMap, orderId)
		}
	}
}

func (db *InMemoryDatabase) SetOrderStatus(orderId common.Hash, status Status, blockNumber uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.OrderMap[orderId] == nil {
		return errors.New(fmt.Sprintf("Invalid orderId %s", orderId.Hex()))
	}
	db.OrderMap[orderId].LifecycleList = append(db.OrderMap[orderId].LifecycleList, Lifecycle{blockNumber, status})
	log.Info("SetOrderStatus", "orderId", orderId.String(), "status", status, "updated state", db.OrderMap[orderId].LifecycleList)
	return nil
}

func (db *InMemoryDatabase) RevertLastStatus(orderId common.Hash) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.OrderMap[orderId] == nil {
		return errors.New(fmt.Sprintf("Invalid orderId %s", orderId.Hex()))
	}

	lifeCycleList := db.OrderMap[orderId].LifecycleList
	if len(lifeCycleList) > 0 {
		db.OrderMap[orderId].LifecycleList = lifeCycleList[:len(lifeCycleList)-1]
	}
	return nil
}

func (db *InMemoryDatabase) GetAllOrders() []LimitOrder {
	db.mu.Lock()
	defer db.mu.Unlock()

	allOrders := []LimitOrder{}
	for _, order := range db.OrderMap {
		allOrders = append(allOrders, *order)
	}
	return allOrders
}

func (db *InMemoryDatabase) Add(orderId common.Hash, order *LimitOrder) {
	db.mu.Lock()
	defer db.mu.Unlock()

	order.LifecycleList = append(order.LifecycleList, Lifecycle{order.BlockNumber.Uint64(), Placed})
	db.OrderMap[orderId] = order
}

func (db *InMemoryDatabase) Delete(orderId common.Hash) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.OrderMap, orderId)
}

func (db *InMemoryDatabase) UpdateFilledBaseAssetQuantity(quantity *big.Int, orderId common.Hash, blockNumber uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	limitOrder := db.OrderMap[orderId]

	if limitOrder.PositionType == "long" {
		limitOrder.FilledBaseAssetQuantity.Add(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled + quantity
	}
	if limitOrder.PositionType == "short" {
		limitOrder.FilledBaseAssetQuantity.Sub(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled - quantity
	}

	if limitOrder.BaseAssetQuantity.Cmp(limitOrder.FilledBaseAssetQuantity) == 0 {
		limitOrder.LifecycleList = append(limitOrder.LifecycleList, Lifecycle{blockNumber, FulFilled})
	}

	if quantity.Cmp(big.NewInt(0)) == -1 && limitOrder.getOrderStatus().Status == FulFilled {
		// handling reorgs
		limitOrder.LifecycleList = limitOrder.LifecycleList[:len(limitOrder.LifecycleList)-1]
	}
}

func (db *InMemoryDatabase) GetNextFundingTime() uint64 {
	return db.NextFundingTime
}

func (db *InMemoryDatabase) UpdateNextFundingTime(nextFundingTime uint64) {
	db.NextFundingTime = nextFundingTime
}

func (db *InMemoryDatabase) GetLongOrders(market Market) []LimitOrder {
	var longOrders []LimitOrder
	for _, order := range db.OrderMap {
		if order.PositionType == "long" &&
			order.Market == market &&
			order.getOrderStatus().Status == Placed { // &&
			// order.Price.Cmp(big.NewInt(20e6)) <= 0 { // hardcode amm spread check eligibility for now
			longOrders = append(longOrders, *order)
		}
	}
	sortLongOrders(longOrders)
	return longOrders
}

func (db *InMemoryDatabase) GetShortOrders(market Market) []LimitOrder {
	var shortOrders []LimitOrder
	for _, order := range db.OrderMap {
		if order.PositionType == "short" &&
			order.Market == market &&
			order.getOrderStatus().Status == Placed {
			shortOrders = append(shortOrders, *order)
		}
	}
	sortShortOrders(shortOrders)
	return shortOrders
}

func (db *InMemoryDatabase) UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int) {
	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = &Trader{
			Positions: map[Market]*Position{},
			Margins:   map[Collateral]*big.Int{},
		}
	}

	if _, ok := db.TraderMap[trader].Margins[collateral]; !ok {
		db.TraderMap[trader].Margins[collateral] = big.NewInt(0)
	}

	db.TraderMap[trader].Margins[collateral].Add(db.TraderMap[trader].Margins[collateral], addAmount)
	log.Info("UpdateMargin", "trader", trader.String(), "collateral", collateral, "updated margin", db.TraderMap[trader].Margins[collateral].Uint64())
}

func (db *InMemoryDatabase) UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool) {
	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = &Trader{
			Positions: map[Market]*Position{},
			Margins:   map[Collateral]*big.Int{},
		}
	}

	if _, ok := db.TraderMap[trader].Positions[market]; !ok {
		db.TraderMap[trader].Positions[market] = &Position{}
	}

	db.TraderMap[trader].Positions[market].Size = size
	db.TraderMap[trader].Positions[market].OpenNotional = openNotional
	db.TraderMap[trader].Positions[market].LastPremiumFraction = big.NewInt(0)

	if !isLiquidation {
		db.TraderMap[trader].Positions[market].LiquidationThreshold = getLiquidationThreshold(size)
	}
	// adjust the liquidation threshold if > resultant position size (for both isLiquidation = true/false)
	threshold := utils.BigIntMinAbs(db.TraderMap[trader].Positions[market].LiquidationThreshold, size)
	db.TraderMap[trader].Positions[market].LiquidationThreshold.Mul(threshold, big.NewInt(int64(size.Sign()))) // same sign as size
}

func (db *InMemoryDatabase) UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int) {
	for _, trader := range db.TraderMap {
		position := trader.Positions[market]
		if position != nil {
			position.UnrealisedFunding = dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size))
		}
	}
}

func (db *InMemoryDatabase) ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int) {
	if db.TraderMap[trader] != nil {
		if _, ok := db.TraderMap[trader].Positions[market]; ok {
			db.TraderMap[trader].Positions[market].UnrealisedFunding = big.NewInt(0)
			db.TraderMap[trader].Positions[market].LastPremiumFraction = cumulativePremiumFraction
		}
	}
}

func (db *InMemoryDatabase) UpdateLastPrice(market Market, lastPrice *big.Int) {
	db.LastPrice[market] = lastPrice
}

func (db *InMemoryDatabase) GetLastPrice(market Market) *big.Int {
	return db.LastPrice[market]
}

func (db *InMemoryDatabase) GetAllTraders() map[common.Address]Trader {
	traderMap := map[common.Address]Trader{}
	for address, trader := range db.TraderMap {
		traderMap[address] = *trader
	}
	return traderMap
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

func deleteOrder(db *InMemoryDatabase, id common.Hash) {
	log.Info("#### deleting order", "orderId", id)
	delete(db.OrderMap, id)
}

func (db *InMemoryDatabase) GetOrderBookData() InMemoryDatabase {
	return *db
}

func getLiquidationThreshold(size *big.Int) *big.Int {
	absSize := big.NewInt(0).Abs(size)
	maxLiquidationSize := divideByBasePrecision(big.NewInt(0).Mul(absSize, maxLiquidationRatio))
	threshold := big.NewInt(0).Add(maxLiquidationSize, big.NewInt(1))
	liquidationThreshold := utils.BigIntMax(threshold, minSizeRequirement)
	return big.NewInt(0).Mul(liquidationThreshold, big.NewInt(int64(size.Sign()))) // same sign as size
}
