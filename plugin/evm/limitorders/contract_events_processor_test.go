package limitorders

import (
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var orderBookTestContractFileLocation = "../../../contract-examples/artifacts/contracts/hubble-v2/OrderBook.sol/OrderBook.json"
var marginAccountTestContractFileLocation = "../../../contract-examples/artifacts/contracts/hubble-v2/MarginAccount.sol/MarginAccount.json"
var clearingHouseTestContractFileLocation = "../../../contract-examples/artifacts/contracts/hubble-v2/ClearingHouse.sol/ClearingHouse.json"

func TestProcessEvents(t *testing.T) {
	t.Run("when events are from orderbook contract", func(t *testing.T) {
		TestHandleOrderBookEvent(t)
	})
	t.Run("when events are from marginAccount contract", func(t *testing.T) {
		TestHandleMarginAccountEvent(t)
	})
	t.Run("when events are from clearinghouse contract", func(t *testing.T) {
		TestHandleClearingHouseEvent(t)
	})
	t.Run("when events has one orderbook, one marginaccount, one clearing house events", func(t *testing.T) {
		TestOrderBookMarginAccountClearingHouseEventInLog(t)
	})
	t.Run("it sorts events by blockNumber and executes in order", func(t *testing.T) {
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		orderBookABI := getABIfromJson(orderBookTestContractFileLocation)

		traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
		ammIndex := big.NewInt(0)
		baseAssetQuantity = big.NewInt(5000000000000000000)
		price := big.NewInt(1000000000)
		longSignature := []byte("long")
		salt1 := big.NewInt(1675239557437)
		longOrder := getOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt1)
		salt2 := big.NewInt(0).Add(salt1, big.NewInt(1))
		shortOrder := getOrder(ammIndex, traderAddress, big.NewInt(0).Neg(baseAssetQuantity), price, salt2)
		ordersPlacedBlockNumber := uint64(12)
		orderPlacedEvent := getEventFromABI(orderBookABI, "OrderPlaced")
		orderPlacedEventTopics := []common.Hash{orderPlacedEvent.ID, traderAddress.Hash()}
		longOrderPlacedEventData, _ := orderPlacedEvent.Inputs.NonIndexed().Pack(longOrder, longSignature)
		longOrderPlacedEventLog := getEventLog(OrderBookContractAddress, orderPlacedEventTopics, longOrderPlacedEventData, ordersPlacedBlockNumber)

		shortSignature := []byte("short")
		shortOrderPlacedEventData, _ := orderPlacedEvent.Inputs.NonIndexed().Pack(shortOrder, shortSignature)
		shortOrderPlacedEventLog := getEventLog(OrderBookContractAddress, orderPlacedEventTopics, shortOrderPlacedEventData, ordersPlacedBlockNumber)

		ordersMatchedBlockNumber := uint64(14)
		ordersMatchedEvent := getEventFromABI(orderBookABI, "OrdersMatched")
		ordersMatchedEventTopics := []common.Hash{ordersMatchedEvent.ID}
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		fillAmount := big.NewInt(3000000000000000000)
		orders := []Order{longOrder, shortOrder}
		signatures := [][]byte{longSignature, shortSignature}
		ordersMatchedEventData, _ := ordersMatchedEvent.Inputs.NonIndexed().Pack(orders, signatures, fillAmount, relayer)
		ordersMatchedEventLog := getEventLog(OrderBookContractAddress, ordersMatchedEventTopics, ordersMatchedEventData, ordersMatchedBlockNumber)
		cep.ProcessEvents([]*types.Log{ordersMatchedEventLog, longOrderPlacedEventLog, shortOrderPlacedEventLog})

		actualLongOrder := db.OrderMap[getIdFromOrder(longOrder)]
		assert.Equal(t, fillAmount, actualLongOrder.FilledBaseAssetQuantity)

		actualShortOrder := db.OrderMap[getIdFromOrder(shortOrder)]
		assert.Equal(t, big.NewInt(0).Neg(fillAmount), actualShortOrder.FilledBaseAssetQuantity)
	})
	t.Run("when event is removed it is not processed", func(t *testing.T) {
		traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
		db := NewInMemoryDatabase()
		collateral := HUSD
		originalMargin := multiplyBasePrecision(big.NewInt(100))
		trader := &Trader{
			Margins: map[Collateral]*big.Int{collateral: big.NewInt(0).Set(originalMargin)},
		}
		db.TraderMap[traderAddress] = trader
		blockNumber := uint64(12)

		//MarginAccount Contract log
		marginAccountABI := getABIfromJson(marginAccountTestContractFileLocation)
		marginAccountEvent := getEventFromABI(marginAccountABI, "MarginAdded")
		marginAccountEventTopics := []common.Hash{marginAccountEvent.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
		marginAdded := multiplyBasePrecision(big.NewInt(100))
		timeStamp := big.NewInt(time.Now().Unix())
		marginAddedEventData, _ := marginAccountEvent.Inputs.NonIndexed().Pack(marginAdded, timeStamp)
		marginAddedLog := getEventLog(MarginAccountContractAddress, marginAccountEventTopics, marginAddedEventData, blockNumber)
		marginAddedLog.Removed = true
		cep := newcep(t, db)

		cep.ProcessEvents([]*types.Log{marginAddedLog})
		assert.Equal(t, originalMargin, db.TraderMap[traderAddress].Margins[collateral])
	})
}

func TestOrderBookMarginAccountClearingHouseEventInLog(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	blockNumber := uint64(12)
	db := NewInMemoryDatabase()
	cep := newcep(t, db)
	collateral := HUSD
	openNotional := multiplyBasePrecision(big.NewInt(100))
	size := multiplyPrecisionSize(big.NewInt(10))
	lastPremiumFraction := multiplyBasePrecision(big.NewInt(1))
	liquidationThreshold := multiplyBasePrecision(big.NewInt(1))
	unrealisedFunding := multiplyBasePrecision(big.NewInt(1))
	market := AvaxPerp
	position := &Position{
		OpenNotional:         openNotional,
		Size:                 size,
		UnrealisedFunding:    unrealisedFunding,
		LastPremiumFraction:  lastPremiumFraction,
		LiquidationThreshold: liquidationThreshold,
	}
	originalMargin := multiplyBasePrecision(big.NewInt(100))
	trader := &Trader{
		Margins:   map[Collateral]*big.Int{collateral: big.NewInt(0).Set(originalMargin)},
		Positions: map[Market]*Position{market: position},
	}
	db.TraderMap[traderAddress] = trader

	//OrderBook Contract log
	ammIndex := big.NewInt(0)
	baseAssetQuantity = big.NewInt(5000000000000000000)
	price := big.NewInt(1000000000)
	salt := big.NewInt(1675239557437)
	signature := []byte("signature")
	order := getOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt)
	orderBookABI := getABIfromJson(orderBookTestContractFileLocation)
	orderBookEvent := getEventFromABI(orderBookABI, "OrderPlaced")
	orderPlacedEventData, _ := orderBookEvent.Inputs.NonIndexed().Pack(order, signature)
	orderBookEventTopics := []common.Hash{orderBookEvent.ID, traderAddress.Hash()}
	orderBookLog := getEventLog(OrderBookContractAddress, orderBookEventTopics, orderPlacedEventData, blockNumber)

	//MarginAccount Contract log
	marginAccountABI := getABIfromJson(marginAccountTestContractFileLocation)
	marginAccountEvent := getEventFromABI(marginAccountABI, "MarginAdded")
	marginAccountEventTopics := []common.Hash{marginAccountEvent.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
	marginAdded := multiplyBasePrecision(big.NewInt(100))
	timeStamp := big.NewInt(time.Now().Unix())
	marginAddedEventData, _ := marginAccountEvent.Inputs.NonIndexed().Pack(marginAdded, timeStamp)
	marginAccountLog := getEventLog(MarginAccountContractAddress, marginAccountEventTopics, marginAddedEventData, blockNumber)

	//ClearingHouse Contract log
	clearingHouseABI := getABIfromJson(clearingHouseTestContractFileLocation)
	clearingHouseEvent := getEventFromABI(clearingHouseABI, "FundingRateUpdated")
	clearingHouseEventTopics := []common.Hash{clearingHouseEvent.ID, common.BigToHash(big.NewInt(int64(market)))}

	nextFundingTime := big.NewInt(time.Now().Unix())
	premiumFraction := multiplyBasePrecision(big.NewInt(10))
	underlyingPrice := multiplyBasePrecision(big.NewInt(100))
	cumulativePremiumFraction := multiplyBasePrecision(big.NewInt(10))
	timestamp := big.NewInt(time.Now().Unix())
	fundingRateUpdated, _ := clearingHouseEvent.Inputs.NonIndexed().Pack(premiumFraction, underlyingPrice, cumulativePremiumFraction, nextFundingTime, timestamp, big.NewInt(int64(blockNumber)))
	clearingHouseLog := getEventLog(ClearingHouseContractAddress, clearingHouseEventTopics, fundingRateUpdated, blockNumber)

	logs := []*types.Log{orderBookLog, marginAccountLog, clearingHouseLog}
	cep.ProcessEvents(logs)

	//OrderBook log - OrderPlaced
	actualLimitOrder := *db.GetOrderBookData().OrderMap[getOrderId(traderAddress, salt)]
	args := map[string]interface{}{}
	orderBookABI.UnpackIntoMap(args, "OrderPlaced", orderPlacedEventData)
	assert.Equal(t, Market(ammIndex.Int64()), actualLimitOrder.Market)
	assert.Equal(t, "long", actualLimitOrder.PositionType)
	assert.Equal(t, traderAddress.String(), actualLimitOrder.UserAddress)
	assert.Equal(t, *baseAssetQuantity, *actualLimitOrder.BaseAssetQuantity)
	assert.Equal(t, *price, *actualLimitOrder.Price)
	assert.Equal(t, Status("placed"), actualLimitOrder.Status)
	assert.Equal(t, signature, actualLimitOrder.Signature)
	assert.Equal(t, big.NewInt(int64(blockNumber)), actualLimitOrder.BlockNumber)
	assert.Equal(t, args["order"], actualLimitOrder.RawOrder)

	//ClearingHouse log - FundingRateUpdated
	expectedUnrealisedFunding := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size))
	assert.Equal(t, expectedUnrealisedFunding, db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)

	//MarginAccount log - marginAdded
	actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margins[collateral]
	assert.Equal(t, big.NewInt(0).Add(marginAdded, originalMargin), actualMargin)

}

func TestHandleOrderBookEvent(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	ammIndex := big.NewInt(0)
	baseAssetQuantity = big.NewInt(5000000000000000000)
	price := big.NewInt(1000000000)
	salt := big.NewInt(1675239557437)
	signature := []byte("signature")
	order := getOrder(ammIndex, traderAddress, baseAssetQuantity, price, salt)
	blockNumber := uint64(12)
	orderBookABI := getABIfromJson(orderBookTestContractFileLocation)

	t.Run("When event is orderPlaced", func(t *testing.T) {
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrderPlaced")
		topics := []common.Hash{event.ID, traderAddress.Hash()}
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderPlacedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, orderPlacedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[getOrderId(traderAddress, salt)]
			assert.Nil(t, actualLimitOrder)
		})
		t.Run("When data in log unpack suceeds", func(t *testing.T) {
			orderPlacedEventData, _ := event.Inputs.NonIndexed().Pack(order, signature)
			log := getEventLog(OrderBookContractAddress, topics, orderPlacedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})

			actualLimitOrder := *db.GetOrderBookData().OrderMap[getOrderId(traderAddress, salt)]
			args := map[string]interface{}{}
			orderBookABI.UnpackIntoMap(args, "OrderPlaced", orderPlacedEventData)
			assert.Equal(t, Market(ammIndex.Int64()), actualLimitOrder.Market)
			assert.Equal(t, "long", actualLimitOrder.PositionType)
			assert.Equal(t, traderAddress.String(), actualLimitOrder.UserAddress)
			assert.Equal(t, *baseAssetQuantity, *actualLimitOrder.BaseAssetQuantity)
			assert.Equal(t, *price, *actualLimitOrder.Price)
			assert.Equal(t, Status("placed"), actualLimitOrder.Status)
			assert.Equal(t, signature, actualLimitOrder.Signature)
			assert.Equal(t, big.NewInt(int64(blockNumber)), actualLimitOrder.BlockNumber)
			assert.Equal(t, args["order"], actualLimitOrder.RawOrder)
		})
	})
	t.Run("When event is OrderCancelled", func(t *testing.T) {
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrderCancelled")
		topics := []common.Hash{event.ID, traderAddress.Hash()}
		blockNumber := uint64(4)
		limitOrder := &LimitOrder{
			Market:            Market(ammIndex.Int64()),
			PositionType:      "long",
			UserAddress:       traderAddress.String(),
			BaseAssetQuantity: baseAssetQuantity,
			Price:             price,
			Status:            Placed,
			Signature:         signature,
			BlockNumber:       big.NewInt(1),
			Salt:              salt,
		}
		db.Add(limitOrder)
		t.Run("When data in log unpack fails", func(t *testing.T) {
			orderCancelledEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, orderCancelledEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[getIdFromLimitOrder(*limitOrder)]
			assert.Equal(t, limitOrder, actualLimitOrder)
		})
		t.Run("When data in log unpack suceeds", func(t *testing.T) {
			orderCancelledEventData, _ := event.Inputs.NonIndexed().Pack(order)
			log := getEventLog(OrderBookContractAddress, topics, orderCancelledEventData, blockNumber)
			orderId := getIdFromLimitOrder(*limitOrder)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[orderId]
			assert.Nil(t, actualLimitOrder)
		})
	})
	t.Run("When event is OrderMatched", func(t *testing.T) {
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "OrdersMatched")
		topics := []common.Hash{event.ID}
		signature1 := []byte("longOrder")
		signature2 := []byte("shortOrder")
		longOrder := &LimitOrder{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            "long",
			UserAddress:             traderAddress.String(),
			BaseAssetQuantity:       baseAssetQuantity,
			Price:                   price,
			Status:                  Placed,
			Signature:               signature1,
			BlockNumber:             big.NewInt(1),
			FilledBaseAssetQuantity: big.NewInt(0),
			Salt:                    salt,
		}
		shortOrder := &LimitOrder{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            "short",
			UserAddress:             traderAddress.String(),
			BaseAssetQuantity:       big.NewInt(0).Mul(baseAssetQuantity, big.NewInt(-1)),
			Price:                   price,
			Status:                  Placed,
			Signature:               signature2,
			BlockNumber:             big.NewInt(1),
			FilledBaseAssetQuantity: big.NewInt(0),
			Salt:                    big.NewInt(0).Add(salt, big.NewInt(1000)),
		}
		db.Add(longOrder)
		db.Add(shortOrder)
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		fillAmount := big.NewInt(10)
		t.Run("When data in log unpack fails", func(t *testing.T) {
			ordersMatchedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, ordersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, int64(0), longOrder.FilledBaseAssetQuantity.Int64())
			assert.Equal(t, int64(0), shortOrder.FilledBaseAssetQuantity.Int64())
		})
		t.Run("When data in log unpack suceeds", func(t *testing.T) {
			order1 := getOrderFromLimitOrder(*longOrder)
			order2 := getOrderFromLimitOrder(*shortOrder)
			orders := []Order{order1, order2}
			signatures := [][]byte{signature1, signature2}
			ordersMatchedEventData, _ := event.Inputs.NonIndexed().Pack(orders, signatures, fillAmount, relayer)
			log := getEventLog(OrderBookContractAddress, topics, ordersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, big.NewInt(fillAmount.Int64()), longOrder.FilledBaseAssetQuantity)
			assert.Equal(t, big.NewInt(-fillAmount.Int64()), shortOrder.FilledBaseAssetQuantity)
		})
	})
	t.Run("When event is LiquidationOrderMatched", func(t *testing.T) {
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(orderBookABI, "LiquidationOrderMatched")
		topics := []common.Hash{event.ID, traderAddress.Hash()}
		signature := []byte("longOrder")
		longOrder := &LimitOrder{
			Market:                  Market(ammIndex.Int64()),
			PositionType:            "long",
			UserAddress:             traderAddress.String(),
			BaseAssetQuantity:       baseAssetQuantity,
			Price:                   price,
			Status:                  Placed,
			Signature:               signature,
			Salt:                    salt,
			BlockNumber:             big.NewInt(1),
			FilledBaseAssetQuantity: big.NewInt(0),
		}
		db.Add(longOrder)
		relayer := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		fillAmount := big.NewInt(10)
		t.Run("When data in log unpack fails", func(t *testing.T) {
			ordersMatchedEventData := []byte{}
			log := getEventLog(OrderBookContractAddress, topics, ordersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualLimitOrder := db.GetOrderBookData().OrderMap[getIdFromLimitOrder(*longOrder)]
			assert.Equal(t, longOrder, actualLimitOrder)
		})
		t.Run("When data in log unpack suceeds", func(t *testing.T) {
			order := getOrder(ammIndex, traderAddress, longOrder.BaseAssetQuantity, price, salt)
			ordersMatchedEventData, _ := event.Inputs.NonIndexed().Pack(order, signature, fillAmount, relayer)
			log := getEventLog(OrderBookContractAddress, topics, ordersMatchedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, big.NewInt(fillAmount.Int64()), longOrder.FilledBaseAssetQuantity)
		})
	})
}

func TestHandleMarginAccountEvent(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	blockNumber := uint64(12)
	collateral := HUSD

	marginAccountABI := getABIfromJson(marginAccountTestContractFileLocation)

	t.Run("when event is MarginAdded", func(t *testing.T) {
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(marginAccountABI, "MarginAdded")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
		t.Run("When event parsing fails", func(t *testing.T) {
			marginAddedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, marginAddedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			marginAdded := big.NewInt(10000)
			timeStamp := big.NewInt(time.Now().Unix())
			marginAddedEventData, _ := event.Inputs.NonIndexed().Pack(marginAdded, timeStamp)
			log := getEventLog(MarginAccountContractAddress, topics, marginAddedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margins[collateral]
			assert.Equal(t, marginAdded, actualMargin)
		})
	})
	t.Run("when event is MarginRemoved", func(t *testing.T) {
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		event := getEventFromABI(marginAccountABI, "MarginRemoved")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(collateral)))}
		t.Run("When event parsing fails", func(t *testing.T) {
			marginRemovedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, marginRemovedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			marginRemoved := big.NewInt(10000)
			timeStamp := big.NewInt(time.Now().Unix())
			marginRemovedEventData, _ := event.Inputs.NonIndexed().Pack(marginRemoved, timeStamp)
			log := getEventLog(MarginAccountContractAddress, topics, marginRemovedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margins[collateral]
			assert.Equal(t, big.NewInt(0).Neg(marginRemoved), actualMargin)
		})
	})
	t.Run("when event is PnLRealized", func(t *testing.T) {
		event := getEventFromABI(marginAccountABI, "PnLRealized")
		topics := []common.Hash{event.ID, traderAddress.Hash()}
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		t.Run("When event parsing fails", func(t *testing.T) {
			pnlRealizedEventData := []byte{}
			log := getEventLog(MarginAccountContractAddress, topics, pnlRealizedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Nil(t, db.GetOrderBookData().TraderMap[traderAddress])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			pnlRealized := big.NewInt(-10000)
			timeStamp := big.NewInt(time.Now().Unix())
			pnlRealizedEventData, _ := event.Inputs.NonIndexed().Pack(pnlRealized, timeStamp)
			log := getEventLog(MarginAccountContractAddress, topics, pnlRealizedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			actualMargin := db.GetOrderBookData().TraderMap[traderAddress].Margins[collateral]
			assert.Equal(t, pnlRealized, actualMargin)
		})
	})
}
func TestHandleClearingHouseEvent(t *testing.T) {
	traderAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	blockNumber := uint64(12)
	collateral := HUSD
	market := AvaxPerp
	clearingHouseABI := getABIfromJson(clearingHouseTestContractFileLocation)
	openNotional := multiplyBasePrecision(big.NewInt(100))
	size := multiplyPrecisionSize(big.NewInt(10))
	lastPremiumFraction := multiplyBasePrecision(big.NewInt(1))
	liquidationThreshold := multiplyBasePrecision(big.NewInt(1))
	unrealisedFunding := multiplyBasePrecision(big.NewInt(1))
	t.Run("when event is FundingRateUpdated", func(t *testing.T) {
		event := getEventFromABI(clearingHouseABI, "FundingRateUpdated")
		topics := []common.Hash{event.ID, common.BigToHash(big.NewInt(int64(market)))}
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margins:   map[Collateral]*big.Int{collateral: big.NewInt(100)},
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
			timestamp := big.NewInt(time.Now().Unix())
			fundingRateUpdated, _ := event.Inputs.NonIndexed().Pack(premiumFraction, underlyingPrice, cumulativePremiumFraction, nextFundingTime, timestamp, big.NewInt(int64(blockNumber)))
			log := getEventLog(ClearingHouseContractAddress, topics, fundingRateUpdated, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			expectedUnrealisedFunding := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size))
			assert.Equal(t, expectedUnrealisedFunding, db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)
		})
	})
	t.Run("When event is FundingPaid", func(t *testing.T) {
		event := getEventFromABI(clearingHouseABI, "FundingPaid")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(market)))}
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margins:   map[Collateral]*big.Int{collateral: big.NewInt(100)},
			Positions: map[Market]*Position{market: position},
		}
		db.TraderMap[traderAddress] = trader

		t.Run("When event parsing fails", func(t *testing.T) {
			pnlRealizedEventData := []byte{}
			log := getEventLog(ClearingHouseContractAddress, topics, pnlRealizedEventData, blockNumber)
			cep.ProcessEvents([]*types.Log{log})

			assert.Equal(t, unrealisedFunding, db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)
			assert.Equal(t, lastPremiumFraction, db.TraderMap[traderAddress].Positions[market].LastPremiumFraction)
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			takerFundingPayment := multiplyBasePrecision(big.NewInt(10))
			cumulativePremiumFraction := multiplyBasePrecision(big.NewInt(10))
			fundingPaidEvent, _ := event.Inputs.NonIndexed().Pack(takerFundingPayment, cumulativePremiumFraction)
			log := getEventLog(ClearingHouseContractAddress, topics, fundingPaidEvent, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, big.NewInt(0), db.TraderMap[traderAddress].Positions[market].UnrealisedFunding)
			assert.Equal(t, cumulativePremiumFraction, db.TraderMap[traderAddress].Positions[market].LastPremiumFraction)
		})
	})
	t.Run("When event is PositionModified", func(t *testing.T) {
		event := getEventFromABI(clearingHouseABI, "PositionModified")
		topics := []common.Hash{event.ID, traderAddress.Hash(), common.BigToHash(big.NewInt(int64(market)))}
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margins:   map[Collateral]*big.Int{collateral: big.NewInt(100)},
			Positions: map[Market]*Position{market: position},
		}
		db.TraderMap[traderAddress] = trader

		t.Run("When event parsing fails", func(t *testing.T) {
			positionModifiedEvent := []byte{}
			log := getEventLog(ClearingHouseContractAddress, topics, positionModifiedEvent, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, big.NewInt(0), db.LastPrice[market])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			baseAsset := multiplyPrecisionSize(big.NewInt(10))
			quoteAsset := multiplyBasePrecision(big.NewInt(1000))
			realizedPnl := multiplyBasePrecision(big.NewInt(20))
			openNotional := multiplyBasePrecision(big.NewInt(4000))
			timestamp := multiplyBasePrecision(big.NewInt(time.Now().Unix()))
			size := multiplyPrecisionSize(big.NewInt(40))

			positionModifiedEvent, _ := event.Inputs.NonIndexed().Pack(baseAsset, quoteAsset, realizedPnl, size, openNotional, timestamp)
			log := getEventLog(ClearingHouseContractAddress, topics, positionModifiedEvent, blockNumber)
			cep.ProcessEvents([]*types.Log{log})

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
		db := NewInMemoryDatabase()
		cep := newcep(t, db)
		position := &Position{
			OpenNotional:         openNotional,
			Size:                 size,
			UnrealisedFunding:    unrealisedFunding,
			LastPremiumFraction:  lastPremiumFraction,
			LiquidationThreshold: liquidationThreshold,
		}
		trader := &Trader{
			Margins:   map[Collateral]*big.Int{collateral: big.NewInt(100)},
			Positions: map[Market]*Position{market: position},
		}
		db.TraderMap[traderAddress] = trader

		t.Run("When event parsing fails", func(t *testing.T) {
			positionLiquidatedEvent := []byte{}
			log := getEventLog(ClearingHouseContractAddress, topics, positionLiquidatedEvent, blockNumber)
			cep.ProcessEvents([]*types.Log{log})
			assert.Equal(t, big.NewInt(0), db.LastPrice[market])
		})
		t.Run("When event parsing succeeds", func(t *testing.T) {
			baseAsset := multiplyPrecisionSize(big.NewInt(10))
			quoteAsset := multiplyBasePrecision(big.NewInt(1000))
			realizedPnl := multiplyBasePrecision(big.NewInt(20))
			openNotional := multiplyBasePrecision(big.NewInt(4000))
			timestamp := multiplyBasePrecision(big.NewInt(time.Now().Unix()))
			size := multiplyPrecisionSize(big.NewInt(40))

			positionLiquidatedEvent, _ := event.Inputs.NonIndexed().Pack(baseAsset, quoteAsset, realizedPnl, size, openNotional, timestamp)
			log := getEventLog(ClearingHouseContractAddress, topics, positionLiquidatedEvent, blockNumber)
			cep.ProcessEvents([]*types.Log{log})

			// quoteAsset/(baseAsset / 1e 18)
			expectedLastPrice := big.NewInt(100000000)
			assert.Equal(t, expectedLastPrice, db.LastPrice[market])
			assert.Equal(t, size, db.TraderMap[traderAddress].Positions[market].Size)
			assert.Equal(t, openNotional, db.TraderMap[traderAddress].Positions[market].OpenNotional)
		})
	})
}

func newcep(t *testing.T, db LimitOrderDatabase) *ContractEventsProcessor {
	SetContractFilesLocation(orderBookTestContractFileLocation, marginAccountTestContractFileLocation, clearingHouseTestContractFileLocation)
	return NewContractEventsProcessor(db)
}

func getABIfromJson(fileLocation string) abi.ABI {
	jsonBytes, _ := ioutil.ReadFile(fileLocation)
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

func getOrder(ammIndex *big.Int, traderAddress common.Address, baseAssetQuantity *big.Int, price *big.Int, salt *big.Int) Order {
	return Order{
		AmmIndex:          ammIndex,
		Trader:            traderAddress,
		BaseAssetQuantity: baseAssetQuantity,
		Price:             price,
		Salt:              salt,
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

func getOrderId(trader common.Address, salt *big.Int) string {
	return trader.String() + salt.String()
}
