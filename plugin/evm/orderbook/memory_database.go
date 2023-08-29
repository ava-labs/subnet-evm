package orderbook

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type InMemoryDatabase struct {
	mu                        *sync.RWMutex              `json:"-"`
	OrderMap                  map[common.Hash]*Order     `json:"order_map"`  // ID => order
	TraderMap                 map[common.Address]*Trader `json:"trader_map"` // address => trader info
	NextFundingTime           uint64                     `json:"next_funding_time"`
	LastPrice                 map[Market]*big.Int        `json:"last_price"`
	CumulativePremiumFraction map[Market]*big.Int        `json:"cumulative_last_premium_fraction"`
	configService             IConfigService
}

func NewInMemoryDatabase(configService IConfigService) *InMemoryDatabase {
	orderMap := map[common.Hash]*Order{}
	lastPrice := map[Market]*big.Int{}
	traderMap := map[common.Address]*Trader{}

	return &InMemoryDatabase{
		OrderMap:                  orderMap,
		TraderMap:                 traderMap,
		NextFundingTime:           0,
		LastPrice:                 lastPrice,
		CumulativePremiumFraction: map[Market]*big.Int{},
		mu:                        &sync.RWMutex{},
		configService:             configService,
	}
}

var (
	_1e18 = big.NewInt(1e18)
	_1e6  = big.NewInt(1e6)
)

const (
	RETRY_AFTER_BLOCKS = 10
)

type Market int64

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

type OrderType uint8

const (
	LimitOrderType OrderType = iota
	IOCOrderType
)

func (o OrderType) String() string {
	return [...]string{"limit", "ioc"}[o]
}

type Lifecycle struct {
	BlockNumber uint64
	Status      Status
	Info        string
}

type Order struct {
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
	BlockNumber             *big.Int      // block number order was placed on
	RawOrder                ContractOrder `json:"-"`
	OrderType               OrderType
}

func (order *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Market                  Market      `json:"market"`
		PositionType            string      `json:"position_type"`
		UserAddress             string      `json:"user_address"`
		BaseAssetQuantity       string      `json:"base_asset_quantity"`
		FilledBaseAssetQuantity string      `json:"filled_base_asset_quantity"`
		Salt                    string      `json:"salt"`
		Price                   string      `json:"price"`
		LifecycleList           []Lifecycle `json:"lifecycle_list"`
		BlockNumber             uint64      `json:"block_number"` // block number order was placed on
		ReduceOnly              bool        `json:"reduce_only"`
		OrderType               string      `json:"order_type"`
	}{
		Market:                  order.Market,
		PositionType:            order.PositionType.String(),
		UserAddress:             order.UserAddress,
		BaseAssetQuantity:       order.BaseAssetQuantity.String(),
		FilledBaseAssetQuantity: order.FilledBaseAssetQuantity.String(),
		Salt:                    order.Salt.String(),
		Price:                   order.Price.String(),
		LifecycleList:           order.LifecycleList,
		BlockNumber:             order.BlockNumber.Uint64(),
		ReduceOnly:              order.ReduceOnly,
		OrderType:               order.OrderType.String(),
	})
}

func (order Order) GetUnFilledBaseAssetQuantity() *big.Int {
	return big.NewInt(0).Sub(order.BaseAssetQuantity, order.FilledBaseAssetQuantity)
}

func (order Order) getOrderStatus() Lifecycle {
	lifecycle := order.LifecycleList
	return lifecycle[len(lifecycle)-1]
}

func (order Order) getExpireAt() *big.Int {
	if order.OrderType == IOCOrderType {
		return order.RawOrder.(*IOCOrder).ExpireAt
	}
	return big.NewInt(0)
}

func (order Order) String() string {
	return fmt.Sprintf("Order: Id: %s, OrderType: %s, Market: %v, PositionType: %v, UserAddress: %v, BaseAssetQuantity: %s, FilledBaseAssetQuantity: %s, Salt: %v, Price: %s, ReduceOnly: %v, BlockNumber: %s", order.Id, order.OrderType, order.Market, order.PositionType, order.UserAddress, prettifyScaledBigInt(order.BaseAssetQuantity, 18), prettifyScaledBigInt(order.FilledBaseAssetQuantity, 18), order.Salt, prettifyScaledBigInt(order.Price, 6), order.ReduceOnly, order.BlockNumber)
}

func (order Order) ToOrderMin() OrderMin {
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

func (p *Position) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		OpenNotional         string `json:"open_notional"`
		Size                 string `json:"size"`
		UnrealisedFunding    string `json:"unrealised_funding"`
		LastPremiumFraction  string `json:"last_premium_fraction"`
		LiquidationThreshold string `json:"liquidation_threshold"`
	}{
		OpenNotional:         p.OpenNotional.String(),
		Size:                 p.Size.String(),
		UnrealisedFunding:    p.UnrealisedFunding.String(),
		LastPremiumFraction:  p.LastPremiumFraction.String(),
		LiquidationThreshold: p.LiquidationThreshold.String(),
	})
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
	LoadFromSnapshot(snapshot Snapshot) error
	GetAllOrders() []Order
	Add(order *Order)
	Delete(orderId common.Hash)
	UpdateFilledBaseAssetQuantity(quantity *big.Int, orderId common.Hash, blockNumber uint64)
	GetLongOrders(market Market, lowerbound *big.Int, blockNumber *big.Int) []Order
	GetShortOrders(market Market, upperbound *big.Int, blockNumber *big.Int) []Order
	UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool, blockNumber uint64)
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
	GetOrderBookDataCopy() (*InMemoryDatabase, error)
	Accept(blockNumber uint64, blockTimestamp uint64)
	SetOrderStatus(orderId common.Hash, status Status, info string, blockNumber uint64) error
	RevertLastStatus(orderId common.Hash) error
	GetNaughtyTraders(oraclePrices map[Market]*big.Int, markets []Market) ([]LiquidablePosition, map[common.Address][]Order)
	GetAllOpenOrdersForTrader(trader common.Address) []Order
	GetOpenOrdersForTraderByType(trader common.Address, orderType OrderType) []Order
	UpdateLastPremiumFraction(market Market, trader common.Address, lastPremiumFraction *big.Int, cumlastPremiumFraction *big.Int)
	GetOrderById(orderId common.Hash) *Order
	GetTraderInfo(trader common.Address) *Trader
}

type Snapshot struct {
	Data                *InMemoryDatabase
	AcceptedBlockNumber *big.Int // data includes this block number too
}

func (db *InMemoryDatabase) LoadFromSnapshot(snapshot Snapshot) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if snapshot.Data == nil || snapshot.Data.OrderMap == nil || snapshot.Data.TraderMap == nil || snapshot.Data.LastPrice == nil {
		return fmt.Errorf("invalid snapshot; snapshot=%+v", snapshot)
	}

	db.OrderMap = snapshot.Data.OrderMap
	db.TraderMap = snapshot.Data.TraderMap
	db.LastPrice = snapshot.Data.LastPrice
	db.NextFundingTime = snapshot.Data.NextFundingTime

	return nil
}

// assumes that lock is held by the caller
func (db *InMemoryDatabase) Accept(blockNumber uint64, blockTimestamp uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for orderId, order := range db.OrderMap {
		lifecycle := order.getOrderStatus()
		if (lifecycle.Status == FulFilled || lifecycle.Status == Cancelled) && lifecycle.BlockNumber <= blockNumber {
			delete(db.OrderMap, orderId)
			continue
		}
		expireAt := order.getExpireAt()
		if expireAt.Sign() > 0 && expireAt.Int64() < int64(blockTimestamp) {
			delete(db.OrderMap, orderId)
		}

	}
}

func (db *InMemoryDatabase) SetOrderStatus(orderId common.Hash, status Status, info string, blockNumber uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.OrderMap[orderId] == nil {
		return fmt.Errorf("invalid orderId %s", orderId.Hex())
	}
	db.OrderMap[orderId].LifecycleList = append(db.OrderMap[orderId].LifecycleList, Lifecycle{blockNumber, status, info})
	return nil
}

func (db *InMemoryDatabase) RevertLastStatus(orderId common.Hash) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.OrderMap[orderId] == nil {
		return fmt.Errorf("invalid orderId %s", orderId.Hex())
	}

	lifeCycleList := db.OrderMap[orderId].LifecycleList
	if len(lifeCycleList) > 0 {
		db.OrderMap[orderId].LifecycleList = lifeCycleList[:len(lifeCycleList)-1]
	}
	return nil
}

func (db *InMemoryDatabase) GetAllOrders() []Order {
	db.mu.RLock() // only read lock required
	defer db.mu.RUnlock()

	allOrders := []Order{}
	for _, order := range db.OrderMap {
		allOrders = append(allOrders, deepCopyOrder(order))
	}
	return allOrders
}

func (db *InMemoryDatabase) Add(order *Order) {
	db.mu.Lock()
	defer db.mu.Unlock()

	order.LifecycleList = append(order.LifecycleList, Lifecycle{order.BlockNumber.Uint64(), Placed, ""})
	db.OrderMap[order.Id] = order
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
	if limitOrder == nil {
		log.Error("In UpdateFilledBaseAssetQuantity - orderId does not exist in the database", "orderId", orderId.Hex())
		metrics.GetOrRegisterCounter("update_filled_base_asset_quantity_order_id_not_found", nil).Inc(1)
		return
	}
	if limitOrder.PositionType == LONG {
		limitOrder.FilledBaseAssetQuantity.Add(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled + quantity
	}
	if limitOrder.PositionType == SHORT {
		limitOrder.FilledBaseAssetQuantity.Sub(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled - quantity
	}

	if limitOrder.BaseAssetQuantity.Cmp(limitOrder.FilledBaseAssetQuantity) == 0 {
		limitOrder.LifecycleList = append(limitOrder.LifecycleList, Lifecycle{blockNumber, FulFilled, ""})
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

func (db *InMemoryDatabase) GetLongOrders(market Market, lowerbound *big.Int, blockNumber *big.Int) []Order {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var longOrders []Order
	for _, order := range db.OrderMap {
		if order.PositionType == LONG && order.Market == market && (lowerbound == nil || order.Price.Cmp(lowerbound) >= 0) {
			if _order := db.getCleanOrder(order, blockNumber); _order != nil {
				longOrders = append(longOrders, *_order)
			}
		}
	}
	sortLongOrders(longOrders)
	return longOrders
}

func (db *InMemoryDatabase) GetShortOrders(market Market, upperbound *big.Int, blockNumber *big.Int) []Order {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var shortOrders []Order
	for _, order := range db.OrderMap {
		if order.PositionType == SHORT && order.Market == market && (upperbound == nil || order.Price.Cmp(upperbound) <= 0) {
			if _order := db.getCleanOrder(order, blockNumber); _order != nil {
				shortOrders = append(shortOrders, *_order)
			}
		}
	}
	sortShortOrders(shortOrders)
	return shortOrders
}

func (db *InMemoryDatabase) getCleanOrder(order *Order, blockNumber *big.Int) *Order {
	eligibleForExecution := false
	orderStatus := order.getOrderStatus()
	switch orderStatus.Status {
	case Placed:
		eligibleForExecution = true
	case Execution_Failed:
		// ideally these orders should have been auto-cancelled (by the validator) at the same time that they were fulfilling the criteria to fail
		// However, there are several reasons why this might not have happened
		// 1. A particular cancellation strategy is not implemented yet for e.g. reduce only orders with order.BaseAssetQuantity > position.size are not being auto-cancelled as of Jun 20, 23. This is a @todo
		// 2. There might be a scenarios that the order was not deemed cancellable at the time of checking and was hence used for matching; but then eventually failed execution
		//		a. a tx before the order in the same block, changed their PnL which caused them to have insufficient margin to execute the order
		//		b. specially true in multi-collateral, where the price of 1 collateral dipped but recovered again after the order was taken for matching (but failed execution)
		// 3. There might be a bug in the cancellation logic in either of EVM or smart contract code
		// 4. We might have made margin requirements for order fulfillment more liberal at a later stage
		// Hence, in view of the above and to serve as a catch-all we retry failed orders after every 100 blocks
		// Note at if an order is failing multiple times and it is also not being caught in the auto-cancel logic, then something/somewhere definitely needs fixing
		if blockNumber != nil {
			if orderStatus.BlockNumber+RETRY_AFTER_BLOCKS <= blockNumber.Uint64() {
				eligibleForExecution = true
			} else if blockNumber.Uint64()%10 == 0 {
				// to not make the log too noisy
				log.Warn("eligible order is in Execution_Failed state", "orderId", order.String(), "retryInBlocks", orderStatus.BlockNumber+RETRY_AFTER_BLOCKS-blockNumber.Uint64())
			}
		}
	}

	expireAt := order.getExpireAt()
	if expireAt.Sign() == 1 && expireAt.Int64() <= time.Now().Unix() {
		eligibleForExecution = false
	}

	if eligibleForExecution {
		if order.ReduceOnly {
			return db.getReduceOnlyOrderDisplay(order)
		}
		_order := deepCopyOrder(order)
		return &_order
	}
	return nil
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
}

func (db *InMemoryDatabase) UpdateReservedMargin(trader common.Address, addAmount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = getBlankTrader()
	}

	db.TraderMap[trader].Margin.Reserved.Add(db.TraderMap[trader].Margin.Reserved, addAmount)
}

func (db *InMemoryDatabase) UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool, blockNumber uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = getBlankTrader()
	}

	if _, ok := db.TraderMap[trader].Positions[market]; !ok {
		db.TraderMap[trader].Positions[market] = &Position{}
	}

	if db.CumulativePremiumFraction[market] == nil {
		db.CumulativePremiumFraction[market] = big.NewInt(0)
	}

	previousSize := db.TraderMap[trader].Positions[market].Size
	if previousSize == nil || previousSize.Sign() == 0 {
		// this is also set in the AMM contract when a new position is opened, without emitting a FundingPaid event
		db.TraderMap[trader].Positions[market].LastPremiumFraction = db.CumulativePremiumFraction[market]
		db.TraderMap[trader].Positions[market].UnrealisedFunding = big.NewInt(0)
	}

	// before the rc9 release (completed at block 1530589) hubble-protocol, it was possible that the lastPremiumFraction for a trader was updated without emitting a corresponding event.
	// This only happened in markets for which trader had a 0 position.
	// Since we build the entire memory db state based on events alone, we miss these updates and hence "forcibly" set LastPremiumFraction = CumulativePremiumFraction for a trader in all markets
	// note that in rc9 release this was changed and the "FundingPaid" event will always be emitted whenever the lastPremiumFraction is updated (EXCEPT for the case when trader opens a new position in the market - handled above)
	// so while we still need this update for backwards compatibility, it can be removed when there is a fresh deployment of the entire system.
	if blockNumber <= 1530589 {
		for market, position := range db.TraderMap[trader].Positions {
			if db.CumulativePremiumFraction[market] == nil {
				db.CumulativePremiumFraction[market] = big.NewInt(0)
			}
			if position.LastPremiumFraction == nil {
				position.LastPremiumFraction = big.NewInt(0)
			}
			if position.LastPremiumFraction.Cmp(db.CumulativePremiumFraction[market]) != 0 {
				if position.Size == nil || position.Size.Sign() == 0 || calcPendingFunding(db.CumulativePremiumFraction[market], position.LastPremiumFraction, position.Size).Sign() == 0 {
					// expected scenario
					position.LastPremiumFraction = db.CumulativePremiumFraction[market]
				} else {
					log.Error("pendingFunding is not 0", "trader", trader.String(), "market", market, "position", position, "pendingFunding", calcPendingFunding(db.CumulativePremiumFraction[market], position.LastPremiumFraction, position.Size), "lastPremiumFraction", position.LastPremiumFraction, "cumulativePremiumFraction", db.CumulativePremiumFraction[market])
				}
			}
		}
	}

	db.TraderMap[trader].Positions[market].Size = size
	db.TraderMap[trader].Positions[market].OpenNotional = openNotional

	if !isLiquidation {
		db.TraderMap[trader].Positions[market].LiquidationThreshold = getLiquidationThreshold(db.configService.getMaxLiquidationRatio(market), db.configService.getMinSizeRequirement(market), size)
	}

	// adjust the liquidation threshold if > resultant position size (for both isLiquidation = true/false)
	threshold := utils.BigIntMinAbs(db.TraderMap[trader].Positions[market].LiquidationThreshold, size)
	db.TraderMap[trader].Positions[market].LiquidationThreshold.Mul(threshold, big.NewInt(int64(size.Sign()))) // same sign as size
}

func (db *InMemoryDatabase) UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.CumulativePremiumFraction[market] = cumulativePremiumFraction
	for _, trader := range db.TraderMap {
		position := trader.Positions[market]
		if position != nil {
			position.UnrealisedFunding = calcPendingFunding(cumulativePremiumFraction, position.LastPremiumFraction, position.Size)
		}
	}
}

func calcPendingFunding(cumulativePremiumFraction, lastPremiumFraction, size *big.Int) *big.Int {
	if size == nil || size.Sign() == 0 {
		return big.NewInt(0)
	}

	if cumulativePremiumFraction == nil {
		cumulativePremiumFraction = big.NewInt(0)
	}

	if lastPremiumFraction == nil {
		lastPremiumFraction = big.NewInt(0)
	}

	// Calculate difference
	diff := new(big.Int).Sub(cumulativePremiumFraction, lastPremiumFraction)

	// Multiply by size
	result := new(big.Int).Mul(diff, size)

	// Handle negative rounding
	if result.Sign() < 0 {
		result.Add(result, big.NewInt(1e18-1))
	}

	// Divide by 1e18
	return result.Div(result, SIZE_BASE_PRECISION)
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

func (db *InMemoryDatabase) UpdateLastPremiumFraction(market Market, trader common.Address, lastPremiumFraction *big.Int, cumulativePremiumFraction *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.TraderMap[trader]; !ok {
		db.TraderMap[trader] = getBlankTrader()
	}

	if _, ok := db.TraderMap[trader].Positions[market]; !ok {
		db.TraderMap[trader].Positions[market] = &Position{}
	}

	db.TraderMap[trader].Positions[market].LastPremiumFraction = lastPremiumFraction
	db.TraderMap[trader].Positions[market].UnrealisedFunding = dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, lastPremiumFraction), db.TraderMap[trader].Positions[market].Size))
}

func (db *InMemoryDatabase) GetLastPrice(market Market) *big.Int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return big.NewInt(0).Set(db.LastPrice[market])
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

func (db *InMemoryDatabase) GetOpenOrdersForTraderByType(trader common.Address, orderType OrderType) []Order {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.getTraderOrders(trader, orderType)
}

func (db *InMemoryDatabase) GetAllOpenOrdersForTrader(trader common.Address) []Order {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.getAllTraderOrders(trader)
}

func (db *InMemoryDatabase) GetOrderById(orderId common.Hash) *Order {
	db.mu.RLock()
	defer db.mu.RUnlock()

	order := db.OrderMap[orderId]
	if order == nil {
		return nil
	}

	orderCopy := deepCopyOrder(order)
	return &orderCopy
}

func (db *InMemoryDatabase) GetTraderInfo(trader common.Address) *Trader {
	db.mu.RLock()
	defer db.mu.RUnlock()

	traderInfo := db.TraderMap[trader]
	if traderInfo == nil {
		return nil
	}

	traderCopy := deepCopyTrader(traderInfo)
	return traderCopy
}

func determinePositionToLiquidate(trader *Trader, addr common.Address, marginFraction *big.Int, markets []Market, minSizes []*big.Int) LiquidablePosition {
	liquidable := LiquidablePosition{}
	// iterate through the markets and return the first one with an open position
	// @todo when we introduce multiple markets, we will have to implement a more sophisticated liquidation strategy
	for i, market := range markets {
		position := trader.Positions[market]
		if position == nil || position.Size.Sign() == 0 {
			continue
		}
		liquidable = LiquidablePosition{
			Address:        addr,
			Market:         market,
			Size:           new(big.Int).Abs(position.LiquidationThreshold), // position.LiquidationThreshold is a pointer, to want to avoid unintentional mutation if/when we mutate liquidable.Size
			MarginFraction: new(big.Int).Set(marginFraction),
			FilledSize:     big.NewInt(0),
		}
		// while setting liquidation threshold of a position, we do not ensure whether it is a multiple of minSize.
		// we will take care of that here
		liquidable.Size.Div(liquidable.Size, minSizes[i])
		liquidable.Size.Mul(liquidable.Size, minSizes[i])
		if position.Size.Sign() == -1 {
			liquidable.PositionType = SHORT
			liquidable.Size.Neg(liquidable.Size)
		} else {
			liquidable.PositionType = LONG
		}
	}
	return liquidable
}

func (db *InMemoryDatabase) GetNaughtyTraders(oraclePrices map[Market]*big.Int, markets []Market) ([]LiquidablePosition, map[common.Address][]Order) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	liquidablePositions := []LiquidablePosition{}
	ordersToCancel := map[common.Address][]Order{}
	count := 0

	// will be updated lazily only if liquidablePositions are found
	minSizes := []*big.Int{}

	for addr, trader := range db.TraderMap {
		pendingFunding := getTotalFunding(trader, markets)
		marginFraction := calcMarginFraction(trader, pendingFunding, oraclePrices, db.LastPrice, markets)
		if marginFraction.Cmp(db.configService.getMaintenanceMargin()) == -1 {
			log.Info("below maintenanceMargin", "trader", addr.String(), "marginFraction", prettifyScaledBigInt(marginFraction, 6))
			if len(minSizes) == 0 {
				for _, market := range markets {
					minSizes = append(minSizes, db.configService.getMinSizeRequirement(market))
				}
			}
			liquidablePositions = append(liquidablePositions, determinePositionToLiquidate(trader, addr, marginFraction, markets, minSizes))
			continue // we do not check for their open orders yet. Maybe liquidating them first will make available margin positive
		}
		if trader.Margin.Reserved.Sign() == 0 {
			continue
		}
		// has orders that might be cancellable
		availableMargin := getAvailableMargin(trader, pendingFunding, oraclePrices, db.LastPrice, db.configService.getMinAllowableMargin(), markets)
		if availableMargin.Cmp(big.NewInt(0)) == -1 {
			foundCancellableOrders := db.determineOrdersToCancel(addr, trader, availableMargin, oraclePrices, ordersToCancel)
			if foundCancellableOrders {
				log.Info("negative available margin", "trader", addr.String(), "availableMargin", prettifyScaledBigInt(availableMargin, 6))
			} else {
				count++
			}
		}
	}
	if count > 0 {
		log.Info("#traders that have -ve margin but no orders to cancel", "count", count)
	}
	// lower margin fraction positions should be liquidated first
	sortLiquidableSliceByMarginFraction(liquidablePositions)
	return liquidablePositions, ordersToCancel
}

// assumes db.mu.RLock has been held by the caller
func (db *InMemoryDatabase) determineOrdersToCancel(addr common.Address, trader *Trader, availableMargin *big.Int, oraclePrices map[Market]*big.Int, ordersToCancel map[common.Address][]Order) bool {
	traderOrders := db.getTraderOrders(addr, LimitOrderType)
	sort.Slice(traderOrders, func(i, j int) bool {
		// higher diff comes first
		iDiff := big.NewInt(0).Abs(big.NewInt(0).Sub(traderOrders[i].Price, oraclePrices[traderOrders[i].Market]))
		jDiff := big.NewInt(0).Abs(big.NewInt(0).Sub(traderOrders[j].Price, oraclePrices[traderOrders[j].Market]))
		return iDiff.Cmp(jDiff) > 0
	})

	_availableMargin := new(big.Int).Set(availableMargin)
	if len(traderOrders) > 0 {
		// cancel orders until available margin is positive
		ordersToCancel[addr] = []Order{}
		for _, order := range traderOrders {
			// cannot cancel ReduceOnly orders or Market orders because no margin is reserved for them
			if order.ReduceOnly || order.OrderType != LimitOrderType {
				continue
			}
			ordersToCancel[addr] = append(ordersToCancel[addr], order)
			orderNotional := big.NewInt(0).Abs(big.NewInt(0).Div(big.NewInt(0).Mul(order.GetUnFilledBaseAssetQuantity(), order.Price), _1e18)) // | size * current price |
			marginReleased := divideByBasePrecision(big.NewInt(0).Mul(orderNotional, db.configService.getMinAllowableMargin()))
			_availableMargin.Add(_availableMargin, marginReleased)
			if _availableMargin.Sign() >= 0 {
				break
			}
		}
		return true
	}
	return false
}

func (db *InMemoryDatabase) getTraderOrders(trader common.Address, orderType OrderType) []Order {
	traderOrders := []Order{}
	_trader := trader.String()
	for _, order := range db.OrderMap {
		if strings.EqualFold(order.UserAddress, _trader) && order.OrderType == orderType {
			traderOrders = append(traderOrders, deepCopyOrder(order))
		}
	}
	return traderOrders
}

func (db *InMemoryDatabase) getAllTraderOrders(trader common.Address) []Order {
	traderOrders := []Order{}
	_trader := trader.String()
	for _, order := range db.OrderMap {
		if strings.EqualFold(order.UserAddress, _trader) {
			traderOrders = append(traderOrders, deepCopyOrder(order))
		}
	}
	return traderOrders
}

func (db *InMemoryDatabase) getReduceOnlyOrderDisplay(order *Order) *Order {
	trader := common.HexToAddress(order.UserAddress)
	if db.TraderMap[trader] == nil {
		return nil
	}
	positions := db.TraderMap[trader].Positions
	if position, ok := positions[order.Market]; ok {
		// position.Size, order.BaseAssetQuantity need to be of opposite sign and abs(position.Size) >= abs(order.BaseAssetQuantity)
		if position.Size.Sign() == 0 || position.Size.Sign() == order.BaseAssetQuantity.Sign() {
			return nil
		}
		if position.Size.CmpAbs(order.GetUnFilledBaseAssetQuantity()) >= 0 {
			// position is bigger than unfilled order size
			orderCopy := deepCopyOrder(order)
			return &orderCopy
		} else {
			// position is smaller than unfilled order
			// increase the filled quantity so that unfilled amount is equal to position size
			orderCopy := deepCopyOrder(order)
			orderCopy.FilledBaseAssetQuantity = big.NewInt(0).Add(orderCopy.BaseAssetQuantity, position.Size) // both have opposite sign, therefore we add
			return &orderCopy
		}
	} else {
		return nil
	}
}

func sortLongOrders(orders []Order) []Order {
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

func sortShortOrders(orders []Order) []Order {
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

func (db *InMemoryDatabase) GetOrderBookDataCopy() (*InMemoryDatabase, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(db)
	if err != nil {
		return nil, fmt.Errorf("error encoding database: %v", err)
	}

	buf2 := bytes.NewBuffer(buf.Bytes())
	var memoryDBCopy *InMemoryDatabase
	err = gob.NewDecoder(buf2).Decode(&memoryDBCopy)
	if err != nil {
		return nil, fmt.Errorf("error decoding database: %v", err)
	}

	memoryDBCopy.mu = &sync.RWMutex{}
	return memoryDBCopy, nil
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

func getAvailableMargin(trader *Trader, pendingFunding *big.Int, oraclePrices map[Market]*big.Int, lastPrices map[Market]*big.Int, minAllowableMargin *big.Int, markets []Market) *big.Int {
	// log.Info("in getAvailableMargin", "trader", trader, "pendingFunding", pendingFunding, "oraclePrices", oraclePrices, "lastPrices", lastPrices)
	margin := new(big.Int).Sub(getNormalisedMargin(trader), pendingFunding)
	notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(trader, margin, Min_Allowable_Margin, oraclePrices, lastPrices, markets)
	utilisedMargin := divideByBasePrecision(new(big.Int).Mul(notionalPosition, minAllowableMargin))
	// print margin, notionalPosition, unrealizePnL, utilisedMargin
	// log.Info("stats", "margin", margin, "notionalPosition", notionalPosition, "unrealizePnL", unrealizePnL, "utilisedMargin", utilisedMargin, "Reserved", trader.Margin.Reserved)
	return new(big.Int).Sub(
		new(big.Int).Add(margin, unrealizePnL),
		new(big.Int).Add(utilisedMargin, trader.Margin.Reserved),
	)
}

// deepCopyOrder deep copies the LimitOrder struct
func deepCopyOrder(order *Order) Order {
	lifecycleList := &order.LifecycleList
	return Order{
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
		BlockNumber:             big.NewInt(0).Set(order.BlockNumber),
		RawOrder:                order.RawOrder,
	}
}

func deepCopyTrader(order *Trader) *Trader {
	positions := map[Market]*Position{}
	for market, position := range order.Positions {
		positions[market] = &Position{
			OpenNotional:         big.NewInt(0).Set(position.OpenNotional),
			Size:                 big.NewInt(0).Set(position.Size),
			UnrealisedFunding:    big.NewInt(0).Set(position.UnrealisedFunding),
			LastPremiumFraction:  big.NewInt(0).Set(position.LastPremiumFraction),
			LiquidationThreshold: big.NewInt(0).Set(position.LiquidationThreshold),
		}
	}

	margin := Margin{
		Reserved:  big.NewInt(0).Set(order.Margin.Reserved),
		Deposited: map[Collateral]*big.Int{},
	}
	for collateral, amount := range order.Margin.Deposited {
		margin.Deposited[collateral] = big.NewInt(0).Set(amount)
	}
	return &Trader{
		Positions: positions,
		Margin:    margin,
	}
}
