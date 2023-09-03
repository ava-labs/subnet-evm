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
		orderPlacedEvent := getEventFromABI(limitOrderBookABI, "OrderPlaced")
		longOrderPlacedEventTopics := []common.Hash{orderPlacedEvent.ID, traderAddress.Hash(), longOrderId}
		longOrderPlacedEventData, err := orderPlacedEvent.Inputs.NonIndexed().Pack(longOrder, timestamp)
		if err != nil {
			t.Fatalf("%s", err)
		}
		longOrderPlacedEventLog := getEventLog(OrderBookContractAddress, longOrderPlacedEventTopics, longOrderPlacedEventData, ordersPlacedBlockNumber)

		shortOrderPlacedEventTopics := []common.Hash{orderPlacedEvent.ID, traderAddress.Hash(), shortOrderId}
		shortOrderPlacedEventData, err := orderPlacedEvent.Inputs.NonIndexed().Pack(shortOrder, timestamp)
		if err != nil {
			t.Fatalf("%s", err)
		}
		shortOrderPlacedEventLog := getEventLog(OrderBookContractAddress, shortOrderPlacedEventTopics, shortOrderPlacedEventData, ordersPlacedBlockNumber)

		ordersMatchedBlockNumber := uint64(14)
		ordersMatchedEvent := getEventFromABI(orderBookABI, "OrdersMatched")
		ordersMatchedEventTopics := []common.Hash{ordersMatchedEvent.ID, longOrderId, shortOrderId}
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		fillAmount := big.NewInt(3000000000000000000)
		fmt.Printf("sending matched event %s and %s", longOrderId.String(), shortOrderId.String())
		ordersMatchedEventData, _ := ordersMatchedEvent.Inputs.NonIndexed().Pack(fillAmount, price, big.NewInt(0), relayer, timestamp)
		ordersMatchedEventLog := getEventLog(OrderBookContractAddress, ordersMatchedEventTopics, ordersMatchedEventData, ordersMatchedBlockNumber)
		cep.ProcessEvents([]*types.Log{longOrderPlacedEventLog, shortOrderPlacedEventLog, ordersMatchedEventLog})
		// changed from the following which waws meaning to test the sorted-ness of the events before processing
		// cep.ProcessEvents([]*types.Log{ordersMatchedEventLog, longOrderPlacedEventLog, shortOrderPlacedEventLog})

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
	orderBookEvent := getEventFromABI(limitOrderBookABI, "OrderPlaced")
	orderPlacedEventData, _ := orderBookEvent.Inputs.NonIndexed().Pack(order, timestamp)
	orderBookEventTopics := []common.Hash{orderBookEvent.ID, traderAddress.Hash(), getIdFromLimitOrder(order)}
	orderBookLog := getEventLog(OrderBookContractAddress, orderBookEventTopics, orderPlacedEventData, blockNumber)

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

	//OrderBook log - OrderPlaced
	actualLimitOrder := *db.GetOrderBookData().OrderMap[getIdFromLimitOrder(order)]
	args := map[string]interface{}{}
	orderBookABI.UnpackIntoMap(args, "OrderPlaced", orderPlacedEventData)
	assert.Equal(t, Market(ammIndex.Int64()), actualLimitOrder.Market)
	assert.Equal(t, LONG, actualLimitOrder.PositionType)
	assert.Equal(t, traderAddress.String(), actualLimitOrder.UserAddress)
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
	order := getLimitOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt)
	blockNumber := uint64(12)
	orderBookABI := getABIfromJson(abis.OrderBookAbi)

	t.Run("When event is orderPlaced", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrderPlaced")
		topics := []common.Hash{event.ID, traderAddress.Hash(), getIdFromLimitOrder(order)}
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderPlacedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, orderPlacedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[getIdFromLimitOrder(order)]
			assert.Nil(t, actualLimitOrder)
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			orderId := crypto.Keccak256Hash([]byte(order.Trader.String() + order.Salt.String()))
			orderPlacedEventData, err := event.Inputs.NonIndexed().Pack(order, timestamp)
			if err != nil {
				t.Fatalf("%s", err)
			}
			log := getEventLog(OrderBookContractAddress, topics, orderPlacedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})

			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			args := map[string]interface{}{}
			orderBookABI.UnpackIntoMap(args, "OrderPlaced", orderPlacedEventData)
			assert.Equal(t, Market(ammIndex.Int64()), actualLimitOrder.Market)
			assert.Equal(t, LONG, actualLimitOrder.PositionType)
			assert.Equal(t, traderAddress.String(), actualLimitOrder.UserAddress)
			assert.Equal(t, *baseAssetQuantity, *actualLimitOrder.BaseAssetQuantity)
			assert.Equal(t, *price, *actualLimitOrder.Price)
			assert.Equal(t, Placed, actualLimitOrder.getOrderStatus().Status)
			assert.Equal(t, big.NewInt(int64(blockNumber)), actualLimitOrder.BlockNumber)
			rawOrder := &LimitOrder{}
			rawOrder.DecodeFromRawOrder(args["order"])
			assert.Equal(t, rawOrder, actualLimitOrder.RawOrder.(*LimitOrder))
		})
	})

	t.Run("When event is OrderAccepted", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrderAccepted")
		orderV2 := getLimitOrderV2(ammIndex, traderAddress, baseAssetQuantity, price, salt)
		orderId := getIdFromLimitOrderV2(orderV2)
		topics := []common.Hash{event.ID, traderAddress.Hash(), orderId}
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderAcceptedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, orderAcceptedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Nil(t, actualLimitOrder)
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			orderAcceptedEventData, err := event.Inputs.NonIndexed().Pack(orderV2, timestamp)
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
			assert.Equal(t, traderAddress.String(), actualLimitOrder.UserAddress)
			assert.Equal(t, *baseAssetQuantity, *actualLimitOrder.BaseAssetQuantity)
			assert.Equal(t, false, actualLimitOrder.ReduceOnly)
			assert.Equal(t, false, actualLimitOrder.isPostOnly())
			assert.Equal(t, *price, *actualLimitOrder.Price)
			assert.Equal(t, Placed, actualLimitOrder.getOrderStatus().Status)
			assert.Equal(t, big.NewInt(int64(blockNumber)), actualLimitOrder.BlockNumber)
			rawOrder := &LimitOrderV2{}
			rawOrder.DecodeFromRawOrder(args["order"])
			assert.Equal(t, rawOrder, actualLimitOrder.RawOrder.(*LimitOrderV2))
		})
	})
	t.Run("When event is OrderCancelled", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrderCancelled")
		topics := []common.Hash{event.ID, traderAddress.Hash(), getIdFromLimitOrder(order)}
		blockNumber := uint64(4)
		limitOrder := &Order{
			Market:            Market(ammIndex.Int64()),
			PositionType:      LONG,
			UserAddress:       traderAddress.String(),
			BaseAssetQuantity: baseAssetQuantity,
			Price:             price,
			BlockNumber:       big.NewInt(1),
			Salt:              salt,
		}
		limitOrder.Id = getIdFromOrder(*limitOrder)
		db.Add(limitOrder)
		// t.Run("When data in log unpack fails", func(t *testing.T) {
		// 	orderCancelledEventData := []byte{}
		// 	log := getEventLog(OrderBookContractAddress, topics, orderCancelledEventData, blockNumber)
		// 	cep.ProcessEvents([]*types.Log{log})
		// 	actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
		// 	assert.Equal(t, limitOrder, actualLimitOrder)
		// })
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			orderCancelledEventData, _ := event.Inputs.NonIndexed().Pack(timestamp)
			log := getEventLog(OrderBookContractAddress, topics, orderCancelledEventData, blockNumber)
			orderId := getIdFromOrder(*limitOrder)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Equal(t, Cancelled, actualLimitOrder.getOrderStatus().Status)
		})
	})
	t.Run("When event is OrderCancelAccepted", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		orderV2 := getLimitOrderV2(ammIndex, traderAddress, baseAssetQuantity, price, salt)
		event := getEventFromABI(orderBookABI, "OrderCancelAccepted")
		topics := []common.Hash{event.ID, traderAddress.Hash(), getIdFromLimitOrderV2(orderV2)}
		blockNumber := uint64(4)
		limitOrder := &Order{
			Market:            Market(ammIndex.Int64()),
			PositionType:      LONG,
			UserAddress:       traderAddress.String(),
			BaseAssetQuantity: baseAssetQuantity,
			Price:             price,
			BlockNumber:       big.NewInt(1),
			Salt:              salt,
			RawOrder:          &orderV2,
		}
		limitOrder.Id = getIdFromOrder(*limitOrder)
		db.Add(limitOrder)
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderCancelledEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, orderCancelledEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			orderId := getIdFromOrder(*limitOrder)
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Equal(t, limitOrder, actualLimitOrder)
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			orderCancelledEventData, _ := event.Inputs.NonIndexed().Pack(timestamp)
			log := getEventLog(OrderBookContractAddress, topics, orderCancelledEventData, blockNumber)
			orderId := getIdFromOrder(*limitOrder)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Equal(t, Cancelled, actualLimitOrder.getOrderStatus().Status)
		})
	})
	t.Run("When event is OrderMatched", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrdersMatched")
		longOrder := &Order{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            LONG,
			UserAddress:             traderAddress.String(),
			BaseAssetQuantity:       baseAssetQuantity,
			Price:                   price,
			BlockNumber:             big.NewInt(1),
			FilledBaseAssetQuantity: big.NewInt(0),
			Salt:                    salt,
		}
		shortOrder := &Order{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            SHORT,
			UserAddress:             traderAddress.String(),
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
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		fillAmount := big.NewInt(10)
		topics := []common.Hash{event.ID, longOrder.Id, shortOrder.Id}
		t.Run("When data in log unpack fails", func(t *testing.T) {
			ordersMatchedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, ordersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, int64(0), longOrder.FilledBaseAssetQuantity.Int64())
			assert.Equal(t, int64(0), shortOrder.FilledBaseAssetQuantity.Int64())
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			ordersMatchedEventData, _ := event.Inputs.NonIndexed().Pack(fillAmount, price, big.NewInt(0).Mul(fillAmount, price), relayer, timestamp)
			log := getEventLog(OrderBookContractAddress, topics, ordersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, big.NewInt(fillAmount.Int64()), longOrder.FilledBaseAssetQuantity)
			assert.Equal(t, big.NewInt(-fillAmount.Int64()), shortOrder.FilledBaseAssetQuantity)
		})
	})
	t.Run("When event is LiquidationOrderMatched", func(t *testing.T) {
		db := getDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "LiquidationOrderMatched")
		longOrder := &Order{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            LONG,
			UserAddress:             traderAddress.String(),
			BaseAssetQuantity:       baseAssetQuantity,
			Price:                   price,
			Salt:                    salt,
			BlockNumber:             big.NewInt(1),
			FilledBaseAssetQuantity: big.NewInt(0),
		}
		longOrder.Id = getIdFromOrder(*longOrder)
		db.Add(longOrder)
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		fillAmount := big.NewInt(10)
		topics := []common.Hash{event.ID, traderAddress.Hash(), longOrder.Id}
		t.Run("When data in log unpack fails", func(t *testing.T) {
			liquidationOrdersMatchedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, liquidationOrdersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[longOrder.Id]
			assert.Equal(t, longOrder, actualLimitOrder)
		})
		t.Run("When data in log unpack succeeds", func(t *testing.T) {
			// order := getOrder(ammIndex, traderAddress, longOrder.BaseAssetQuantity, price, salt)
			liquidationOrdersMatchedEventData, _ := event.Inputs.NonIndexed().Pack(fillAmount, price, big.NewInt(0).Mul(fillAmount, price), relayer, timestamp)
			log := getEventLog(OrderBookContractAddress, topics, liquidationOrdersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, fillAmount, db.OrderMap[longOrder.Id].FilledBaseAssetQuantity)
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

	orderPlacedEvent := getEventFromABI(limitOrderrderBookABI, "OrderPlaced")
	longOrder := getLimitOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt1)
	longOrderId := getIdFromLimitOrder(longOrder)
	longOrderPlacedEventTopics := []common.Hash{orderPlacedEvent.ID, traderAddress.Hash(), longOrderId}
	longOrderPlacedEventData, _ := orderPlacedEvent.Inputs.NonIndexed().Pack(longOrder, timestamp)

	shortOrder := getLimitOrder(ammIndex, traderAddress2, big.NewInt(0).Neg(baseAssetQuantity), price, salt2)
	shortOrderId := getIdFromLimitOrder(shortOrder)
	shortOrderPlacedEventTopics := []common.Hash{orderPlacedEvent.ID, traderAddress2.Hash(), shortOrderId}
	shortOrderPlacedEventData, _ := orderPlacedEvent.Inputs.NonIndexed().Pack(shortOrder, timestamp)

	t.Run("delete order when OrderPlaced is removed", func(t *testing.T) {
		longOrderPlacedEventLog := getEventLog(OrderBookContractAddress, longOrderPlacedEventTopics, longOrderPlacedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderPlacedEventLog})

		// order exists in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)

		// order should be deleted if OrderPlaced log is removed
		longOrderPlacedEventLog.Removed = true
		cep.ProcessEvents([]*types.Log{longOrderPlacedEventLog})
		assert.Nil(t, db.OrderMap[longOrderId])
	})

	t.Run("un-cancel an order when OrderCancelled is removed", func(t *testing.T) {
		longOrderPlacedEventLog := getEventLog(OrderBookContractAddress, longOrderPlacedEventTopics, longOrderPlacedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderPlacedEventLog})

		// order exists in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)

		// cancel it
		orderCancelledEvent := getEventFromABI(limitOrderrderBookABI, "OrderCancelled")
		orderCancelledEventTopics := []common.Hash{orderCancelledEvent.ID, traderAddress.Hash(), longOrderId}
		orderCancelledEventData, _ := orderCancelledEvent.Inputs.NonIndexed().Pack(timestamp)
		orderCancelledLog := getEventLog(OrderBookContractAddress, orderCancelledEventTopics, orderCancelledEventData, blockNumber.Uint64()+2)
		cep.ProcessEvents([]*types.Log{orderCancelledLog})

		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Cancelled)

		// now uncancel it
		orderCancelledLog.Removed = true
		cep.ProcessEvents([]*types.Log{orderCancelledLog})
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Placed)
	})

	t.Run("un-fulfill an order when OrdersMatched is removed", func(t *testing.T) {
		longOrderPlacedEventLog := getEventLog(OrderBookContractAddress, longOrderPlacedEventTopics, longOrderPlacedEventData, blockNumber.Uint64())
		shortOrderPlacedEventLog := getEventLog(OrderBookContractAddress, shortOrderPlacedEventTopics, shortOrderPlacedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderPlacedEventLog, shortOrderPlacedEventLog})

		// orders exist in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)
		assert.Equal(t, db.OrderMap[shortOrderId].Salt, shortOrder.Salt)

		// fulfill them
		ordersMatchedEvent := getEventFromABI(orderBookABI, "OrdersMatched")
		ordersMatchedEventTopics := []common.Hash{ordersMatchedEvent.ID, longOrderId, shortOrderId}
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		ordersMatchedEventData, _ := ordersMatchedEvent.Inputs.NonIndexed().Pack(baseAssetQuantity, price, big.NewInt(0).Mul(baseAssetQuantity, price), relayer, timestamp)
		ordersMatchedLog := getEventLog(OrderBookContractAddress, ordersMatchedEventTopics, ordersMatchedEventData, blockNumber.Uint64()+2)
		cep.ProcessEvents([]*types.Log{ordersMatchedLog})

		assert.Equal(t, db.OrderMap[shortOrderId].getOrderStatus().Status, FulFilled)
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, FulFilled)

		// now un-fulfill it
		ordersMatchedLog.Removed = true
		cep.ProcessEvents([]*types.Log{ordersMatchedLog})
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Placed)
		assert.Equal(t, db.OrderMap[shortOrderId].getOrderStatus().Status, Placed)
	})

	t.Run("un-fulfill an order when LiquidationOrderMatched is removed", func(t *testing.T) {
		// change salt to create a new order in memory
		longOrder.Salt.Add(longOrder.Salt, big.NewInt(10))
		longOrderId = getIdFromLimitOrder(longOrder)
		longOrderPlacedEventTopics = []common.Hash{orderPlacedEvent.ID, traderAddress.Hash(), longOrderId}
		longOrderPlacedEventData, _ = orderPlacedEvent.Inputs.NonIndexed().Pack(longOrder, timestamp)
		longOrderPlacedEventLog := getEventLog(OrderBookContractAddress, longOrderPlacedEventTopics, longOrderPlacedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderPlacedEventLog})

		// orders exist in memory now
		assert.Equal(t, db.OrderMap[longOrderId].Salt, longOrder.Salt)

		// fulfill
		liquidationOrderMatchedEvent := getEventFromABI(orderBookABI, "LiquidationOrderMatched")
		liquidationOrderMatchedEventTopics := []common.Hash{liquidationOrderMatchedEvent.ID, traderAddress.Hash(), longOrderId}
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		liquidationOrderMatchedEventData, _ := liquidationOrderMatchedEvent.Inputs.NonIndexed().Pack(baseAssetQuantity, price, big.NewInt(0).Mul(baseAssetQuantity, price), relayer, timestamp)
		liquidationOrderMatchedLog := getEventLog(OrderBookContractAddress, liquidationOrderMatchedEventTopics, liquidationOrderMatchedEventData, blockNumber.Uint64()+2)
		cep.ProcessEvents([]*types.Log{liquidationOrderMatchedLog})

		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, FulFilled)

		// now un-fulfill it
		liquidationOrderMatchedLog.Removed = true
		cep.ProcessEvents([]*types.Log{liquidationOrderMatchedLog})
		assert.Equal(t, db.OrderMap[longOrderId].getOrderStatus().Status, Placed)
	})

	t.Run("revert state of an order when OrderMatchingError is removed", func(t *testing.T) {
		// change salt
		longOrder.Salt.Add(longOrder.Salt, big.NewInt(20))
		longOrderId = getIdFromLimitOrder(longOrder)
		longOrderPlacedEventTopics = []common.Hash{orderPlacedEvent.ID, traderAddress.Hash(), longOrderId}
		longOrderPlacedEventData, _ = orderPlacedEvent.Inputs.NonIndexed().Pack(longOrder, timestamp)
		longOrderPlacedEventLog := getEventLog(OrderBookContractAddress, longOrderPlacedEventTopics, longOrderPlacedEventData, blockNumber.Uint64())
		cep.ProcessEvents([]*types.Log{longOrderPlacedEventLog})

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
		AmmIndex:          ammIndex,
		Trader:            traderAddress,
		BaseAssetQuantity: baseAssetQuantity,
		Price:             price,
		Salt:              salt,
		ReduceOnly:        false,
	}
}

func getLimitOrderV2(ammIndex *big.Int, traderAddress common.Address, baseAssetQuantity *big.Int, price *big.Int, salt *big.Int) LimitOrderV2 {
	return LimitOrderV2{
		LimitOrder: LimitOrder{
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
func getIdFromLimitOrder(order LimitOrder) common.Hash {
	return crypto.Keccak256Hash([]byte(order.Trader.String() + order.Salt.String()))
}

// @todo change this to return the EIP712 hash instead
func getIdFromOrder(order Order) common.Hash {
	return crypto.Keccak256Hash([]byte(order.UserAddress + order.Salt.String()))
}

// @todo change this to return the EIP712 hash instead
func getIdFromLimitOrderV2(order LimitOrderV2) common.Hash {
	return crypto.Keccak256Hash([]byte(order.Trader.String() + order.Salt.String()))
}
