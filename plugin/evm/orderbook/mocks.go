package orderbook

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/core/types"
	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockLimitOrderDatabase struct {
	mock.Mock
}

func NewMockLimitOrderDatabase() *MockLimitOrderDatabase {
	return &MockLimitOrderDatabase{}
}

func (db *MockLimitOrderDatabase) SetOrderStatus(orderId common.Hash, status Status, info string, blockNumber uint64) error {
	return nil
}

func (db *MockLimitOrderDatabase) RevertLastStatus(orderId common.Hash) error {
	return nil
}

func (db *MockLimitOrderDatabase) Accept(blockNumber uint64, blockTimestamp uint64) {
}

func (db *MockLimitOrderDatabase) GetAllOrders() []Order {
	args := db.Called()
	return args.Get(0).([]Order)
}

func (db *MockLimitOrderDatabase) GetMarketOrders(market Market) []Order {
	args := db.Called()
	return args.Get(0).([]Order)
}

func (db *MockLimitOrderDatabase) Add(order *Order) {
}

func (db *MockLimitOrderDatabase) AddSignedOrder(order *Order, requiredMargin *big.Int) {
}

func (db *MockLimitOrderDatabase) UpdateFilledBaseAssetQuantity(quantity *big.Int, orderId common.Hash, blockNumber uint64) {
}

func (db *MockLimitOrderDatabase) Delete(id common.Hash) {
}

func (db *MockLimitOrderDatabase) GetLongOrders(market Market, lowerbound *big.Int, blockNumber *big.Int) []Order {
	args := db.Called()
	return args.Get(0).([]Order)
}

func (db *MockLimitOrderDatabase) GetShortOrders(market Market, upperbound *big.Int, blockNumber *big.Int) []Order {
	args := db.Called()
	return args.Get(0).([]Order)
}

func (db *MockLimitOrderDatabase) UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool, blockNumber uint64) {
}

func (db *MockLimitOrderDatabase) UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int) {
}

func (db *MockLimitOrderDatabase) UpdateReservedMargin(trader common.Address, addAmount *big.Int) {
}

func (db *MockLimitOrderDatabase) UpdateUnrealisedFunding(market Market, fundingRate *big.Int) {
}

func (db *MockLimitOrderDatabase) ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int) {
}

func (db *MockLimitOrderDatabase) UpdateNextFundingTime(uint64) {
}

func (db *MockLimitOrderDatabase) GetNextFundingTime() uint64 {
	return 0
}

func (db *MockLimitOrderDatabase) GetAllTraders() map[common.Address]Trader {
	args := db.Called()
	return args.Get(0).(map[common.Address]Trader)
}

func (db *MockLimitOrderDatabase) UpdateLastPrice(market Market, lastPrice *big.Int) {}

func (db *MockLimitOrderDatabase) GetInProgressBlocks() []*types.Block {
	return []*types.Block{}
}

func (db *MockLimitOrderDatabase) UpdateInProgressState(block *types.Block, quantityMap map[string]*big.Int) {
}

func (db *MockLimitOrderDatabase) RemoveInProgressState(block *types.Block, quantityMap map[string]*big.Int) {
}

func (db *MockLimitOrderDatabase) GetLastPrice(market Market) *big.Int {
	args := db.Called()
	return args.Get(0).(*big.Int)
}

func (db *MockLimitOrderDatabase) UpdateNextSamplePITime(nextSamplePITime uint64) {}

func (db *MockLimitOrderDatabase) GetNextSamplePITime() uint64 {
	return 0
}

func (db *MockLimitOrderDatabase) GetOrdersToCancel(oraclePrice map[Market]*big.Int) map[common.Address][]common.Hash {
	args := db.Called()
	return args.Get(0).(map[common.Address][]common.Hash)
}

func (db *MockLimitOrderDatabase) GetLastPrices() map[Market]*big.Int {
	return map[Market]*big.Int{}
}

func (db *MockLimitOrderDatabase) GetNaughtyTraders(hState *hu.HubbleState) ([]LiquidablePosition, map[common.Address][]Order, map[common.Address]*big.Int) {
	return []LiquidablePosition{}, map[common.Address][]Order{}, map[common.Address]*big.Int{}
}

func (db *MockLimitOrderDatabase) GetOrderBookData() InMemoryDatabase {
	return InMemoryDatabase{}
}

func (db *MockLimitOrderDatabase) GetOrderBookDataCopy() (*InMemoryDatabase, error) {
	return &InMemoryDatabase{}, nil
}

func (db *MockLimitOrderDatabase) LoadFromSnapshot(snapshot Snapshot) error {
	return nil
}

func (db *MockLimitOrderDatabase) GetAllOpenOrdersForTrader(trader common.Address) []Order {
	return nil
}

func (db *MockLimitOrderDatabase) GetOpenOrdersForTraderByType(trader common.Address, orderType OrderType) []Order {
	return nil
}

func (db *MockLimitOrderDatabase) UpdateLastPremiumFraction(market Market, trader common.Address, lastPremiumFraction *big.Int, cumulativePremiumFraction *big.Int) {
}

func (db *MockLimitOrderDatabase) GetOrderById(id common.Hash) *Order {
	return nil
}

func (db *MockLimitOrderDatabase) GetTraderInfo(trader common.Address) *Trader {
	return &Trader{}
}

func (db *MockLimitOrderDatabase) GetSamplePIAttemptedTime() uint64 {
	return 0
}

func (db *MockLimitOrderDatabase) SignalSamplePIAttempted(time uint64) {}

func (db *MockLimitOrderDatabase) GetOrderValidationFields(orderId common.Hash, order *hu.SignedOrder) OrderValidationFields {
	return OrderValidationFields{}
}

func (db *MockLimitOrderDatabase) SampleImpactPrice() (impactBids, impactAsks, midPrices []*big.Int) {
	return []*big.Int{}, []*big.Int{}, []*big.Int{}
}

func (db *MockLimitOrderDatabase) RemoveExpiredSignedOrders() {}

type MockLimitOrderTxProcessor struct {
	mock.Mock
}

func NewMockLimitOrderTxProcessor() *MockLimitOrderTxProcessor {
	return &MockLimitOrderTxProcessor{}
}

func (lotp *MockLimitOrderTxProcessor) ExecuteMatchedOrdersTx(incomingOrder Order, matchedOrder Order, fillAmount *big.Int) error {
	args := lotp.Called(incomingOrder, matchedOrder, fillAmount)
	return args.Error(0)
}

func (lotp *MockLimitOrderTxProcessor) PurgeOrderBookTxs() {
	lotp.Called()
}

func (lotp *MockLimitOrderTxProcessor) SetOrderBookTxsBlockNumber(blockNumber uint64) {
	lotp.Called()
}

func (lotp *MockLimitOrderTxProcessor) GetOrderBookTxsCount() uint64 {
	args := lotp.Called()
	return uint64(args.Int(0))
}

func (lotp *MockLimitOrderTxProcessor) ExecuteFundingPaymentTx() error {
	return nil
}

func (lotp *MockLimitOrderTxProcessor) ExecuteSamplePITx() error {
	return nil
}

func (lotp *MockLimitOrderTxProcessor) ExecuteLiquidation(trader common.Address, matchedOrder Order, fillAmount *big.Int) error {
	args := lotp.Called(trader, matchedOrder, fillAmount)
	return args.Error(0)
}

func (lotp *MockLimitOrderTxProcessor) ExecuteLimitOrderCancel(orderIds []LimitOrder) error {
	args := lotp.Called(orderIds)
	return args.Error(0)
}

func (lotp *MockLimitOrderTxProcessor) HandleOrderBookEvent(event *types.Log) {
}

func (lotp *MockLimitOrderTxProcessor) HandleMarginAccountEvent(event *types.Log) {
}

func (lotp *MockLimitOrderTxProcessor) HandleClearingHouseEvent(event *types.Log) {
}

func (lotp *MockLimitOrderTxProcessor) GetUnderlyingPrice() (map[Market]*big.Int, error) {
	return nil, nil
}

func (lotp *MockLimitOrderTxProcessor) UpdateMetrics(block *types.Block) {
	lotp.Called()
}

type MockConfigService struct {
	mock.Mock
}

func (mcs *MockConfigService) GetAcceptableBounds(market Market) (*big.Int, *big.Int) {
	args := mcs.Called()
	return args.Get(0).(*big.Int), args.Get(1).(*big.Int)
}

func (mcs *MockConfigService) GetAcceptableBoundsForLiquidation(market Market) (*big.Int, *big.Int) {
	args := mcs.Called(market)
	return args.Get(0).(*big.Int), args.Get(1).(*big.Int)
}

func (mcs *MockConfigService) getLiquidationSpreadThreshold(market Market) *big.Int {
	return big.NewInt(1e4)
}

func (mcs *MockConfigService) getMaxLiquidationRatio(market Market) *big.Int {
	args := mcs.Called()
	return args.Get(0).(*big.Int)
}

func (mcs *MockConfigService) GetMinAllowableMargin() *big.Int {
	args := mcs.Called()
	return args.Get(0).(*big.Int)
}

func (mcs *MockConfigService) GetMaintenanceMargin() *big.Int {
	args := mcs.Called()
	return args.Get(0).(*big.Int)
}

func (mcs *MockConfigService) getMinSizeRequirement(market Market) *big.Int {
	return big.NewInt(1)
}

func (cs *MockConfigService) GetActiveMarketsCount() int64 {
	return int64(1)
}

func (cs *MockConfigService) GetUnderlyingPrices() []*big.Int {
	return []*big.Int{}
}

func (cs *MockConfigService) GetMidPrices() []*big.Int {
	return []*big.Int{}
}

func (cs *MockConfigService) GetLastPremiumFraction(market Market, trader *common.Address) *big.Int {
	return big.NewInt(0)
}

func (cs *MockConfigService) GetLastPremiumFractionAtBlock(market Market, trader *common.Address, blockNumber uint64) *big.Int {
	return big.NewInt(0)
}

func (cs *MockConfigService) GetCumulativePremiumFraction(market Market) *big.Int {
	return big.NewInt(0)
}

func (cs *MockConfigService) GetCumulativePremiumFractionAtBlock(market Market, blockNumber uint64) *big.Int {
	return big.NewInt(0)
}

func (cs *MockConfigService) GetCollaterals() []hu.Collateral {
	return []hu.Collateral{{Price: big.NewInt(1e6), Weight: big.NewInt(1e6), Decimals: 6}}
}

func (cs *MockConfigService) GetTakerFee() *big.Int {
	return big.NewInt(0)
}

func (cs *MockConfigService) HasReferrer(trader common.Address) bool {
	return true
}

func (cs *MockConfigService) GetPriceMultiplier(market Market) *big.Int {
	return big.NewInt(1e6)
}

func (cs *MockConfigService) GetSignedOrderStatus(orderHash common.Hash) int64 {
	return 0
}

func (cs *MockConfigService) IsTradingAuthority(trader, signer common.Address) bool {
	return false
}

func NewMockConfigService() *MockConfigService {
	return &MockConfigService{}
}

func (cs *MockConfigService) GetSignedOrderbookContract() common.Address {
	return common.Address{}
}

func (cs *MockConfigService) GetUpgradeVersion() hu.UpgradeVersion {
	return hu.V2
}

func (cs *MockConfigService) GetMarketAddressFromMarketID(marketId int64) common.Address {
	return common.Address{}
}

func (cs *MockConfigService) GetImpactMarginNotional(ammAddress common.Address) *big.Int {
	return big.NewInt(500e6)
}

func (cs *MockConfigService) GetReduceOnlyAmounts(trader common.Address) []*big.Int {
	return []*big.Int{big.NewInt(0)}
}
