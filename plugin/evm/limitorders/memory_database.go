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

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var _1e18 = big.NewInt(1e18)
var _1e6 = big.NewInt(1e6)

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

type PositionType int

const (
	LONG PositionType = iota
	SHORT
)

func (p PositionType) String() string {
	return [...]string{"long", "short"}[p]
}

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
	Id                      common.Hash
	Market                  Market
	PositionType            PositionType
	UserAddress             string
	BaseAssetQuantity       *big.Int
	FilledBaseAssetQuantity *big.Int
	Salt                    *big.Int
	Price                   *big.Int
	ReduceOnly              bool
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
	ReduceOnly              bool        `json:"reduce_only"`
}

func (order *LimitOrder) MarshalJSON() ([]byte, error) {
	limitOrderJson := LimitOrderJson{
		Market:                  order.Market,
		PositionType:            order.PositionType.String(),
		UserAddress:             strings.ToLower(order.UserAddress),
		BaseAssetQuantity:       order.BaseAssetQuantity,
		FilledBaseAssetQuantity: order.FilledBaseAssetQuantity,
		Salt:                    order.Salt,
		Price:                   order.Price,
		LifecycleList:           order.LifecycleList,
		Signature:               hex.EncodeToString(order.Signature),
		BlockNumber:             order.BlockNumber,
		ReduceOnly:              order.ReduceOnly,
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
	return fmt.Sprintf("LimitOrder: Market: %v, PositionType: %v, UserAddress: %v, BaseAssetQuantity: %s, FilledBaseAssetQuantity: %s, Salt: %v, Price: %s, ReduceOnly: %v, Signature: %v, BlockNumber: %s", order.Market, order.PositionType, order.UserAddress, prettifyScaledBigInt(order.BaseAssetQuantity, 18), prettifyScaledBigInt(order.FilledBaseAssetQuantity, 18), order.Salt, prettifyScaledBigInt(order.Price, 6), order.ReduceOnly, hex.EncodeToString(order.Signature), order.BlockNumber)
}

func (order LimitOrder) ToOrderMin() OrderMin {
	return OrderMin{
		Market:  order.Market,
		Price:   order.Price.String(),
		Size:    order.GetUnFilledBaseAssetQuantity().String(),
		Signer:  order.UserAddress,
		OrderId: order.Id.String(),
	}
}

type Position struct {
	OpenNotional         *big.Int `json:"open_notional"`
	Size                 *big.Int `json:"size"`
	UnrealisedFunding    *big.Int `json:"unrealised_funding"`
	LastPremiumFraction  *big.Int `json:"last_premium_fraction"`
	LiquidationThreshold *big.Int `json:"liquidation_threshold"`
}

type Margin struct {
	Reserved  *big.Int                `json:"reserved"`
	Deposited map[Collateral]*big.Int `json:"deposited"`
}

type Trader struct {
	Positions map[Market]*Position `json:"positions"` // position for every market
	Margin    Margin               `json:"margin"`    // available margin/balance for every market
}

type LimitOrderDatabase interface {
	GetAllOrders() []LimitOrder
	Add(orderId common.Hash, order *LimitOrder)
	Delete(orderId common.Hash)
	UpdateFilledBaseAssetQuantity(quantity *big.Int, orderId common.Hash, blockNumber uint64)
	GetLongOrders(market Market, cutoff *big.Int) []LimitOrder
	GetShortOrders(market Market, cutoff *big.Int) []LimitOrder
	UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool)
	UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int)
	UpdateReservedMargin(trader common.Address, addAmount *big.Int)
	UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int)
	ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int)
	UpdateNextFundingTime(nextFundingTime uint64)
	GetNextFundingTime() uint64
	UpdateLastPrice(market Market, lastPrice *big.Int)
	GetLastPrice(market Market) *big.Int
	GetLastPrices() map[Market]*big.Int
	GetAllTraders() map[common.Address]Trader
	GetOrderBookData() InMemoryDatabase
	Accept(blockNumber uint64)
	SetOrderStatus(orderId common.Hash, status Status, blockNumber uint64) error
	RevertLastStatus(orderId common.Hash) error
	GetNaughtyTraders(oraclePrices map[Market]*big.Int) ([]LiquidablePosition, map[common.Address][]common.Hash)
	GetOpenOrdersForTrader(trader common.Address) []LimitOrder
}

type InMemoryDatabase struct {
	mu              *sync.RWMutex               `json:"-"`
	OrderMap        map[common.Hash]*LimitOrder `json:"order_map"`  // ID => order
	TraderMap       map[common.Address]*Trader  `json:"trader_map"` // address => trader info
	NextFundingTime uint64                      `json:"next_funding_time"`
	LastPrice       map[Market]*big.Int         `json:"last_price"`
	configService   IConfigService
}

func NewInMemoryDatabase(configService IConfigService) *InMemoryDatabase {
	orderMap := map[common.Hash]*LimitOrder{}
	lastPrice := map[Market]*big.Int{AvaxPerp: big.NewInt(0)}
	traderMap := map[common.Address]*Trader{}

	return &InMemoryDatabase{
		OrderMap:        orderMap,
		TraderMap:       traderMap,
		NextFundingTime: 0,
		LastPrice:       lastPrice,
		mu:              &sync.RWMutex{},
		configService:   configService,
	}
}

func (db *InMemoryDatabase) Accept(blockNumber uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

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
	db.mu.RLock() // only read lock required
	defer db.mu.RUnlock()

	allOrders := []LimitOrder{}
	for _, order := range db.OrderMap {
		allOrders = append(allOrders, deepCopyOrder(*order))
	}
	return allOrders
}

func (db *InMemoryDatabase) Add(orderId common.Hash, order *LimitOrder) {
	db.mu.Lock()
	defer db.mu.Unlock()

	order.Id = orderId
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

	if limitOrder.PositionType == LONG {
		limitOrder.FilledBaseAssetQuantity.Add(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled + quantity
	}
	if limitOrder.PositionType == SHORT {
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
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.NextFundingTime
}

func (db *InMemoryDatabase) UpdateNextFundingTime(nextFundingTime uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.NextFundingTime = nextFundingTime
}

func (db *InMemoryDatabase) GetLongOrders(market Market, cutoff *big.Int) []LimitOrder {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var longOrders []LimitOrder
	for _, order := range db.OrderMap {
		if order.PositionType == LONG &&
			order.Market == market &&
			order.getOrderStatus().Status == Placed &&
			(cutoff == nil || order.Price.Cmp(cutoff) <= 0) &&
			// this will filter orders that are reduce only but with size > current position size (basically no partial fills) - @todo: think if this is correct
			(!order.ReduceOnly || db.willReducePosition(order)) {
			longOrders = append(longOrders, deepCopyOrder(*order))
		}
	}
	sortLongOrders(longOrders)
	return longOrders
}

func (db *InMemoryDatabase) GetShortOrders(market Market, cutoff *big.Int) []LimitOrder {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var shortOrders []LimitOrder
	for _, order := range db.OrderMap {
		if order.PositionType == SHORT &&
			order.Market == market &&
			order.getOrderStatus().Status == Placed &&
			(cutoff == nil || order.Price.Cmp(cutoff) >= 0) &&
			// this will filter orders that are reduce only but with size > current position size (basically no partial fills) - @todo: think if this is correct
			(!order.ReduceOnly || db.willReducePosition(order)) {
			shortOrders = append(shortOrders, deepCopyOrder(*order))
		}
	}
	sortShortOrders(shortOrders)
	return shortOrders
}

func (db *InMemoryDatabase) UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = getBlankTrader()
	}

	if _, ok := db.TraderMap[trader].Margin.Deposited[collateral]; !ok {
		db.TraderMap[trader].Margin.Deposited[collateral] = big.NewInt(0)
	}

	db.TraderMap[trader].Margin.Deposited[collateral].Add(db.TraderMap[trader].Margin.Deposited[collateral], addAmount)
	log.Info("UpdateMargin", "trader", trader.String(), "collateral", collateral, "updated margin", db.TraderMap[trader].Margin.Deposited[collateral].Uint64())
}

func (db *InMemoryDatabase) UpdateReservedMargin(trader common.Address, addAmount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = getBlankTrader()
	}

	db.TraderMap[trader].Margin.Reserved.Add(db.TraderMap[trader].Margin.Reserved, addAmount)
	log.Info("UpdateReservedMargin", "trader", trader.String(), "updated reserved margin", db.TraderMap[trader].Margin.Reserved.Uint64())
}

func (db *InMemoryDatabase) UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = getBlankTrader()
	}

	if _, ok := db.TraderMap[trader].Positions[market]; !ok {
		db.TraderMap[trader].Positions[market] = &Position{}
	}

	db.TraderMap[trader].Positions[market].Size = size
	db.TraderMap[trader].Positions[market].OpenNotional = openNotional
	db.TraderMap[trader].Positions[market].LastPremiumFraction = big.NewInt(0)

	if !isLiquidation {
		db.TraderMap[trader].Positions[market].LiquidationThreshold = getLiquidationThreshold(db.configService.getMaxLiquidationRatio(), db.configService.getMinSizeRequirement(), size)
	}

	if db.TraderMap[trader].Positions[market].UnrealisedFunding == nil {
		db.TraderMap[trader].Positions[market].UnrealisedFunding = big.NewInt(0)
	}
	// adjust the liquidation threshold if > resultant position size (for both isLiquidation = true/false)
	threshold := utils.BigIntMinAbs(db.TraderMap[trader].Positions[market].LiquidationThreshold, size)
	db.TraderMap[trader].Positions[market].LiquidationThreshold.Mul(threshold, big.NewInt(int64(size.Sign()))) // same sign as size
}

func (db *InMemoryDatabase) UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, trader := range db.TraderMap {
		position := trader.Positions[market]
		if position != nil {
			position.UnrealisedFunding = dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size))
		}
	}
}

func (db *InMemoryDatabase) ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.TraderMap[trader] != nil {
		if _, ok := db.TraderMap[trader].Positions[market]; ok {
			db.TraderMap[trader].Positions[market].UnrealisedFunding = big.NewInt(0)
			db.TraderMap[trader].Positions[market].LastPremiumFraction = cumulativePremiumFraction
		}
	}
}

func (db *InMemoryDatabase) UpdateLastPrice(market Market, lastPrice *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.LastPrice[market] = lastPrice
}

func (db *InMemoryDatabase) GetLastPrice(market Market) *big.Int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.LastPrice[market]
}

func (db *InMemoryDatabase) GetLastPrices() map[Market]*big.Int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.LastPrice
}

func (db *InMemoryDatabase) GetAllTraders() map[common.Address]Trader {
	db.mu.RLock()
	defer db.mu.RUnlock()

	traderMap := map[common.Address]Trader{}
	for address, trader := range db.TraderMap {
		traderMap[address] = *trader
	}
	return traderMap
}

func (db *InMemoryDatabase) GetOpenOrdersForTrader(trader common.Address) []LimitOrder {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.getTraderOrders(trader)
}

func determinePositionToLiquidate(trader *Trader, addr common.Address, marginFraction *big.Int) LiquidablePosition {
	liquidable := LiquidablePosition{}
	// iterate through the markets and return the first one with an open position
	// @todo when we introduce multiple markets, we will have to implement a more sophisticated liquidation strategy
	for _, market := range GetActiveMarkets() {
		position := trader.Positions[market]
		if position == nil || position.Size.Sign() == 0 {
			continue
		}
		liquidable = LiquidablePosition{
			Address:        addr,
			Market:         market,
			Size:           position.LiquidationThreshold,
			MarginFraction: marginFraction,
			FilledSize:     big.NewInt(0),
		}
		if position.Size.Sign() == -1 {
			liquidable.PositionType = SHORT
		} else {
			liquidable.PositionType = LONG
		}
	}
	return liquidable
}

func (db *InMemoryDatabase) GetNaughtyTraders(oraclePrices map[Market]*big.Int) ([]LiquidablePosition, map[common.Address][]common.Hash) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	liquidablePositions := []LiquidablePosition{}
	ordersToCancel := map[common.Address][]common.Hash{}

	for addr, trader := range db.TraderMap {
		pendingFunding := getTotalFunding(trader)
		marginFraction := calcMarginFraction(trader, pendingFunding, oraclePrices, db.LastPrice)
		if marginFraction.Cmp(db.configService.getMaintenanceMargin()) == -1 {
			log.Info("below maintenanceMargin", "trader", addr.String(), "marginFraction", prettifyScaledBigInt(marginFraction, 6))
			liquidablePositions = append(liquidablePositions, determinePositionToLiquidate(trader, addr, marginFraction))
			continue // we do not check for their open orders yet. Maybe liquidating them first will make available margin positive
		}
		availableMargin := getAvailableMargin(trader, pendingFunding, oraclePrices, db.LastPrice, db.configService.getMinAllowableMargin())
		// log.Info("getAvailableMargin", "trader", addr.String(), "availableMargin", prettifyScaledBigInt(availableMargin, 6))
		if availableMargin.Cmp(big.NewInt(0)) == -1 {
			log.Info("negative available margin", "trader", addr.String(), "availableMargin", prettifyScaledBigInt(availableMargin, 6))
			db.determineOrdersToCancel(addr, trader, availableMargin, oraclePrices, ordersToCancel)
		}
	}

	// lower margin fraction positions should be liquidated first
	sortLiquidableSliceByMarginFraction(liquidablePositions)
	return liquidablePositions, ordersToCancel
}

// assumes db.mu.RLock has been held by the caller
func (db *InMemoryDatabase) determineOrdersToCancel(addr common.Address, trader *Trader, availableMargin *big.Int, oraclePrices map[Market]*big.Int, ordersToCancel map[common.Address][]common.Hash) {
	traderOrders := db.getTraderOrders(addr)
	sort.Slice(traderOrders, func(i, j int) bool {
		// higher diff comes first
		iDiff := big.NewInt(0).Abs(big.NewInt(0).Sub(traderOrders[i].Price, oraclePrices[traderOrders[i].Market]))
		jDiff := big.NewInt(0).Abs(big.NewInt(0).Sub(traderOrders[j].Price, oraclePrices[traderOrders[j].Market]))
		return iDiff.Cmp(jDiff) > 0
	})

	_availableMargin := new(big.Int).Set(availableMargin)
	if len(traderOrders) > 0 {
		// cancel orders until available margin is positive
		ordersToCancel[addr] = []common.Hash{}
		for _, order := range traderOrders {
			// cannot cancel ReduceOnly orders because no margin is reserved for them
			if order.ReduceOnly {
				continue
			}
			ordersToCancel[addr] = append(ordersToCancel[addr], order.Id)
			orderNotional := big.NewInt(0).Abs(big.NewInt(0).Div(big.NewInt(0).Mul(order.GetUnFilledBaseAssetQuantity(), order.Price), _1e18)) // | size * current price |
			marginReleased := divideByBasePrecision(big.NewInt(0).Mul(orderNotional, db.configService.getMinAllowableMargin()))
			_availableMargin.Add(_availableMargin, marginReleased)
			log.Info("in determineOrdersToCancel loop", "availableMargin", prettifyScaledBigInt(_availableMargin, 6), "marginReleased", prettifyScaledBigInt(marginReleased, 6), "orderNotional", prettifyScaledBigInt(orderNotional, 6))
			if _availableMargin.Cmp(big.NewInt(0)) >= 0 {
				break
			}
		}
	}
}

func (db *InMemoryDatabase) getTraderOrders(trader common.Address) []LimitOrder {
	traderOrders := []LimitOrder{}
	for _, order := range db.OrderMap {
		if strings.EqualFold(order.UserAddress, trader.String()) {
			traderOrders = append(traderOrders, deepCopyOrder(*order))
		}
	}
	return traderOrders
}

func (db *InMemoryDatabase) willReducePosition(order *LimitOrder) bool {
	trader := common.HexToAddress(order.UserAddress)
	if db.TraderMap[trader] == nil {
		return false
	}
	positions := db.TraderMap[trader].Positions
	if position, ok := positions[order.Market]; ok {
		// position.Size, order.BaseAssetQuantity need to be of opposite sign and abs(position.Size) >= abs(order.BaseAssetQuantity)
		return position.Size.Sign() != order.BaseAssetQuantity.Sign() && big.NewInt(0).Abs(position.Size).Cmp(big.NewInt(0).Abs(order.BaseAssetQuantity)) != -1
	} else {
		return false
	}
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

func (db *InMemoryDatabase) GetOrderBookData() InMemoryDatabase {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return *db
}

func getLiquidationThreshold(maxLiquidationRatio *big.Int, minSizeRequirement *big.Int, size *big.Int) *big.Int {
	absSize := big.NewInt(0).Abs(size)
	maxLiquidationSize := divideByBasePrecision(big.NewInt(0).Mul(absSize, maxLiquidationRatio))
	liquidationThreshold := utils.BigIntMax(maxLiquidationSize, minSizeRequirement)
	return big.NewInt(0).Mul(liquidationThreshold, big.NewInt(int64(size.Sign()))) // same sign as size
}

func getBlankTrader() *Trader {
	return &Trader{
		Positions: map[Market]*Position{},
		Margin: Margin{
			Reserved: big.NewInt(0),
			Deposited: map[Collateral]*big.Int{
				0: big.NewInt(0),
			},
		},
	}
}

func getAvailableMargin(trader *Trader, pendingFunding *big.Int, oraclePrices map[Market]*big.Int, lastPrices map[Market]*big.Int, minAllowableMargin *big.Int) *big.Int {
	// log.Info("in getAvailableMargin", "trader", trader, "pendingFunding", pendingFunding, "oraclePrices", oraclePrices, "lastPrices", lastPrices)
	margin := new(big.Int).Sub(getNormalisedMargin(trader), pendingFunding)
	notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(trader, margin, Min_Allowable_Margin, oraclePrices, lastPrices)
	utilisedMargin := divideByBasePrecision(new(big.Int).Mul(notionalPosition, minAllowableMargin))
	// print margin, notionalPosition, unrealizePnL, utilisedMargin
	// log.Info("stats", "margin", margin, "notionalPosition", notionalPosition, "unrealizePnL", unrealizePnL, "utilisedMargin", utilisedMargin)
	return new(big.Int).Sub(
		new(big.Int).Add(margin, unrealizePnL),
		new(big.Int).Add(utilisedMargin, trader.Margin.Reserved),
	)
}

// deepCopyOrder deep copies the LimitOrder struct
func deepCopyOrder(order LimitOrder) LimitOrder {
	lifecycleList := &order.LifecycleList
	return LimitOrder{
		Id:                      order.Id,
		Market:                  order.Market,
		PositionType:            order.PositionType,
		UserAddress:             order.UserAddress,
		BaseAssetQuantity:       big.NewInt(0).Set(order.BaseAssetQuantity),
		FilledBaseAssetQuantity: big.NewInt(0).Set(order.FilledBaseAssetQuantity),
		Salt:                    big.NewInt(0).Set(order.Salt),
		Price:                   big.NewInt(0).Set(order.Price),
		ReduceOnly:              order.ReduceOnly,
		LifecycleList:           *lifecycleList,
		Signature:               order.Signature,
		BlockNumber:             big.NewInt(0).Set(order.BlockNumber),
		RawOrder:                order.RawOrder,
	}
}
