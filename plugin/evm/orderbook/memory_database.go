package orderbook

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/metrics"
	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type InMemoryDatabase struct {
	mu                        *sync.RWMutex              `json:"-"`
	Orders                    map[common.Hash]*Order     `json:"order_map"` // ID => order
	LongOrders                map[Market][]*Order        `json:"long_orders"`
	ShortOrders               map[Market][]*Order        `json:"short_orders"`
	TraderMap                 map[common.Address]*Trader `json:"trader_map"` // address => trader info
	NextFundingTime           uint64                     `json:"next_funding_time"`
	LastPrice                 map[Market]*big.Int        `json:"last_price"`
	CumulativePremiumFraction map[Market]*big.Int        `json:"cumulative_last_premium_fraction"`
	NextSamplePITime          uint64                     `json:"next_sample_pi_time"`
	SamplePIAttemptedTime     uint64                     `json:"sample_pi_attempted_time"`
	configService             IConfigService
}

func NewInMemoryDatabase(configService IConfigService) *InMemoryDatabase {
	lastPrice := map[Market]*big.Int{}
	traderMap := map[common.Address]*Trader{}

	return &InMemoryDatabase{
		Orders:                    map[common.Hash]*Order{},
		LongOrders:                map[Market][]*Order{},
		ShortOrders:               map[Market][]*Order{},
		TraderMap:                 traderMap,
		LastPrice:                 lastPrice,
		CumulativePremiumFraction: map[Market]*big.Int{},
		mu:                        &sync.RWMutex{},
		configService:             configService,
	}
}

const (
	RETRY_AFTER_BLOCKS = 10
)

type Market = hu.Market

type Collateral = int

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

type OrderType = hu.OrderType

const (
	Limit  = hu.Limit
	IOC    = hu.IOC
	Signed = hu.Signed
)

type Lifecycle struct {
	BlockNumber uint64
	Status      Status
	Info        string
}

type Order struct {
	Id                      common.Hash
	Market                  Market
	PositionType            PositionType
	Trader                  common.Address
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
		Trader                  string      `json:"trader"`
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
		Trader:                  order.Trader.String(),
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
	if order.OrderType == IOC {
		return order.RawOrder.(*IOCOrder).ExpireAt
	}
	if order.OrderType == Signed {
		return order.RawOrder.(*hu.SignedOrder).ExpireAt
	}
	return big.NewInt(0)
}

func (order Order) isPostOnly() bool {
	if order.OrderType == Limit {
		if rawOrder, ok := order.RawOrder.(*LimitOrder); ok {
			return rawOrder.PostOnly
		}
	}
	if order.OrderType == Signed {
		if rawOrder, ok := order.RawOrder.(*hu.SignedOrder); ok {
			return rawOrder.PostOnly
		}
	}
	return false
}

func (order Order) String() string {
	t := time.Unix(order.getExpireAt().Int64(), 0)
	return fmt.Sprintf("Order: Id: %s, OrderType: %s, Market: %v, PositionType: %v, UserAddress: %v, BaseAssetQuantity: %s, FilledBaseAssetQuantity: %s, Salt: %v, Price: %s, ReduceOnly: %v, PostOnly: %v, expireAt %s, BlockNumber: %s", order.Id, order.OrderType, order.Market, order.PositionType, order.Trader.String(), prettifyScaledBigInt(order.BaseAssetQuantity, 18), prettifyScaledBigInt(order.FilledBaseAssetQuantity, 18), order.Salt, prettifyScaledBigInt(order.Price, 6), order.ReduceOnly, order.isPostOnly(), t.UTC(), order.BlockNumber)
}

func (order Order) ToOrderMin() OrderMin {
	return OrderMin{
		Market:  order.Market,
		Price:   order.Price.String(),
		Size:    order.GetUnFilledBaseAssetQuantity().String(),
		Signer:  order.Trader.String(),
		OrderId: order.Id.String(),
	}
}

type Position struct {
	hu.Position
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
	GetMarketOrders(market Market) []Order
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
	UpdateNextSamplePITime(nextSamplePITime uint64)
	GetNextSamplePITime() uint64
	GetSamplePIAttemptedTime() uint64
	SignalSamplePIAttempted(time uint64)
	UpdateLastPrice(market Market, lastPrice *big.Int)
	GetLastPrices() map[Market]*big.Int
	GetAllTraders() map[common.Address]Trader
	GetOrderBookData() InMemoryDatabase
	GetOrderBookDataCopy() (*InMemoryDatabase, error)
	Accept(acceptedBlockNumber uint64, blockTimestamp uint64)
	SetOrderStatus(orderId common.Hash, status Status, info string, blockNumber uint64) error
	RevertLastStatus(orderId common.Hash) error
	GetNaughtyTraders(hState *hu.HubbleState) ([]LiquidablePosition, map[common.Address][]Order, map[common.Address]*big.Int)
	GetAllOpenOrdersForTrader(trader common.Address) []Order
	GetOpenOrdersForTraderByType(trader common.Address, orderType OrderType) []Order
	UpdateLastPremiumFraction(market Market, trader common.Address, lastPremiumFraction *big.Int, cumlastPremiumFraction *big.Int)
	GetOrderById(orderId common.Hash) *Order
	GetTraderInfo(trader common.Address) *Trader
	GetOrderValidationFields(
		orderId common.Hash,
		trader common.Address,
		marketId int,
	) OrderValidationFields
}

type Snapshot struct {
	Data                *InMemoryDatabase
	AcceptedBlockNumber *big.Int // data includes this block number too
}

func (db *InMemoryDatabase) LoadFromSnapshot(snapshot Snapshot) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if snapshot.Data == nil || snapshot.Data.Orders == nil || snapshot.Data.TraderMap == nil || snapshot.Data.LastPrice == nil ||
		snapshot.Data.CumulativePremiumFraction == nil {
		return fmt.Errorf("invalid snapshot; snapshot=%+v", snapshot)
	}

	db.Orders = snapshot.Data.Orders
	db.TraderMap = snapshot.Data.TraderMap
	db.LastPrice = snapshot.Data.LastPrice
	db.NextFundingTime = snapshot.Data.NextFundingTime
	db.NextSamplePITime = snapshot.Data.NextSamplePITime
	db.CumulativePremiumFraction = snapshot.Data.CumulativePremiumFraction

	for _, order := range db.Orders {
		db.AddInSortedArray(order)
	}
	return nil
}

func (db *InMemoryDatabase) Accept(acceptedBlockNumber, blockTimestamp uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	log.Info("Accept", "acceptedBlockNumber", acceptedBlockNumber, "blockTimestamp", blockTimestamp)
	count := db.configService.GetActiveMarketsCount()
	for m := int64(0); m < count; m++ {
		longOrders := db.getLongOrdersWithoutLock(Market(m), nil, nil, false)
		shortOrders := db.getShortOrdersWithoutLock(Market(m), nil, nil, false)

		for _, longOrder := range longOrders {
			status := shouldRemove(acceptedBlockNumber, blockTimestamp, longOrder)
			log.Info("evaluating order...", "longOrder", longOrder, "status", status)
			if status == KEEP_IF_MATCHEABLE {
				matchFound := false
				for _, shortOrder := range shortOrders {
					if longOrder.Price.Cmp(shortOrder.Price) < 0 {
						break // because the short orders are sorted in ascending order of price, there is no point in checking further
					}
					// an IOC order even if has a price overlap can only be matched if the order came before it (or same block)
					if longOrder.BlockNumber.Uint64() >= shortOrder.BlockNumber.Uint64() {
						matchFound = true
						break
					} /* else {
						dont break here because there might be an a short order with higher price that came before the IOC longOrder in question
					} */
				}
				if !matchFound {
					status = REMOVE
				}
			}

			if status == REMOVE {
				db.deleteOrderWithoutLock(longOrder.Id)
			}
		}

		for _, shortOrder := range shortOrders {
			status := shouldRemove(acceptedBlockNumber, blockTimestamp, shortOrder)
			log.Info("Accept", "shortOrder", shortOrder, "status", status)
			if status == KEEP_IF_MATCHEABLE {
				matchFound := false
				for _, longOrder := range longOrders {
					if longOrder.Price.Cmp(shortOrder.Price) < 0 {
						break // because the long orders are sorted in descending order of price, there is no point in checking further
					}
					// an IOC order even if has a price overlap can only be matched if the order came before it (or same block)
					if shortOrder.BlockNumber.Uint64() >= longOrder.BlockNumber.Uint64() {
						matchFound = true
						break
					}
					/* else {
						dont break here because there might be an a long order with lower price that came before the IOC shortOrder in question
					} */
				}
				if !matchFound {
					status = REMOVE
				}
			}

			if status == REMOVE {
				db.deleteOrderWithoutLock(shortOrder.Id)
			}
		}
	}
}

type OrderStatus uint8

const (
	KEEP OrderStatus = iota
	REMOVE
	KEEP_IF_MATCHEABLE
)

func shouldRemove(acceptedBlockNumber, blockTimestamp uint64, order Order) OrderStatus {
	// check if there is any criteria to delete the order
	// 1. Order is fulfilled or cancelled
	lifecycle := order.getOrderStatus()
	if (lifecycle.Status == FulFilled || lifecycle.Status == Cancelled) && lifecycle.BlockNumber <= acceptedBlockNumber {
		return REMOVE
	}

	if order.OrderType == Limit {
		return KEEP
	}

	// remove if order is expired; valid for both IOC and Signed orders
	expireAt := order.getExpireAt()
	if expireAt.Sign() > 0 && expireAt.Int64() < int64(blockTimestamp) {
		return REMOVE
	}

	// IOC order can not matched with any order that came after it (same block is allowed)
	// we can only surely say about orders that came at <= acceptedBlockNumber
	if order.OrderType == IOC {
		if order.BlockNumber.Uint64() > acceptedBlockNumber {
			return KEEP
		}
		return KEEP_IF_MATCHEABLE
	}
	return KEEP
}

func (db *InMemoryDatabase) SetOrderStatus(orderId common.Hash, status Status, info string, blockNumber uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.Orders[orderId] == nil {
		return fmt.Errorf("invalid orderId %s", orderId.Hex())
	}
	db.Orders[orderId].LifecycleList = append(db.Orders[orderId].LifecycleList, Lifecycle{blockNumber, status, info})
	return nil
}

func (db *InMemoryDatabase) RevertLastStatus(orderId common.Hash) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.Orders[orderId] == nil {
		return fmt.Errorf("invalid orderId %s", orderId.Hex())
	}

	lifeCycleList := db.Orders[orderId].LifecycleList
	if len(lifeCycleList) > 0 {
		db.Orders[orderId].LifecycleList = lifeCycleList[:len(lifeCycleList)-1]
	}
	return nil
}

func (db *InMemoryDatabase) GetAllOrders() []Order {
	db.mu.RLock() // only read lock required
	defer db.mu.RUnlock()

	allOrders := []Order{}
	for _, order := range db.Orders {
		allOrders = append(allOrders, deepCopyOrder(order))
	}
	return allOrders
}

func (db *InMemoryDatabase) GetMarketOrders(market Market) []Order {
	db.mu.RLock() // only read lock required
	defer db.mu.RUnlock()

	allOrders := []Order{}
	for _, order := range db.LongOrders[market] {
		allOrders = append(allOrders, deepCopyOrder(order))
	}

	for _, order := range db.ShortOrders[market] {
		allOrders = append(allOrders, deepCopyOrder(order))
	}

	return allOrders
}

func (db *InMemoryDatabase) Add(order *Order) {
	db.mu.Lock()
	defer db.mu.Unlock()

	log.Info("Adding order to memdb", "order", order)
	order.LifecycleList = append(order.LifecycleList, Lifecycle{order.BlockNumber.Uint64(), Placed, ""})
	db.AddInSortedArray(order)
	db.Orders[order.Id] = order
}

// caller is expected to acquire db.mu before calling this function
func (db *InMemoryDatabase) AddInSortedArray(order *Order) {
	market := order.Market

	var orders []*Order
	var position int
	if order.PositionType == LONG {
		orders = db.LongOrders[market]
		position = sort.Search(len(orders), func(i int) bool {
			priceDiff := order.Price.Cmp(orders[i].Price)
			if priceDiff == 1 {
				return true
			} else if priceDiff == 0 {
				blockDiff := order.BlockNumber.Cmp(orders[i].BlockNumber)
				if blockDiff == -1 { // order was placed before i
					return true
				} else if blockDiff == 0 { // order and i were placed in the same block
					if order.OrderType == IOC {
						// prioritize fulfilling IOC orders first, because they are short-lived
						return true
					}
				}
			}
			return false
		})
	} else {
		orders = db.ShortOrders[market]
		position = sort.Search(len(orders), func(i int) bool {
			priceDiff := order.Price.Cmp(orders[i].Price)
			if priceDiff == -1 {
				return true
			} else if priceDiff == 0 {
				blockDiff := order.BlockNumber.Cmp(orders[i].BlockNumber)
				if blockDiff == -1 { // order was placed before i
					return true
				} else if blockDiff == 0 { // order and i were placed in the same block
					if order.OrderType == IOC {
						// prioritize fulfilling IOC orders first, because they are short-lived
						return true
					}
				}
			}
			return false
		})
	}

	// Insert the order at the determined position
	orders = append(orders, &Order{})            // Add an empty order to the end
	copy(orders[position+1:], orders[position:]) // Shift orders to the right
	orders[position] = order                     // Insert new Order at the right position

	if order.PositionType == LONG {
		db.LongOrders[market] = orders
	} else {
		db.ShortOrders[market] = orders
	}
}

func (db *InMemoryDatabase) Delete(orderId common.Hash) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.deleteOrderWithoutLock(orderId)
}

func (db *InMemoryDatabase) deleteOrderWithoutLock(orderId common.Hash) {
	order := db.Orders[orderId]
	if order == nil {
		log.Error("In Delete - orderId does not exist in the db.Orders", "orderId", orderId.Hex())
		deleteOrderIdNotFoundCounter.Inc(1)
		return
	}

	market := order.Market
	if order.PositionType == LONG {
		orders := db.LongOrders[market]
		idx := getOrderIdx(orders, orderId)
		if idx == -1 {
			log.Error("In Delete - orderId does not exist in the db.LongOrders", "orderId", orderId.Hex())
			deleteOrderIdNotFoundCounter.Inc(1)
		} else {
			orders = append(orders[:idx], orders[idx+1:]...)
			db.LongOrders[market] = orders
		}
	} else {
		orders := db.ShortOrders[market]
		idx := getOrderIdx(orders, orderId)
		if idx == -1 {
			log.Error("In Delete - orderId does not exist in the db.ShortOrders", "orderId", orderId.Hex())
			deleteOrderIdNotFoundCounter.Inc(1)
		} else {
			orders = append(orders[:idx], orders[idx+1:]...)
			db.ShortOrders[market] = orders
		}
	}

	delete(db.Orders, orderId)
}

func (db *InMemoryDatabase) UpdateFilledBaseAssetQuantity(quantity *big.Int, orderId common.Hash, blockNumber uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	order := db.Orders[orderId]
	if order == nil {
		log.Error("In UpdateFilledBaseAssetQuantity - orderId does not exist in the database", "orderId", orderId.Hex())
		metrics.GetOrRegisterCounter("update_filled_base_asset_quantity_order_id_not_found", nil).Inc(1)
		return
	}
	if order.PositionType == LONG {
		order.FilledBaseAssetQuantity.Add(order.FilledBaseAssetQuantity, quantity) // filled = filled + quantity
	}
	if order.PositionType == SHORT {
		order.FilledBaseAssetQuantity.Sub(order.FilledBaseAssetQuantity, quantity) // filled = filled - quantity
	}

	if order.BaseAssetQuantity.Cmp(order.FilledBaseAssetQuantity) == 0 {
		order.LifecycleList = append(order.LifecycleList, Lifecycle{blockNumber, FulFilled, ""})
	}

	if quantity.Cmp(big.NewInt(0)) == -1 && order.getOrderStatus().Status == FulFilled {
		// handling reorgs
		order.LifecycleList = order.LifecycleList[:len(order.LifecycleList)-1]
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

func (db *InMemoryDatabase) GetNextSamplePITime() uint64 {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.NextSamplePITime
}

func (db *InMemoryDatabase) GetSamplePIAttemptedTime() uint64 {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.SamplePIAttemptedTime
}

func (db *InMemoryDatabase) SignalSamplePIAttempted(time uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.SamplePIAttemptedTime = time
}

func (db *InMemoryDatabase) UpdateNextSamplePITime(nextSamplePITime uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.NextSamplePITime = nextSamplePITime
}

func (db *InMemoryDatabase) GetLongOrders(market Market, lowerbound *big.Int, blockNumber *big.Int) []Order {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.getLongOrdersWithoutLock(market, lowerbound, blockNumber, true)
}

func (db *InMemoryDatabase) getLongOrdersWithoutLock(market Market, lowerbound *big.Int, blockNumber *big.Int, shouldClean bool) []Order {
	var longOrders []Order

	marketOrders := db.LongOrders[market]
	// log.Info("getLongOrdersWithoutLock", "marketOrders", marketOrders, "lowerbound", lowerbound, "blockNumber", blockNumber)
	for _, order := range marketOrders {
		if lowerbound != nil && order.Price.Cmp(lowerbound) < 0 {
			// because the long orders are sorted in descending order of price, there is no point in checking further
			break
		}

		if shouldClean {
			if _order := db.getCleanOrder(order, blockNumber); _order != nil {
				longOrders = append(longOrders, *_order)
			}
		} else {
			longOrders = append(longOrders, deepCopyOrder(order))
		}
	}
	return longOrders
}

func (db *InMemoryDatabase) GetShortOrders(market Market, upperbound *big.Int, blockNumber *big.Int) []Order {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.getShortOrdersWithoutLock(market, upperbound, blockNumber, true)
}

func (db *InMemoryDatabase) getShortOrdersWithoutLock(market Market, upperbound *big.Int, blockNumber *big.Int, shouldClean bool) []Order {
	var shortOrders []Order

	marketOrders := db.ShortOrders[market]

	for _, order := range marketOrders {
		if upperbound != nil && order.Price.Cmp(upperbound) > 0 {
			// short orders are sorted in ascending order of price
			break
		}
		if shouldClean {
			if _order := db.getCleanOrder(order, blockNumber); _order != nil {
				shortOrders = append(shortOrders, *_order)
			}
		} else {
			shortOrders = append(shortOrders, deepCopyOrder(order))
		}
	}
	return shortOrders
}

func (db *InMemoryDatabase) getCleanOrder(order *Order, blockNumber *big.Int) *Order {
	// log.Info("getCleanOrder", "order", order, "blockNumber", blockNumber)
	eligibleForExecution := false
	orderStatus := order.getOrderStatus()
	// log.Info("getCleanOrder", "orderStatus", orderStatus)
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
	// log.Info("getCleanOrder", "expireAt", expireAt, "eligibleForExecution", eligibleForExecution)

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
	return hu.Div1e18(result)
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
	db.TraderMap[trader].Positions[market].UnrealisedFunding = hu.Div1e18(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, lastPremiumFraction), db.TraderMap[trader].Positions[market].Size))
}

func (db *InMemoryDatabase) GetLastPrices() map[Market]*big.Int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	copyMap := make(map[Market]*big.Int)
	for k, v := range db.LastPrice {
		copyMap[k] = new(big.Int).Set(v)
	}
	return copyMap
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

	order := db.Orders[orderId]
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

func (db *InMemoryDatabase) GetNaughtyTraders(hState *hu.HubbleState) ([]LiquidablePosition, map[common.Address][]Order, map[common.Address]*big.Int) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	liquidablePositions := []LiquidablePosition{}
	ordersToCancel := map[common.Address][]Order{}
	marginMap := map[common.Address]*big.Int{}
	count := 0

	// will be updated lazily only if liquidablePositions are found
	minSizes := []*big.Int{}

	for addr, trader := range db.TraderMap {
		userState := &hu.UserState{
			Positions:      translatePositions(trader.Positions),
			Margins:        getMargins(trader, len(hState.Assets)),
			PendingFunding: getTotalFunding(trader, hState.ActiveMarkets),
			ReservedMargin: new(big.Int).Set(trader.Margin.Reserved),
		}
		marginFraction := hu.GetMarginFraction(hState, userState)
		marginMap[addr] = hu.GetAvailableMargin(hState, userState)
		if marginFraction.Cmp(hState.MaintenanceMargin) == -1 {
			log.Info("below maintenanceMargin", "trader", addr.String(), "marginFraction", prettifyScaledBigInt(marginFraction, 6))
			if len(minSizes) == 0 {
				for _, market := range hState.ActiveMarkets {
					minSizes = append(minSizes, db.configService.getMinSizeRequirement(market))
				}
			}
			liquidablePositions = append(liquidablePositions, determinePositionToLiquidate(trader, addr, marginFraction, hState.ActiveMarkets, minSizes))
			continue // we do not check for their open orders yet. Maybe liquidating them first will make available margin positive
		}
		if trader.Margin.Reserved.Sign() == 0 {
			continue
		}
		// has orders that might be cancellable
		availableMargin := new(big.Int).Set(marginMap[addr])
		if availableMargin.Sign() == -1 {
			foundCancellableOrders := false
			foundCancellableOrders = db.determineOrdersToCancel(addr, trader, availableMargin, hState.OraclePrices, ordersToCancel, hState.MinAllowableMargin)
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
	return liquidablePositions, ordersToCancel, marginMap
}

// assumes db.mu.RLock has been held by the caller
func (db *InMemoryDatabase) determineOrdersToCancel(addr common.Address, trader *Trader, availableMargin *big.Int, oraclePrices map[Market]*big.Int, ordersToCancel map[common.Address][]Order, minAllowableMargin *big.Int) bool {
	traderOrders := db.getTraderOrders(addr, Limit)
	if len(traderOrders) == 0 {
		return false
	}

	sort.Slice(traderOrders, func(i, j int) bool {
		// higher diff comes first
		iDiff := big.NewInt(0).Abs(big.NewInt(0).Sub(traderOrders[i].Price, oraclePrices[traderOrders[i].Market]))
		jDiff := big.NewInt(0).Abs(big.NewInt(0).Sub(traderOrders[j].Price, oraclePrices[traderOrders[j].Market]))
		return iDiff.Cmp(jDiff) > 0
	})

	_availableMargin := new(big.Int).Set(availableMargin)
	// cancel orders until available margin is positive
	ordersToCancel[addr] = []Order{}
	foundCancellableOrders := false
	for _, order := range traderOrders {
		// @todo how are reduce only orders that are not fillable cancelled?
		if order.ReduceOnly || order.OrderType != Limit {
			continue
		}
		ordersToCancel[addr] = append(ordersToCancel[addr], order)
		foundCancellableOrders = true
		orderNotional := big.NewInt(0).Abs(hu.Div1e18(hu.Mul(order.GetUnFilledBaseAssetQuantity(), order.Price))) // | size * current price |
		marginReleased := hu.Div1e6(hu.Mul(orderNotional, db.configService.getMinAllowableMargin()))
		_availableMargin.Add(_availableMargin, marginReleased)
		if _availableMargin.Sign() >= 0 {
			break
		}
	}
	return foundCancellableOrders
}

func (db *InMemoryDatabase) getTraderOrders(trader common.Address, orderType OrderType) []Order {
	traderOrders := []Order{}
	for _, order := range db.Orders {
		if order.Trader == trader && order.OrderType == orderType {
			traderOrders = append(traderOrders, deepCopyOrder(order))
		}
	}
	return traderOrders
}

func (db *InMemoryDatabase) getAllTraderOrders(trader common.Address) []Order {
	traderOrders := []Order{}
	for _, order := range db.Orders {
		if order.Trader == trader {
			traderOrders = append(traderOrders, deepCopyOrder(order))
		}
	}
	return traderOrders
}

func (db *InMemoryDatabase) getReduceOnlyOrderDisplay(order *Order) *Order {
	trader := order.Trader
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
	maxLiquidationSize := hu.Div1e6(big.NewInt(0).Mul(absSize, maxLiquidationRatio))
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

func getAvailableMargin(trader *Trader, hState *hu.HubbleState) *big.Int {
	return hu.GetAvailableMargin(
		hState,
		&hu.UserState{
			Positions:      translatePositions(trader.Positions),
			Margins:        getMargins(trader, len(hState.Assets)),
			PendingFunding: getTotalFunding(trader, hState.ActiveMarkets),
			ReservedMargin: trader.Margin.Reserved,
		},
	)
}

// deepCopyOrder deep copies the LimitOrder struct
func deepCopyOrder(order *Order) Order {
	lifecycleList := &order.LifecycleList
	return Order{
		Id:                      order.Id,
		Market:                  order.Market,
		PositionType:            order.PositionType,
		Trader:                  order.Trader,
		BaseAssetQuantity:       big.NewInt(0).Set(order.BaseAssetQuantity),
		FilledBaseAssetQuantity: big.NewInt(0).Set(order.FilledBaseAssetQuantity),
		Salt:                    big.NewInt(0).Set(order.Salt),
		Price:                   big.NewInt(0).Set(order.Price),
		ReduceOnly:              order.ReduceOnly,
		LifecycleList:           *lifecycleList,
		BlockNumber:             big.NewInt(0).Set(order.BlockNumber),
		RawOrder:                order.RawOrder,
		OrderType:               order.OrderType,
	}
}

func deepCopyTrader(order *Trader) *Trader {
	positions := map[Market]*Position{}
	for market, position := range order.Positions {
		positions[market] = &Position{
			Position: hu.Position{
				OpenNotional: big.NewInt(0).Set(position.OpenNotional),
				Size:         big.NewInt(0).Set(position.Size),
			},
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

func getOrderIdx(orders []*Order, orderId common.Hash) int {
	for i, order := range orders {
		if order.Id == orderId {
			return i
		}
	}
	return -1
}

type OrderValidationFields struct {
	Exists   bool
	PosSize  *big.Int
	AsksHead *big.Int
	BidsHead *big.Int
}

func (db *InMemoryDatabase) GetOrderValidationFields(
	orderId common.Hash,
	trader common.Address,
	marketId int,
) OrderValidationFields {
	db.mu.RLock()
	defer db.mu.RUnlock()

	posSize := big.NewInt(0)
	if db.TraderMap[trader] != nil && db.TraderMap[trader].Positions[marketId] != nil && db.TraderMap[trader].Positions[marketId].Size != nil {
		posSize = db.TraderMap[trader].Positions[marketId].Size
	}
	asksHead := big.NewInt(0)
	if len(db.ShortOrders[marketId]) > 0 {
		asksHead = db.ShortOrders[marketId][0].Price
	}
	bidsHead := big.NewInt(0)
	if len(db.LongOrders[marketId]) > 0 {
		bidsHead = db.LongOrders[marketId][0].Price
	}
	fields := OrderValidationFields{
		PosSize:  posSize,
		AsksHead: asksHead,
		BidsHead: bidsHead,
	}
	if db.Orders[orderId] != nil {
		fields.Exists = true
	}
	return fields
}
