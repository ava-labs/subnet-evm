package limitorders

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func (lotp *limitOrderTxProcessor) HandleOrderBookEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case lotp.orderBookABI.Events["OrderPlaced"].ID:
		err := lotp.orderBookABI.UnpackIntoMap(args, "OrderPlaced", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderPlaced", "err", err)
		}
		log.Info("HandleOrderBookEvent", "orderplaced args", args)
		order := getOrderFromRawOrder(args["order"])

		lotp.memoryDb.Add(&LimitOrder{
			Market:            Market(order.AmmIndex.Int64()),
			PositionType:      getPositionTypeBasedOnBaseAssetQuantity(order.BaseAssetQuantity),
			UserAddress:       getAddressFromTopicHash(event.Topics[1]).String(),
			BaseAssetQuantity: order.BaseAssetQuantity,
			Price:             order.Price,
			Status:            Placed,
			RawOrder:          args["order"],
			Signature:         args["signature"].([]byte),
			BlockNumber:       big.NewInt(int64(event.BlockNumber)),
		})
	case lotp.orderBookABI.Events["OrderCancelled"].ID:
		err := lotp.orderBookABI.UnpackIntoMap(args, "OrderCancelled", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderCancelled", "err", err)
		}
		log.Info("HandleOrderBookEvent", "OrderCancelled args", args)
		signature := args["signature"].([]byte)

		lotp.memoryDb.Delete(signature)
	case lotp.orderBookABI.Events["OrdersMatched"].ID:
		log.Info("OrdersMatched event")
		err := lotp.orderBookABI.UnpackIntoMap(args, "OrdersMatched", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrdersMatched", "err", err)
		}
		log.Info("HandleOrderBookEvent", "OrdersMatched args", args)
		signatures := args["signatures"].([][]byte)
		fillAmount := args["fillAmount"].(*big.Int)
		lotp.memoryDb.UpdateFilledBaseAssetQuantity(fillAmount, signatures[0])
		lotp.memoryDb.UpdateFilledBaseAssetQuantity(fillAmount, signatures[1])
	case lotp.orderBookABI.Events["LiquidationOrderMatched"].ID:
		log.Info("LiquidationOrderMatched event")
		err := lotp.orderBookABI.UnpackIntoMap(args, "LiquidationOrderMatched", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "LiquidationOrderMatched", "err", err)
		}
		log.Info("HandleOrderBookEvent", "LiquidationOrderMatched args", args)
		signature := args["signature"].([]byte)
		fillAmount := args["fillAmount"].(*big.Int)
		lotp.memoryDb.UpdateFilledBaseAssetQuantity(fillAmount, signature)
	}
	log.Info("Log found", "log_.Address", event.Address.String(), "log_.BlockNumber", event.BlockNumber, "log_.Index", event.Index, "log_.TxHash", event.TxHash.String())

}

func (lotp *limitOrderTxProcessor) HandleMarginAccountEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case lotp.marginAccountABI.Events["MarginAdded"].ID:
		err := lotp.marginAccountABI.UnpackIntoMap(args, "MarginAdded", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginAdded", "err", err)
		}
		collateral := event.Topics[2].Big().Int64()
		lotp.memoryDb.UpdateMargin(getAddressFromTopicHash(event.Topics[1]), Collateral(collateral), args["amount"].(*big.Int))
	case lotp.marginAccountABI.Events["MarginRemoved"].ID:
		err := lotp.marginAccountABI.UnpackIntoMap(args, "MarginRemoved", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginRemoved", "err", err)
		}
		collateral := event.Topics[2].Big().Int64()
		lotp.memoryDb.UpdateMargin(getAddressFromTopicHash(event.Topics[1]), Collateral(collateral), big.NewInt(0).Neg(args["amount"].(*big.Int)))
	case lotp.marginAccountABI.Events["PnLRealized"].ID:
		err := lotp.marginAccountABI.UnpackIntoMap(args, "PnLRealized", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "PnLRealized", "err", err)
		}
		realisedPnL := args["realizedPnl"].(*big.Int)

		lotp.memoryDb.UpdateMargin(getAddressFromTopicHash(event.Topics[1]), HUSD, realisedPnL)
	}
	log.Info("Log found", "log_.Address", event.Address.String(), "log_.BlockNumber", event.BlockNumber, "log_.Index", event.Index, "log_.TxHash", event.TxHash.String())
}

func (lotp *limitOrderTxProcessor) HandleClearingHouseEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case lotp.clearingHouseABI.Events["FundingRateUpdated"].ID:
		log.Info("FundingRateUpdated event")
		err := lotp.clearingHouseABI.UnpackIntoMap(args, "FundingRateUpdated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingRateUpdated", "err", err)
		}
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		nextFundingTime := args["nextFundingTime"].(*big.Int)
		market := Market(int(event.Topics[1].Big().Int64()))
		lotp.memoryDb.UpdateUnrealisedFunding(Market(market), cumulativePremiumFraction)
		lotp.memoryDb.UpdateNextFundingTime(nextFundingTime.Uint64())

	case lotp.clearingHouseABI.Events["FundingPaid"].ID:
		log.Info("FundingPaid event")
		err := lotp.clearingHouseABI.UnpackIntoMap(args, "FundingPaid", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingPaid", "err", err)
		}
		market := Market(int(event.Topics[2].Big().Int64()))
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		lotp.memoryDb.ResetUnrealisedFunding(Market(market), getAddressFromTopicHash(event.Topics[1]), cumulativePremiumFraction)

	// both PositionModified and PositionLiquidated have the exact same signature
	case lotp.clearingHouseABI.Events["PositionModified"].ID:
		log.Info("PositionModified event")
		err := lotp.clearingHouseABI.UnpackIntoMap(args, "PositionModified", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionModified", "err", err)
		}

		market := Market(int(event.Topics[2].Big().Int64()))
		baseAsset := args["baseAsset"].(*big.Int)
		quoteAsset := args["quoteAsset"].(*big.Int)
		lastPrice := big.NewInt(0).Div(big.NewInt(0).Mul(quoteAsset, big.NewInt(1e18)), baseAsset)
		lotp.memoryDb.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		lotp.memoryDb.UpdatePosition(getAddressFromTopicHash(event.Topics[1]), market, size, openNotional, false)
	case lotp.clearingHouseABI.Events["PositionLiquidated"].ID:
		log.Info("PositionLiquidated event")
		err := lotp.clearingHouseABI.UnpackIntoMap(args, "PositionLiquidated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionLiquidated", "err", err)
		}

		market := Market(int(event.Topics[2].Big().Int64()))
		baseAsset := args["baseAsset"].(*big.Int)
		quoteAsset := args["quoteAsset"].(*big.Int)
		lastPrice := big.NewInt(0).Div(big.NewInt(0).Mul(quoteAsset, big.NewInt(1e18)), baseAsset)
		lotp.memoryDb.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		lotp.memoryDb.UpdatePosition(getAddressFromTopicHash(event.Topics[1]), market, size, openNotional, true)
	}
}
