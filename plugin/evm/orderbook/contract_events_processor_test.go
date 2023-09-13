package orderbook

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook/abis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var timestamp = big.NewInt(time.Now().Unix())

func TestProcessEvents(t *testing.T) {
	// this test is obsolete because we expect the events to automatically come in sorted order
	t.Run("it sorts events by blockNumber and executes in order", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		orderBookABI := getABIfromJson(abis.OrderBookAbi)
		limitOrderBookABI := getABIfromJson(abis.OrderBookAbi)

		traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
		ammIndex := big.NewInt(0)
		baseAssetQuantity := big.NewInt(5000000000000000000)
		price := big.NewInt(1000000000)
		salt1 := big.NewInt(1675239557437)
		longOrder := getLimitOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt1)
		longOrderId := getIdFromLimitOrder(longOrder)

		salt2 := big.NewInt(0).Add(salt1, big.NewInt(1))
		shortOrder := getLimitOrder(ammIndex, traderAddress, big.NewInt(0).Neg(baseAssetQuantity), price, salt2)
		shortOrderId := getIdFromLimitOrder(shortOrder)

		ordersPlacedBlockNumber := uint64(12)
		orderAcceptedEvent := getEventFromABI(limitOrderBookABI, "OrderAccepted")
		longOrderAcceptedEventTopics := []common.Hash{orderAcceptedEvent.ID, traderAddress.Hash(), longOrderId}
		longOrderAcceptedEventData, err := orderAcceptedEvent.Inputs.NonIndexed().Pack(longOrder, timestamp)
		if err != nil {
			t.Fatalf("%s", err)
		}
		longOrderAcceptedEventLog := getEventLog(OrderBookContractAddress, longOrderAcceptedEventTopics, longOrderAcceptedEventData, ordersPlacedBlockNumber)

		shortOrderAcceptedEventTopics := []common.Hash{orderAcceptedEvent.ID, traderAddress.Hash(), shortOrderId}
		shortOrderAcceptedEventData, err := orderAcceptedEvent.Inputs.NonIndexed().Pack(shortOrder, timestamp)
		if err != nil {
			t.Fatalf("%s", err)
		}
		shortOrderAcceptedEventLog := getEventLog(OrderBookContractAddress, shortOrderAcceptedEventTopics, shortOrderAcceptedEventData, ordersPlacedBlockNumber)

		orderMatchedBlockNumber := uint64(14)
		orderMatchedEvent0 := getEventFromABI(orderBookABI, "OrderMatched")
		orderMatchedEvent1 := getEventFromABI(orderBookABI, "OrderMatched")
		orderMatchedEventTopics0 := []common.Hash{orderMatchedEvent0.ID, traderAddress.Hash(), longOrderId}
		orderMatchedEventTopics1 := []common.Hash{orderMatchedEvent1.ID, traderAddress.Hash(), shortOrderId}
		fillAmount := big.NewInt(3000000000000000000)
		fmt.Printf("sending matched event %s and %s", longOrderId.String(), shortOrderId.String())
		orderMatchedEventData0, _ := orderMatchedEvent0.Inputs.NonIndexed().Pack(fillAmount, price, big.NewInt(0), timestamp)
		orderMatchedEventData1, _ := orderMatchedEvent0.Inputs.NonIndexed().Pack(fillAmount, price, big.NewInt(0), timestamp)
		orderMatchedEventLog0 := getEventLog(OrderBookContractAddress, orderMatchedEventTopics0, orderMatchedEventData0, orderMatchedBlockNumber)
		orderMatchedEventLog1 := getEventLog(OrderBookContractAddress, orderMatchedEventTopics1, orderMatchedEventData1, orderMatchedBlockNumber)
		cep.ProcessEvents([]*types.Log{longOrderAcceptedEventLog, shortOrderAcceptedEventLog, orderMatchedEventLog0, orderMatchedEventLog1})

		actualLongOrder := db.OrderMap[getIdFromLimitOrder(longOrder)]
		assert.Equal(t, fillAmount, actualLongOrder.FilledBaseAssetQuantity)

		actualShortOrder := db.OrderMap[getIdFromLimitOrder(shortOrder)]
		assert.Equal(t, big.NewInt(0).Neg(fillAmount), actualShortOrder.FilledBaseAssetQuantity)
	})

	// t.Run("when event is removed it is not processed", func(t *testing.T) {
	// 	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	// 	db := getDatabase()
	// 	collateral := HUSD
	// 	originalMargin := multiplyBasePrecision(big.NewInt(100))
	// 	trader := &Trader{
	// 		Margins: map[Collateral]*big.Int{collateral: big.NewInt(0).Set(originalMargin)},
	// 	}
	// 	db.TraderMap[traderAddress] = trader
	// 	blockNumber := uint64(12)

	// 	//MarginAccount Contract log
	// 	marginAccountABI := getABIfromJson(marginAccountAbi)
	// 	marginAccountEvent := getEventFromABI(marginAccountABI, "MarginAdded")
	// 	marginAccountEventTopics := []common.Hash{marginAccountEvent.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
	// 	marginAdded := multiplyBasePrecision(big.NewInt(100))
	// 	timestamp := big.NewInt(time.Now().Unix())
	// 	marginAddedEventData, _ := marginAccountEvent.Inputs.NonIndexed().Pack(marginAdded, timestamp)
	// 	marginAddedLog := getEventLog(MarginAccountContractAddress, marginAccountEventTopics, marginAddedEventData, blockNumber)
	// 	marginAddedLog.Removed = true
	// 	cep := newcep(t, db)

	// 	cep.ProcessEvents([]*types.Log{marginAddedLog})
	// 	assert.Equal(t, originalMargin, db.TraderMap[traderAddress].Margins[collateral])
	// })
}

func TestOrderBookMarginAccountClearingHouseEventInLog(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	blockNumber := uint64(12)
	db := getDatabase()
	cep := newcep(t, db)
	collateral := HUSD
	openNotional := multiplyBasePrecision(big.NewInt(100))
	size := multiplyPrecisionSize(big.NewInt(10))
	lastPremiumFraction := multiplyBasePrecision(big.NewInt(1))
	liquidationThreshold := multiplyBasePrecision(big.NewInt(1))
	unrealisedFunding := multiplyBasePrecision(big.NewInt(1))
	market := Market(0)
	position := &Position{
		OpenNotional:         openNotional,
		Size:                 size,
		UnrealisedFunding:    unrealisedFunding,
		LastPremiumFraction:  lastPremiumFraction,
		LiquidationThreshold: liquidationThreshold,
	}
	originalMargin := multiplyBasePrecision(big.NewInt(100))
	trader := &Trader{
		Margin:    Margin{Deposited: map[Collateral]*big.Int{collateral: big.NewInt(0).Set(originalMargin)}},
		Positions: map[Market]*Position{market: position},
	}
	db.TraderMap[traderAddress] = trader

	//OrderBook Contract log
	ammIndex := big.NewInt(0)
	baseAssetQuantity := big.NewInt(5000000000000000000)
	price := big.NewInt(1000000000)
	salt := big.NewInt(1675239557437)
	order := getLimitOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt)
	orderBookABI := getABIfromJson(abis.OrderBookAbi)
	limitOrderBookABI := getABIfromJson(abis.OrderBookAbi)
	orderBookEvent := getEventFromABI(limitOrderBookABI, "OrderAccepted")
	orderAcceptedEventData, _ := orderBookEvent.Inputs.NonIndexed().Pack(order, timestamp)
	orderBookEventTopics := []common.Hash{orderBookEvent.ID, traderAddress.Hash(), getIdFromLimitOrder(order)}
	orderBookLog := getEventLog(OrderBookContractAddress, orderBookEventTopics, orderAcceptedEventData, blockNumber)

	//MarginAccount Contract log
	marginAccountABI := getABIfromJson(abis.MarginAccountAbi)
	marginAccountEvent := getEventFromABI(marginAccountABI, "MarginAdded")
	marginAccountEventTopics := []common.Hash{marginAccountEvent.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
	marginAdded := multiplyBasePrecision(big.NewInt(100))
	marginAddedEventData, _ := marginAccountEvent.Inputs.NonIndexed().Pack(marginAdded, timestamp)
	marginAccountLog := getEventLog(MarginAccountContractAddress, marginAccountEventTopics, marginAddedEventData, blockNumber)

	//ClearingHouse Contract log
	clearingHouseABI := getABIfromJson(abis.ClearingHouseAbi)
	clearingHouseEvent := getEventFromABI(clearingHouseABI, "FundingRateUpdated")
	clearingHouseEventTopics := []common.Hash{clearingHouseEvent.ID, common.BigToHash(big.NewInt(int64(market)))}

	nextFundingTime := big.NewInt(time.Now().Unix())
	premiumFraction := multiplyBasePrecision(big.NewInt(10))
	underlyingPrice := multiplyBasePrecision(big.NewInt(100))
	cumulativePremiumFraction := multiplyBasePrecision(big.NewInt(10))
	fundingRateUpdated, _ := clearingHouseEvent.Inputs.NonIndexed().Pack(premiumFraction, underlyingPrice, cumulativePremiumFraction, nextFundingTime, timestamp, big.NewInt(int64(blockNumber)))
	clearingHouseLog := getEventLog(ClearingHouseContractAddress, clearingHouseEventTopics, fundingRateUpdated, blockNumber)

	// logs := []*types.Log{orderBookLog, marginAccountLog, clearingHouseLog}
	cep.ProcessEvents([]*types.Log{orderBookLog})
	cep.ProcessAcceptedEvents([]*types.Log{marginAccountLog, clearingHouseLog}, true)

	//OrderBook log - OrderAccepted
	actualLimitOrder := *db.GetOrderBookData().OrderMap[getIdFromLimitOrder(order)]
	args := map[string]interface{}{}
	orderBookABI.UnpackIntoMap(args, "OrderAccepted", orderAcceptedEventData)
	assert.Equal(t, Market(ammIndex.Int64()), actualLimitOrder.Market)
	assert.Equal(t, LONG, actualLimitOrder.PositionType)
	assert.Equal(t, traderAddress.String(), actualLimitOrder.Trader.String())
	assert.Equal(t, *baseAssetQuantity, *actualLimitOrder.BaseAssetQuantity)
	assert.Equal(t, *price, *actualLimitOrder.Price)
	assert.Equal(t, Placed, actualLimitOrder.getOrderStatus().Status)
	assert.Equal(t, big.NewInt(int64(blockNumber)), actualLimitOrder.BlockNumber)
	rawOrder := &LimitOrder{}
	rawOrder.DecodeFromRawOrder(args["order"])
	assert.Equal(t, rawOrder, actualLimitOrder.RawOrder.(*LimitOrder))

	//ClearingHouse log - FundingRateUpdated
	expectedUnrealisedFunding := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size))
	assert.Equal(t, expectedUnrealisedFunding, db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)

	//MarginAccount log - marginAdded
	actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margin.Deposited[collateral]
	assert.Equal(t, big.NewInt(0).Add(marginAdded, originalMargin), actualMargin)

}

func TestHandleOrderBookEvent(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	ammIndex := big.NewInt(0)
	baseAssetQuantity := big.NewInt(5000000000000000000)
	price := big.NewInt(1000000000)
	salt := big.NewInt(1675239557437)
	blockNumber := uint64(12)
	orderBookABI := getABIfromJson(abis.OrderBookAbi)

	t.Run("When event is OrderAccepted", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrderAccepted")
		order := getLimitOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt)
		orderId := getIdFromLimitOrder(order)
		topics := []common.Hash{event.ID, traderAddress.Hash(), orderId}
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderAcceptedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, orderAcceptedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Nil(t, actualLimitOrder)
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			orderAcceptedEventData, err := event.Inputs.NonIndexed().Pack(order, timestamp)
			if err != nil {
				t.Fatalf("%s", err)
			}
			log := getEventLog(OrderBookContractAddress, topics, orderAcceptedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})

			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			args := map[string]interface{}{}
			orderBookABI.UnpackIntoMap(args, "OrderAccepted", orderAcceptedEventData)
			assert.Equal(t, Market(ammIndex.Int64()), actualLimitOrder.Market)
			assert.Equal(t, LONG, actualLimitOrder.PositionType)
			assert.Equal(t, traderAddress.String(), actualLimitOrder.Trader.String())
			assert.Equal(t, *baseAssetQuantity, *actualLimitOrder.BaseAssetQuantity)
			assert.Equal(t, false, actualLimitOrder.ReduceOnly)
			assert.Equal(t, false, actualLimitOrder.isPostOnly())
			assert.Equal(t, *price, *actualLimitOrder.Price)
			assert.Equal(t, Placed, actualLimitOrder.getOrderStatus().Status)
			assert.Equal(t, big.NewInt(int64(blockNumber)), actualLimitOrder.BlockNumber)
			rawOrder := &LimitOrder{}
			rawOrder.DecodeFromRawOrder(args["order"])
			assert.Equal(t, rawOrder, actualLimitOrder.RawOrder.(*LimitOrder))
		})
	})
	t.Run("When event is OrderCancelAccepted", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		order := getLimitOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt)
		event := getEventFromABI(orderBookABI, "OrderCancelAccepted")
		topics := []common.Hash{event.ID, traderAddress.Hash(), getIdFromLimitOrder(order)}
		blockNumber := uint64(4)
		limitOrder := &Order{
			Market:            Market(ammIndex.Int64()),
			PositionType:      LONG,
			Trader:            traderAddress,
			BaseAssetQuantity: baseAssetQuantity,
			Price:             price,
			BlockNumber:       big.NewInt(1),
			Salt:              salt,
			RawOrder:          &order,
		}
		limitOrder.Id = getIdFromOrder(*limitOrder)
		db.Add(limitOrder)
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderCancelAcceptedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, orderCancelAcceptedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			orderId := getIdFromOrder(*limitOrder)
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Equal(t, limitOrder, actualLimitOrder)
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			orderCancelAcceptedEventData, _ := event.Inputs.NonIndexed().Pack(timestamp)
			log := getEventLog(OrderBookContractAddress, topics, orderCancelAcceptedEventData, blockNumber)
			orderId := getIdFromOrder(*limitOrder)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Equal(t, Cancelled, actualLimitOrder.getOrderStatus().Status)
		})
	})
	t.Run("When event is OrderMatched", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrderMatched")
		longOrder := &Order{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            LONG,
			Trader:                  traderAddress,
			BaseAssetQuantity:       baseAssetQuantity,
			Price:                   price,
			BlockNumber:             big.NewInt(1),
			FilledBaseAssetQuantity: big.NewInt(0),
			Salt:                    salt,
		}
		shortOrder := &Order{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            SHORT,
			Trader:                  traderAddress,
			BaseAssetQuantity:       big.NewInt(0).Mul(baseAssetQuantity, big.NewInt(-1)),
			Price:                   price,
			BlockNumber:             big.NewInt(1),
			FilledBaseAssetQuantity: big.NewInt(0),
			Salt:                    big.NewInt(0).Add(salt, big.NewInt(1000)),
		}

		longOrder.Id = getIdFromOrder(*longOrder)
		shortOrder.Id = getIdFromOrder(*shortOrder)
		db.Add(longOrder)
		db.Add(shortOrder)
		// relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		fillAmount := big.NewInt(10)
		topics0 := []common.Hash{event.ID, traderAddress.Hash(), longOrder.Id}
		topics1 := []common.Hash{event.ID, traderAddress.Hash(), shortOrder.Id}
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderMatchedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics0, orderMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, int64(0), longOrder.FilledBaseAssetQuantity.Int64())
			assert.Equal(t, int64(0), shortOrder.FilledBaseAssetQuantity.Int64())
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			orderMatchedEventData, _ := event.Inputs.NonIndexed().Pack(fillAmount, price, big.NewInt(0).Mul(fillAmount, price), timestamp)
			log0 := getEventLog(OrderBookContractAddress, topics0, orderMatchedEventData, blockNumber)
			log1 := getEventLog(OrderBookContractAddress, topics1, orderMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log0, log1})
			assert.Equal(t, big.NewInt(fillAmount.Int64()), longOrder.FilledBaseAssetQuantity)
			assert.Equal(t, big.NewInt(-fillAmount.Int64()), shortOrder.FilledBaseAssetQuantity)
		})
	})
}

func TestHandleMarginAccountEvent(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	blockNumber := uint64(12)
	collateral := HUSD

	marginAccountABI := getABIfromJson(abis.MarginAccountAbi)

	t.Run("when event is MarginAdded", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(marginAccountABI, "MarginAdded")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
		t.Run("When event parsing fails", func(t *testing.T) {
			marginAddedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, marginAddedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			marginAdded := big.NewInt(10000)
			timestamp := big.NewInt(time.Now().Unix())
			marginAddedEventData, _ := event.Inputs.NonIndexed().Pack(marginAdded, timestamp)
			log := getEventLog(MarginAccountContractAddress, topics, marginAddedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margin.Deposited[collateral]
			assert.Equal(t, marginAdded, actualMargin)
		})
	})
	t.Run("when event is MarginRemoved", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(marginAccountABI, "MarginRemoved")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
		t.Run("When event parsing fails", func(t *testing.T) {
			marginRemovedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, marginRemovedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			marginRemoved := big.NewInt(10000)
			marginRemovedEventData, _ := event.Inputs.NonIndexed().Pack(marginRemoved, timestamp)
			log := getEventLog(MarginAccountContractAddress, topics, marginRemovedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margin.Deposited[collateral]
			assert.Equal(t, big.NewInt(0).Neg(marginRemoved), actualMargin)
		})
	})
	t.Run("when event is PnLRealized", func(t *testing.T) {
		event := getEventFromABI(marginAccountABI, "PnLRealized")
		topics := []common.Hash{event.ID, traderAddress.Hash()}
		db := getDatabase()
		cep := newcep(t, db)
		t.Run("When event parsing fails", func(t *testing.T) {
			pnlRealizedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, pnlRealizedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			pnlRealized := big.NewInt(-10000)
			pnlRealizedEventData, _ := event.Inputs.NonIndexed().Pack(pnlRealized, timestamp)
			log := getEventLog(MarginAccountContractAddress, topics, pnlRealizedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margin.Deposited[collateral]
			assert.Equal(t, pnlRealized, actualMargin)
		})
	})

	t.Run("when event is MarginReserved", func(t *testing.T) {
		event := getEventFromABI(marginAccountABI, "MarginReserved")
		topics := []common.Hash{event.ID, traderAddress.Hash()}
		db := getDatabase()
		cep := newcep(t, db)
		t.Run("When event parsing fails", func(t *testing.T) {
			marginReservedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, marginReservedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			reservedMargin := big.NewInt(10000000)
			marginReservedEventData, _ := event.Inputs.NonIndexed().Pack(reservedMargin)
			log := getEventLog(MarginAccountContractAddress, topics, marginReservedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			reservedMarginInDb := db.GetOrderBookData().TraderMap[traderAddress].Margin.Reserved
			assert.Equal(t, reservedMargin, reservedMarginInDb)
		})
	})

	t.Run("when event is MarginReleased", func(t *testing.T) {
		event := getEventFromABI(marginAccountABI, "MarginReleased")
		topics := []common.Hash{event.ID, traderAddress.Hash()}
		db := getDatabase()
		cep := newcep(t, db)
		t.Run("When event parsing fails", func(t *testing.T) {
			marginReleasedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, marginReleasedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			releasedMargin := big.NewInt(10000000)
			marginReleasedEventData, _ := event.Inputs.NonIndexed().Pack(releasedMargin)
			log := getEventLog(MarginAccountContractAddress, topics, marginReleasedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			releasedMarginInDb := db.GetOrderBookData().TraderMap[traderAddress].Margin.Reserved
			assert.Equal(t, big.NewInt(0).Neg(releasedMargin), releasedMarginInDb)
		})
	})
}
func TestHandleClearingHouseEvent(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	blockNumber := uint64(12)
	collateral := HUSD
	market := Market(0)
	clearingHouseABI := getABIfromJson(abis.ClearingHouseAbi)
	openNotional := multiplyBasePrecision(big.NewInt(100))
	size := multiplyPrecisionSize(big.NewInt(10))
	lastPremiumFraction := multiplyBasePrecision(big.NewInt(1))
	liquidationThreshold := multiplyBasePrecision(big.NewInt(1))
	unrealisedFunding := multiplyBasePrecision(big.NewInt(1))
	t.Run("when event is FundingRateUpdated", func(t *testing.T) {
		event := getEventFromABI(clearingHouseABI, "FundingRateUpdated")
		topics := []common.Hash{event.ID, common.BigToHash(big.NewInt(int64(market)))}
		db := getDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margin:    Margin{Deposited: map[Collateral]*big.Int{collateral: big.NewInt(100)}},
			Positions: map[Market]*Position{market: position},
		}
		db.TraderMap[traderAddress] = trader

		t.Run("When event parsing fails", func(t *testing.T) {
			pnlRealizedEventData := []byte{}
			log := getEventLog(ClearingHouseContractAddress, topics, pnlRealizedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})

			assert.Equal(t, uint64(0), db.NextFundingTime)
			assert.Equal(t, unrealisedFunding, db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			nextFundingTime := big.NewInt(time.Now().Unix())
			premiumFraction := multiplyBasePrecision(big.NewInt(10))
			underlyingPrice := multiplyBasePrecision(big.NewInt(100))
			cumulativePremiumFraction := multiplyBasePrecision(big.NewInt(10))
			fundingRateUpdated, _ := event.Inputs.NonIndexed().Pack(premiumFraction, underlyingPrice, cumulativePremiumFraction, nextFundingTime, timestamp, big.NewInt(int64(blockNumber)))
			log := getEventLog(ClearingHouseContractAddress, topics, fundingRateUpdated, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			expectedUnrealisedFunding := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size))
			assert.Equal(t, expectedUnrealisedFunding, db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)
		})
	})
	t.Run("When event is FundingPaid", func(t *testing.T) {
		event := getEventFromABI(clearingHouseABI, "FundingPaid")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(market)))}
		db := getDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margin:    Margin{Deposited: map[Collateral]*big.Int{collateral: big.NewInt(100)}},
			Positions: map[Market]*Position{market: position},
		}
		db.TraderMap[traderAddress] = trader

		t.Run("When event parsing fails", func(t *testing.T) {
			pnlRealizedEventData := []byte{}
			log := getEventLog(ClearingHouseContractAddress, topics, pnlRealizedEventData, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)

			assert.Equal(t, unrealisedFunding, db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)
			assert.Equal(t, lastPremiumFraction, db.TraderMap[traderAddress].Positions[market].LastPremiumFraction)
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			takerFundingPayment := multiplyBasePrecision(big.NewInt(10))
			cumulativePremiumFraction := multiplyBasePrecision(big.NewInt(10))
			fundingPaidEvent, _ := event.Inputs.NonIndexed().Pack(takerFundingPayment, cumulativePremiumFraction)
			log := getEventLog(ClearingHouseContractAddress, topics, fundingPaidEvent, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Equal(t, big.NewInt(0), db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)
			assert.Equal(t, cumulativePremiumFraction, db.TraderMap[traderAddress].Positions[market].LastPremiumFraction)
		})
	})
	t.Run("When event is PositionModified", func(t *testing.T) {
		event := getEventFromABI(clearingHouseABI, "PositionModified")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(market)))}
		db := getDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margin:    Margin{Deposited: map[Collateral]*big.Int{collateral: big.NewInt(100)}},
			Positions: map[Market]*Position{market: position},
		}
		db.TraderMap[traderAddress] = trader

		t.Run("When event parsing fails", func(t *testing.T) {
			positionModifiedEvent := []byte{}
			log := getEventLog(ClearingHouseContractAddress, topics, positionModifiedEvent, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Nil(t, db.LastPrice[market])
			// assert.Equal(t, big.NewInt(0), db.LastPrice[market])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			baseAsset := multiplyPrecisionSize(big.NewInt(10))
			// quoteAsset := multiplyBasePrecision(big.NewInt(1000))
			realizedPnl := multiplyBasePrecision(big.NewInt(20))
			openNotional := multiplyBasePrecision(big.NewInt(4000))
			timestamp := multiplyBasePrecision(big.NewInt(time.Now().Unix()))
			size := multiplyPrecisionSize(big.NewInt(40))
			price := multiplyBasePrecision(big.NewInt(100)) // baseAsset / quoteAsset

			positionModifiedEvent, err := event.Inputs.NonIndexed().Pack(baseAsset, price, realizedPnl, size, openNotional, big.NewInt(0), uint8(0), timestamp)
			if err != nil {
				t.Fatal(err)
			}
			log := getEventLog(ClearingHouseContractAddress, topics, positionModifiedEvent, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)

			// quoteAsset/(baseAsset / 1e 18)
			expectedLastPrice := big.NewInt(100000000)
			assert.Equal(t, expectedLastPrice, db.LastPrice[market])
			assert.Equal(t, size, db.TraderMap[traderAddress].Positions[market].Size)
			assert.Equal(t, openNotional, db.TraderMap[traderAddress].Positions[market].OpenNotional)
		})
	})
	t.Run("When event is PositionLiquidated", func(t *testing.T) {
		event := getEventFromABI(clearingHouseABI, "PositionLiquidated")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(market)))}
		db := getDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margin:    Margin{Deposited: map[Collateral]*big.Int{collateral: big.NewInt(100)}},
			Positions: map[Market]*Position{market: position},
		}
		db.TraderMap[traderAddress] = trader

		t.Run("When event parsing fails", func(t *testing.T) {
			positionLiquidatedEvent := []byte{}
			log := getEventLog(ClearingHouseContractAddress, topics, positionLiquidatedEvent, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)
			assert.Nil(t, db.LastPrice[market])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			baseAsset := multiplyPrecisionSize(big.NewInt(10))
			// quoteAsset := multiplyBasePrecision(big.NewInt(1000))
			realizedPnl := multiplyBasePrecision(big.NewInt(20))
			openNotional := multiplyBasePrecision(big.NewInt(4000))
			timestamp := multiplyBasePrecision(big.NewInt(time.Now().Unix()))
			size := multiplyPrecisionSize(big.NewInt(40))
			price := multiplyBasePrecision(big.NewInt(100)) // baseAsset / quoteAsset

			positionLiquidatedEvent, _ := event.Inputs.NonIndexed().Pack(baseAsset, price, realizedPnl, size, openNotional, big.NewInt(0), timestamp)
			log := getEventLog(ClearingHouseContractAddress, topics, positionLiquidatedEvent, blockNumber)
			cep.ProcessAcceptedEvents([]*types.Log{log}, true)

			// quoteAsset/(baseAsset / 1e 18)
			expectedLastPrice := big.NewInt(100000000)
			assert.Equal(t, expectedLastPrice, db.LastPrice[market])
			assert.Equal(t, size, db.TraderMap[traderAddress].Positions[market].Size)
			assert.Equal(t, openNotional, db.TraderMap[traderAddress].Positions[market].OpenNotional)
		})
	})
}

func TestRemovedEvents(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	traderAddress2 := common.HexToAddress("0xC348509BD9dD348b963B4ae0CB876782387729a0")
	blockNumber := big.NewInt(12)
	ammIndex := big.NewInt(0)
	baseAssetQuantity := big.NewInt(50)
	salt1 := big.NewInt(1675239557437)
	salt2 := big.NewInt(1675239557439)
	orderBookABI := getABIfromJson(abis.OrderBookAbi)
	limitOrderrderBookABI := getABIfromJson(abis.OrderBookAbi)

	db := getDatabase()
	cep := newcep(t, db)

	orderAcceptedEvent := getEventFromABI(limitOrderrderBookABI, "OrderAccepted")
	longOrder := getLimitOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt1)
	longOrderId := getIdFromLimitOrder(longOrder)
	longOrderAcceptedEventTopics := []common.Hash{orderAcceptedEvent.ID, traderAddress.Hash(), longOrderId}
	longOrderAcceptedEventData, _ := orderAcceptedEvent.Inputs.NonIndexed().Pack(longOrder, timestamp)

	shortOrder := getLimitOrder(ammIndex, traderAddress2, big.NewInt(0).Neg(baseAssetQuantity), price, salt2)
	shortOrderId := getIdFromLimitOrder(shortOrder)
	shortOrderAcceptedEventTopics := []common.Hash{orderAcceptedEvent.ID, traderAddress2.Hash(), shortOrderId}
	shortOrderAcceptedEventData, _ := orderAcceptedEvent.Inputs.NonIndexed().Pack(shortOrder, timestamp)

	t.Run("delete order when OrderAccepted is removed", func(t *testing.T) {
		longOrderAcceptedEventLog := getEventLog(OrderBookContractAddress, longOrderAcceptedEventTopics, longOrderAcceptedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderAcceptedEventLog})

		// order exists in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)

		// order should be deleted if OrderAccepted log is removed
		longOrderAcceptedEventLog.Removed = true
		cep.ProcessEvents([]*types.Log{longOrderAcceptedEventLog})
		assert.Nil(t, db.OrderMap[longOrderId])
	})

	t.Run("un-cancel an order when OrderCancelAccepted is removed", func(t *testing.T) {
		longOrderAcceptedEventLog := getEventLog(OrderBookContractAddress, longOrderAcceptedEventTopics, longOrderAcceptedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderAcceptedEventLog})

		// order exists in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)

		// cancel it
		orderCancelAcceptedEvent := getEventFromABI(limitOrderrderBookABI, "OrderCancelAccepted")
		orderCancelAcceptedEventTopics := []common.Hash{orderCancelAcceptedEvent.ID, traderAddress.Hash(), longOrderId}
		orderCancelAcceptedEventData, _ := orderCancelAcceptedEvent.Inputs.NonIndexed().Pack(timestamp)
		orderCancelAcceptedLog := getEventLog(OrderBookContractAddress, orderCancelAcceptedEventTopics, orderCancelAcceptedEventData, blockNumber.Uint64()+2)
		cep.ProcessEvents([]*types.Log{orderCancelAcceptedLog})

		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Cancelled)

		// now uncancel it
		orderCancelAcceptedLog.Removed = true
		cep.ProcessEvents([]*types.Log{orderCancelAcceptedLog})
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Placed)
	})

	t.Run("un-fulfill an order when OrderMatched is removed", func(t *testing.T) {
		longOrderAcceptedEventLog := getEventLog(OrderBookContractAddress, longOrderAcceptedEventTopics, longOrderAcceptedEventData, blockNumber.Uint64())
		shortOrderAcceptedEventLog := getEventLog(OrderBookContractAddress, shortOrderAcceptedEventTopics, shortOrderAcceptedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderAcceptedEventLog, shortOrderAcceptedEventLog})

		// orders exist in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)
		assert.Equal(t, db.OrderMap[shortOrderId].Salt, shortOrder.Salt)

		// fulfill them
		orderMatchedEvent := getEventFromABI(orderBookABI, "OrderMatched")
		orderMatchedEventTopics := []common.Hash{orderMatchedEvent.ID, traderAddress.Hash(), longOrderId}
		orderMatchedEventData, _ := orderMatchedEvent.Inputs.NonIndexed().Pack(baseAssetQuantity, price, big.NewInt(0).Mul(baseAssetQuantity, price), timestamp)
		orderMatchedLog := getEventLog(OrderBookContractAddress, orderMatchedEventTopics, orderMatchedEventData, blockNumber.Uint64()+2)
		cep.ProcessEvents([]*types.Log{orderMatchedLog})

		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, FulFilled)

		// now un-fulfill it
		orderMatchedLog.Removed = true
		cep.ProcessEvents([]*types.Log{orderMatchedLog})
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Placed)
	})

	t.Run("revert state of an order when OrderMatchingError is removed", func(t *testing.T) {
		// change salt
		longOrder.Salt.Add(longOrder.Salt, big.NewInt(20))
		longOrderId = getIdFromLimitOrder(longOrder)
		longOrderAcceptedEventTopics = []common.Hash{orderAcceptedEvent.ID, traderAddress.Hash(), longOrderId}
		longOrderAcceptedEventData, _ = orderAcceptedEvent.Inputs.NonIndexed().Pack(longOrder, timestamp)
		longOrderAcceptedEventLog := getEventLog(OrderBookContractAddress, longOrderAcceptedEventTopics, longOrderAcceptedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderAcceptedEventLog})

		// orders exist in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Placed)

		// fail matching
		orderMatchingError := getEventFromABI(orderBookABI, "OrderMatchingError")
		orderMatchingErrorTopics := []common.Hash{orderMatchingError.ID, longOrderId}
		orderMatchingErrorData, _ := orderMatchingError.Inputs.NonIndexed().Pack("INSUFFICIENT_MARGIN")
		orderMatchingErrorLog := getEventLog(OrderBookContractAddress, orderMatchingErrorTopics, orderMatchingErrorData, blockNumber.Uint64()+2)
		cep.ProcessEvents([]*types.Log{orderMatchingErrorLog})

		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Execution_Failed)

		// now un-fail it
		orderMatchingErrorLog.Removed = true
		cep.ProcessEvents([]*types.Log{orderMatchingErrorLog})
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Placed)
	})
}

func newcep(t *testing.T, db LimitOrderDatabase) *ContractEventsProcessor {
	return NewContractEventsProcessor(db)
}

func getABIfromJson(jsonBytes []byte) abi.ABI {
	returnedABI, err := abi.FromSolidityJson(string(jsonBytes))
	if err != nil {
		panic(err)
	}
	return returnedABI
}

func getEventFromABI(contractABI abi.ABI, eventName string) abi.Event {
	for _, event := range contractABI.Events {
		if event.Name == eventName {
			return event
		}
	}
	return abi.Event{}
}

func getLimitOrder(ammIndex *big.Int, traderAddress common.Address, baseAssetQuantity *big.Int, price *big.Int, salt *big.Int) LimitOrder {
	return LimitOrder{
		BaseOrder: BaseOrder{
			AmmIndex:          ammIndex,
			Trader:            traderAddress,
			BaseAssetQuantity: baseAssetQuantity,
			Price:             price,
			Salt:              salt,
			ReduceOnly:        false,
		},
		PostOnly: false,
	}
}

func getEventLog(contractAddress common.Address, topics []common.Hash, eventData []byte, blockNumber uint64) *types.Log {
	return &types.Log{
		Address:     contractAddress,
		Topics:      topics,
		Data:        eventData,
		BlockNumber: blockNumber,
	}
}

// @todo change this to return the EIP712 hash instead
func getIdFromOrder(order Order) common.Hash {
	return crypto.Keccak256Hash([]byte(order.Trader.String() + order.Salt.String()))
}

// @todo change this to return the EIP712 hash instead
func getIdFromLimitOrder(order LimitOrder) common.Hash {
	return crypto.Keccak256Hash([]byte(order.Trader.String() + order.Salt.String()))
}
