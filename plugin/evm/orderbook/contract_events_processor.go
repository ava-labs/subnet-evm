package orderbook

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook/abis"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type ContractEventsProcessor struct {
	orderBookABI     abi.ABI
	iocOrderBookABI  abi.ABI
	marginAccountABI abi.ABI
	clearingHouseABI abi.ABI
	database         LimitOrderDatabase
}

func NewContractEventsProcessor(database LimitOrderDatabase) *ContractEventsProcessor {
	orderBookABI, err := abi.FromSolidityJson(string(abis.OrderBookAbi))
	if err != nil {
		panic(err)
	}

	marginAccountABI, err := abi.FromSolidityJson(string(abis.MarginAccountAbi))
	if err != nil {
		panic(err)
	}

	clearingHouseABI, err := abi.FromSolidityJson(string(abis.ClearingHouseAbi))
	if err != nil {
		panic(err)
	}

	iocOrderBookABI, err := abi.FromSolidityJson(string(abis.IOCOrderBookAbi))
	if err != nil {
		panic(err)
	}

	return &ContractEventsProcessor{
		orderBookABI:     orderBookABI,
		marginAccountABI: marginAccountABI,
		clearingHouseABI: clearingHouseABI,
		iocOrderBookABI:  iocOrderBookABI,
		database:         database,
	}
}

func (cep *ContractEventsProcessor) ProcessEvents(logs []*types.Log) {
	var (
		deletedLogs []*types.Log
		rebirthLogs []*types.Log
	)
	for i := 0; i < len(logs); i++ {
		log := logs[i]
		if log.Removed {
			deletedLogs = append(deletedLogs, log)
		} else {
			rebirthLogs = append(rebirthLogs, log)
		}
	}

	// deletedLogs are in descending order by (blockNumber, LogIndex)
	// rebirthLogs should be in ascending order by (blockNumber, LogIndex)
	sort.Slice(deletedLogs, func(i, j int) bool {
		if deletedLogs[i].BlockNumber == deletedLogs[j].BlockNumber {
			return deletedLogs[i].Index > deletedLogs[j].Index
		}
		return deletedLogs[i].BlockNumber > deletedLogs[j].BlockNumber
	})

	sort.Slice(rebirthLogs, func(i, j int) bool {
		if rebirthLogs[i].BlockNumber == rebirthLogs[j].BlockNumber {
			return rebirthLogs[i].Index < rebirthLogs[j].Index
		}
		return rebirthLogs[i].BlockNumber < rebirthLogs[j].BlockNumber
	})

	logs = append(deletedLogs, rebirthLogs...)
	for _, event := range logs {
		switch event.Address {
		case OrderBookContractAddress:
			cep.handleOrderBookEvent(event)
		case IOCOrderBookContractAddress:
			cep.handleIOCOrderBookEvent(event)
		}
	}
}

func (cep *ContractEventsProcessor) ProcessAcceptedEvents(logs []*types.Log, inBootstrap bool) {
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].Index < logs[j].Index
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})

	for _, event := range logs {
		switch event.Address {
		case MarginAccountContractAddress:
			cep.handleMarginAccountEvent(event)
		case ClearingHouseContractAddress:
			cep.handleClearingHouseEvent(event)
		}
	}
	if !inBootstrap {
		// events are applied in sequence during bootstrap also, those shouldn't be updated in metrics as they are already counted
		go cep.updateMetrics(logs)
	}
}

func (cep *ContractEventsProcessor) handleOrderBookEvent(event *types.Log) {
	removed := event.Removed
	args := map[string]interface{}{}
	switch event.Topics[0] {
	// event OrderMatched(address indexed trader, bytes32 indexed orderHash, uint256 fillAmount, uint price, uint openInterestNotional, uint timestamp, bool isLiquidation);
	case cep.orderBookABI.Events["OrderMatched"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderMatched", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderMatched", "err", err)
			return
		}

		trader := getAddressFromTopicHash(event.Topics[1])
		orderId := event.Topics[2]
		fillAmount := args["fillAmount"].(*big.Int)
		if !removed {
			log.Info("OrderMatched", "orderId", orderId.String(), "trader", trader.String(), "args", args, "number", event.BlockNumber)
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount, orderId, event.BlockNumber)
		} else {
			fillAmount.Neg(fillAmount)
			log.Info("OrderMatched removed", "orderId", orderId.String(), "trader", trader.String(), "args", args, "number", event.BlockNumber)
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount, orderId, event.BlockNumber)
		}
	case cep.orderBookABI.Events["OrderMatchingError"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderMatchingError", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderMatchingError", "err", err)
			return
		}
		orderId := event.Topics[1]
		if !removed {
			log.Info("OrderMatchingError", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
			if err := cep.database.SetOrderStatus(orderId, Execution_Failed, args["err"].(string), event.BlockNumber); err != nil {
				log.Error("error in SetOrderStatus", "method", "OrderMatchingError", "err", err)
				return
			}
		} else {
			log.Info("OrderMatchingError removed", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
			if err := cep.database.RevertLastStatus(orderId); err != nil {
				log.Error("error in SetOrderStatus", "method", "OrderMatchingError", "removed", true, "err", err)
				return
			}
		}

	// event OrderAccepted(address indexed trader, bytes32 indexed orderHash, OrderV2 order, uint timestamp);
	case cep.orderBookABI.Events["OrderAccepted"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderAccepted", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderAccepted", "err", err)
			return
		}

		orderId := event.Topics[2]
		if !removed {
			timestamp := args["timestamp"].(*big.Int)
			order := LimitOrder{}
			order.DecodeFromRawOrder(args["order"])

			limitOrder := Order{
				Id:                      orderId,
				Market:                  Market(order.AmmIndex.Int64()),
				PositionType:            getPositionTypeBasedOnBaseAssetQuantity(order.BaseAssetQuantity),
				Trader:                  getAddressFromTopicHash(event.Topics[1]),
				BaseAssetQuantity:       order.BaseAssetQuantity,
				FilledBaseAssetQuantity: big.NewInt(0),
				Price:                   order.Price,
				RawOrder:                &order,
				Salt:                    order.Salt,
				ReduceOnly:              order.ReduceOnly,
				BlockNumber:             big.NewInt(int64(event.BlockNumber)),
				OrderType:               LimitOrderType,
			}
			log.Info("LimitOrder/OrderAccepted", "order", limitOrder, "timestamp", timestamp)
			cep.database.Add(&limitOrder)
		} else {
			log.Info("LimitOrder/OrderAccepted removed", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
			cep.database.Delete(orderId)
		}

	// event OrderRejected(address indexed trader, bytes32 indexed orderHash, OrderV2 order, uint timestamp, string err);
	case cep.orderBookABI.Events["OrderRejected"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderRejected", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderRejected", "err", err)
			return
		}

		orderId := event.Topics[2]
		order := args["order"]
		if !removed {
			log.Info("LimitOrder/OrderRejected", "args", args, "orderId", orderId.String(), "number", event.BlockNumber, "order", order)
		} else {
			log.Info("LimitOrder/OrderRejected removed", "args", args, "orderId", orderId.String(), "number", event.BlockNumber, "order", order)
		}

	// event OrderCancelAccepted(address indexed trader, bytes32 indexed orderHash, uint timestamp);
	case cep.orderBookABI.Events["OrderCancelAccepted"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderCancelAccepted", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderCancelAccepted", "err", err)
			return
		}

		orderId := event.Topics[2]
		if !removed {
			timestamp := args["timestamp"].(*big.Int)
			log.Info("LimitOrder/OrderCancelAccepted", "args", args, "orderId", orderId.String(), "number", event.BlockNumber, "timestamp", timestamp)
			if err := cep.database.SetOrderStatus(orderId, Cancelled, "", event.BlockNumber); err != nil {
				log.Error("error in SetOrderStatus", "method", "OrderCancelAccepted", "err", err)
				return
			}
		} else {
			log.Info("LimitOrder/OrderCancelAccepted removed", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
			if err := cep.database.RevertLastStatus(orderId); err != nil {
				log.Error("error in SetOrderStatus", "method", "OrderCancelAccepted", "removed", true, "err", err)
				return
			}
		}

	// event OrderCancelRejected(address indexed trader, bytes32 indexed orderHash, uint timestamp, string err);
	case cep.orderBookABI.Events["OrderCancelRejected"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderCancelRejected", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderCancelRejected", "err", err)
			return
		}

		orderId := event.Topics[2]
		if !removed {
			log.Info("LimitOrder/OrderCancelRejected", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
		} else {
			log.Info("LimitOrder/OrderCancelRejected removed", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
		}
	}
}

func (cep *ContractEventsProcessor) handleIOCOrderBookEvent(event *types.Log) {
	removed := event.Removed
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.iocOrderBookABI.Events["OrderPlaced"].ID:
		err := cep.iocOrderBookABI.UnpackIntoMap(args, "OrderPlaced", event.Data)
		if err != nil {
			log.Error("error in iocOrderBookABI.UnpackIntoMap", "method", "OrderPlaced", "err", err)
			return
		}
		orderId := event.Topics[2]
		if !removed {
			order := IOCOrder{}
			order.DecodeFromRawOrder(args["order"])
			limitOrder := Order{
				Id:                      orderId,
				Market:                  Market(order.AmmIndex.Int64()),
				PositionType:            getPositionTypeBasedOnBaseAssetQuantity(order.BaseAssetQuantity),
				Trader:                  getAddressFromTopicHash(event.Topics[1]),
				BaseAssetQuantity:       order.BaseAssetQuantity,
				FilledBaseAssetQuantity: big.NewInt(0),
				Price:                   order.Price,
				RawOrder:                &order,
				Salt:                    order.Salt,
				ReduceOnly:              order.ReduceOnly,
				BlockNumber:             big.NewInt(int64(event.BlockNumber)),
				OrderType:               IOCOrderType,
			}
			log.Info("IOCOrder/OrderPlaced", "order", limitOrder, "number", event.BlockNumber)
			cep.database.Add(&limitOrder)
		} else {
			log.Info("IOCOrder/OrderPlaced removed", "orderId", orderId.String(), "block", event.BlockHash.String(), "number", event.BlockNumber)
			cep.database.Delete(orderId)
		}
	}
}

func (cep *ContractEventsProcessor) handleMarginAccountEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.marginAccountABI.Events["MarginAdded"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginAdded", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginAdded", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		collateral := event.Topics[2].Big().Int64()
		amount := args["amount"].(*big.Int)
		log.Info("MarginAdded", "trader", trader, "collateral", collateral, "amount", amount.Uint64(), "number", event.BlockNumber)
		cep.database.UpdateMargin(trader, Collateral(collateral), amount)
	case cep.marginAccountABI.Events["MarginRemoved"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginRemoved", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginRemoved", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		collateral := event.Topics[2].Big().Int64()
		amount := args["amount"].(*big.Int)
		log.Info("MarginRemoved", "trader", trader, "collateral", collateral, "amount", amount.Uint64(), "number", event.BlockNumber)
		cep.database.UpdateMargin(trader, Collateral(collateral), big.NewInt(0).Neg(amount))
	case cep.marginAccountABI.Events["MarginReserved"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginReserved", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginReserved", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		amount := args["amount"].(*big.Int)
		log.Info("MarginReserved", "trader", trader, "amount", amount.Uint64(), "number", event.BlockNumber)
		cep.database.UpdateReservedMargin(trader, amount)
	case cep.marginAccountABI.Events["MarginReleased"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginReleased", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginReleased", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		amount := args["amount"].(*big.Int)
		log.Info("MarginReleased", "trader", trader, "amount", amount.Uint64(), "number", event.BlockNumber)
		cep.database.UpdateReservedMargin(trader, big.NewInt(0).Neg(amount))
	case cep.marginAccountABI.Events["PnLRealized"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "PnLRealized", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "PnLRealized", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		realisedPnL := args["realizedPnl"].(*big.Int)
		log.Info("PnLRealized", "trader", trader, "amount", realisedPnL.Uint64(), "number", event.BlockNumber)
		cep.database.UpdateMargin(trader, HUSD, realisedPnL)
	}
}

func (cep *ContractEventsProcessor) handleClearingHouseEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.clearingHouseABI.Events["FundingRateUpdated"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "FundingRateUpdated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingRateUpdated", "err", err)
			return
		}
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		nextFundingTime := args["nextFundingTime"].(*big.Int)
		market := Market(int(event.Topics[1].Big().Int64()))
		log.Info("FundingRateUpdated", "args", args, "cumulativePremiumFraction", cumulativePremiumFraction, "market", market)
		cep.database.UpdateUnrealisedFunding(market, cumulativePremiumFraction)
		cep.database.UpdateNextFundingTime(nextFundingTime.Uint64())

	case cep.clearingHouseABI.Events["FundingPaid"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "FundingPaid", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingPaid", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		market := Market(int(event.Topics[2].Big().Int64()))
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		log.Info("FundingPaid", "trader", trader, "market", market, "cumulativePremiumFraction", cumulativePremiumFraction)
		cep.database.ResetUnrealisedFunding(market, trader, cumulativePremiumFraction)

	case cep.clearingHouseABI.Events["PositionModified"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PositionModified", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionModified", "err", err)
			return
		}

		trader := getAddressFromTopicHash(event.Topics[1])
		market := Market(int(event.Topics[2].Big().Int64()))
		lastPrice := args["price"].(*big.Int)
		cep.database.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		log.Info("PositionModified", "trader", trader, "market", market, "args", args)
		cep.database.UpdatePosition(trader, market, size, openNotional, false, event.BlockNumber)
	case cep.clearingHouseABI.Events["PositionLiquidated"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PositionLiquidated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionLiquidated", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])

		market := Market(int(event.Topics[2].Big().Int64()))
		lastPrice := args["price"].(*big.Int)
		cep.database.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		log.Info("PositionLiquidated", "market", market, "trader", trader, "args", args)
		cep.database.UpdatePosition(trader, market, size, openNotional, true, event.BlockNumber)

	// event NotifyNextPISample(uint nextSampleTime);
	case cep.clearingHouseABI.Events["NotifyNextPISample"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "NotifyNextPISample", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "NotifyNextPISample", "err", err)
			return
		}
		nextSampleTime := args["nextSampleTime"].(*big.Int)
		log.Info("NotifyNextPISample", "nextSampleTime", nextSampleTime)
		cep.database.UpdateNextSamplePITime(nextSampleTime.Uint64())

	case cep.clearingHouseABI.Events["PISampledUpdated"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PISampledUpdated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PISampledUpdated", "err", err)
			return
		}
		log.Info("PISampledUpdated", "args", args)

	case cep.clearingHouseABI.Events["PISampleSkipped"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PISampleSkipped", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PISampleSkipped", "err", err)
			return
		}
		log.Info("PISampleSkipped", "args", args)
	}
}

type TraderEvent struct {
	Trader          common.Address
	OrderId         common.Hash
	OrderType       string
	Removed         bool
	EventName       string
	Args            map[string]interface{}
	BlockNumber     *big.Int
	BlockStatus     BlockConfirmationLevel
	Timestamp       *big.Int
	TransactionHash common.Hash
}

type MarketFeedEvent struct {
	Trader          common.Address
	Market          Market
	Size            float64
	Price           float64
	Removed         bool
	EventName       string
	BlockNumber     *big.Int
	BlockStatus     BlockConfirmationLevel
	Timestamp       *big.Int
	TransactionHash common.Hash
}

type BlockConfirmationLevel string

const (
	ConfirmationLevelHead     BlockConfirmationLevel = "head"
	ConfirmationLevelAccepted BlockConfirmationLevel = "accepted"
)

func (cep *ContractEventsProcessor) PushToTraderFeed(events []*types.Log, blockStatus BlockConfirmationLevel) {
	for _, event := range events {
		removed := event.Removed
		args := map[string]interface{}{}
		eventName := ""
		var orderId common.Hash
		var orderType string
		var trader common.Address
		txHash := event.TxHash
		switch event.Address {
		case OrderBookContractAddress:
			orderType = "limit"
			switch event.Topics[0] {
			case cep.orderBookABI.Events["OrderAccepted"].ID:
				err := cep.orderBookABI.UnpackIntoMap(args, "OrderAccepted", event.Data)
				if err != nil {
					log.Error("error in orderBookABI.UnpackIntoMap", "method", "OrderAccepted", "err", err)
					continue
				}
				eventName = "OrderAccepted"
				order := LimitOrder{}
				order.DecodeFromRawOrder(args["order"])
				args["order"] = order.Map()
				orderId = event.Topics[2]
				trader = getAddressFromTopicHash(event.Topics[1])

			case cep.orderBookABI.Events["OrderRejected"].ID:
				err := cep.orderBookABI.UnpackIntoMap(args, "OrderRejected", event.Data)
				if err != nil {
					log.Error("error in orderBookABI.UnpackIntoMap", "method", "OrderRejected", "err", err)
					continue
				}
				eventName = "OrderRejected"
				order := LimitOrder{}
				order.DecodeFromRawOrder(args["order"])
				args["order"] = order.Map()
				orderId = event.Topics[2]
				trader = getAddressFromTopicHash(event.Topics[1])

			case cep.orderBookABI.Events["OrderMatched"].ID:
				err := cep.orderBookABI.UnpackIntoMap(args, "OrderMatched", event.Data)
				if err != nil {
					log.Error("error in orderBookABI.UnpackIntoMap", "method", "OrderMatched", "err", err)
					continue
				}
				eventName = "OrderMatched"
				fillAmount := args["fillAmount"].(*big.Int)
				openInterestNotional := args["openInterestNotional"].(*big.Int)
				price := args["price"].(*big.Int)
				args["fillAmount"] = utils.BigIntToFloat(fillAmount, 18)
				args["openInterestNotional"] = utils.BigIntToFloat(openInterestNotional, 18)
				args["price"] = utils.BigIntToFloat(price, 6)
				orderId = event.Topics[2]
				trader = getAddressFromTopicHash(event.Topics[1])

			case cep.orderBookABI.Events["OrderCancelAccepted"].ID:
				err := cep.orderBookABI.UnpackIntoMap(args, "OrderCancelAccepted", event.Data)
				if err != nil {
					log.Error("error in orderBookABI.UnpackIntoMap", "method", "OrderCancelAccepted", "err", err)
					continue
				}
				eventName = "OrderCancelAccepted"
				orderId = event.Topics[2]
				trader = getAddressFromTopicHash(event.Topics[1])

			case cep.orderBookABI.Events["OrderCancelRejected"].ID:
				err := cep.orderBookABI.UnpackIntoMap(args, "OrderCancelRejected", event.Data)
				if err != nil {
					log.Error("error in orderBookABI.UnpackIntoMap", "method", "OrderCancelRejected", "err", err)
					continue
				}
				eventName = "OrderCancelRejected"
				orderId = event.Topics[2]
				trader = getAddressFromTopicHash(event.Topics[1])

			default:
				continue
			}

		case IOCOrderBookContractAddress:
			orderType = "ioc"
			switch event.Topics[0] {
			case cep.iocOrderBookABI.Events["OrderPlaced"].ID:
				err := cep.iocOrderBookABI.UnpackIntoMap(args, "OrderPlaced", event.Data)
				if err != nil {
					log.Error("error in iocOrderBookABI.UnpackIntoMap", "method", "OrderPlaced", "err", err)
					continue
				}
				eventName = "OrderPlaced"
				order := IOCOrder{}
				order.DecodeFromRawOrder(args["order"])
				args["order"] = order.Map()
				orderId = event.Topics[2]
				trader = getAddressFromTopicHash(event.Topics[1])
			}
		default:
			continue
		}

		timestamp := args["timestamp"]
		timestampInt, _ := timestamp.(*big.Int)
		traderEvent := TraderEvent{
			Trader:          trader,
			Removed:         removed,
			EventName:       eventName,
			Args:            args,
			BlockNumber:     big.NewInt(int64(event.BlockNumber)),
			BlockStatus:     blockStatus,
			OrderId:         orderId,
			OrderType:       orderType,
			Timestamp:       timestampInt,
			TransactionHash: txHash,
		}

		traderFeed.Send(traderEvent)
	}
}

func (cep *ContractEventsProcessor) PushToMarketFeed(events []*types.Log, blockStatus BlockConfirmationLevel) {
	for _, event := range events {
		args := map[string]interface{}{}
		switch event.Topics[0] {
		case cep.clearingHouseABI.Events["PositionModified"].ID:
			err := cep.clearingHouseABI.UnpackIntoMap(args, "PositionModified", event.Data)
			if err != nil {
				log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionModified", "err", err)
				return
			}

			trader := getAddressFromTopicHash(event.Topics[1])
			market := Market(int(event.Topics[2].Big().Int64()))
			price := args["price"].(*big.Int)

			size := args["baseAsset"].(*big.Int)

			timestamp := args["timestamp"]
			timestampInt, _ := timestamp.(*big.Int)
			marketFeedEvent := MarketFeedEvent{
				Trader:          trader,
				Market:          market,
				Size:            utils.BigIntToFloat(size, 18),
				Price:           utils.BigIntToFloat(price, 6),
				Removed:         event.Removed,
				EventName:       "PositionModified",
				BlockNumber:     big.NewInt(int64(event.BlockNumber)),
				BlockStatus:     blockStatus,
				Timestamp:       timestampInt,
				TransactionHash: event.TxHash,
			}
			marketFeed.Send(marketFeedEvent)
		}
	}
}

func (cep *ContractEventsProcessor) updateMetrics(logs []*types.Log) {
	var orderPlacedCount int64 = 0
	var orderCancelledCount int64 = 0
	for _, event := range logs {
		var contractABI abi.ABI
		switch event.Address {
		case OrderBookContractAddress:
			contractABI = cep.orderBookABI
		case MarginAccountContractAddress:
			contractABI = cep.marginAccountABI
		case ClearingHouseContractAddress:
			contractABI = cep.clearingHouseABI
		}

		event_, err := contractABI.EventByID(event.Topics[0])
		if err != nil {
			continue
		}

		metricName := fmt.Sprintf("%s/%s", "events", event_.Name)

		if !event.Removed {
			metrics.GetOrRegisterCounter(metricName, nil).Inc(1)
		} else {
			metrics.GetOrRegisterCounter(metricName, nil).Dec(1)
		}

		switch event_.Name {
		// both LimitOrder and IOCOrder's respective events
		case "OrderPlaced", "OrderAccepted":
			orderPlacedCount++
		case "OrderCancelAccepted":
			orderCancelledCount++
		}
	}

	ordersPlacedPerBlock.Update(orderPlacedCount)
	ordersCancelledPerBlock.Update(orderCancelledCount)
}

func getAddressFromTopicHash(topicHash common.Hash) common.Address {
	return common.BytesToAddress(topicHash.Bytes())
}
