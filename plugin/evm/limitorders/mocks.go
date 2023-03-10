package limitorders

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockLimitOrderDatabase struct {
	mock.Mock
}

func NewMockLimitOrderDatabase() *MockLimitOrderDatabase {
	return &MockLimitOrderDatabase{}
}

func (db *MockLimitOrderDatabase) GetAllOrders() []LimitOrder {
	args := db.Called()
	return args.Get(0).([]LimitOrder)
}

func (db *MockLimitOrderDatabase) Add(order *LimitOrder) {
}

func (db *MockLimitOrderDatabase) UpdateFilledBaseAssetQuantity(quantity *big.Int, orderId string) {
}

func (db *MockLimitOrderDatabase) Delete(id string) {
}

func (db *MockLimitOrderDatabase) GetLongOrders(market Market) []LimitOrder {
	args := db.Called()
	return args.Get(0).([]LimitOrder)
}

func (db *MockLimitOrderDatabase) GetShortOrders(market Market) []LimitOrder {
	args := db.Called()
	return args.Get(0).([]LimitOrder)
}

func (db *MockLimitOrderDatabase) UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int, isLiquidation bool) {
}

func (db *MockLimitOrderDatabase) UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int) {
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

func (db *MockLimitOrderDatabase) UpdateLastPrice(market Market, lastPrice *big.Int) {
}

func (db *MockLimitOrderDatabase) GetLastPrice(market Market) *big.Int {
	args := db.Called()
	return args.Get(0).(*big.Int)
}

func (db *MockLimitOrderDatabase) GetOrderBookData() InMemoryDatabase {
	return *&InMemoryDatabase{}
}

type MockLimitOrderTxProcessor struct {
	mock.Mock
}

func NewMockLimitOrderTxProcessor() *MockLimitOrderTxProcessor {
	return &MockLimitOrderTxProcessor{}
}

func (lotp *MockLimitOrderTxProcessor) ExecuteMatchedOrdersTx(incomingOrder LimitOrder, matchedOrder LimitOrder, fillAmount *big.Int) error {
	args := lotp.Called(incomingOrder, matchedOrder, fillAmount)
	return args.Error(0)
}

func (lotp *MockLimitOrderTxProcessor) PurgeLocalTx() {
	lotp.Called()
}

func (lotp *MockLimitOrderTxProcessor) CheckIfOrderBookContractCall(tx *types.Transaction) bool {
	return true
}

func (lotp *MockLimitOrderTxProcessor) ExecuteFundingPaymentTx() error {
	return nil
}

func (lotp *MockLimitOrderTxProcessor) ExecuteLiquidation(trader common.Address, matchedOrder LimitOrder, fillAmount *big.Int) error {
	args := lotp.Called(trader, matchedOrder, fillAmount)
	return args.Error(0)
}

func (lotp *MockLimitOrderTxProcessor) HandleOrderBookEvent(event *types.Log) {
}

func (lotp *MockLimitOrderTxProcessor) HandleMarginAccountEvent(event *types.Log) {
}

func (lotp *MockLimitOrderTxProcessor) HandleClearingHouseEvent(event *types.Log) {
}
