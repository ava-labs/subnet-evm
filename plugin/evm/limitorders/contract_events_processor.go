package limitorders

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type ContractEventsProcessor struct {
	orderBookABI     abi.ABI
	marginAccountABI abi.ABI
	clearingHouseABI abi.ABI
	database         LimitOrderDatabase
}

func NewContractEventsProcessor(database LimitOrderDatabase) *ContractEventsProcessor {
	jsonBytes, _ := ioutil.ReadFile(orderBookContractFileLocation)
	orderBookABI, err := abi.FromSolidityJson(string(jsonBytes))
	if err != nil {
		panic(err)
	}

	jsonBytes, _ = ioutil.ReadFile(marginAccountContractFileLocation)
	marginAccountABI, err := abi.FromSolidityJson(string(jsonBytes))
	if err != nil {
		panic(err)
	}

	jsonBytes, _ = ioutil.ReadFile(clearingHouseContractFileLocation)
	clearingHouseABI, err := abi.FromSolidityJson(string(jsonBytes))
	if err != nil {
		panic(err)
	}
	return &ContractEventsProcessor{
		orderBookABI:     orderBookABI,
		marginAccountABI: marginAccountABI,
		clearingHouseABI: clearingHouseABI,
		database:         database,
	}
}

func (cep *ContractEventsProcessor) HandleOrderBookEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.orderBookABI.Events["OrderPlaced"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderPlaced", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderPlaced", "err", err)
			return
		}
		log.Info("HandleOrderBookEvent", "orderplaced args", args)
		order := getOrderFromRawOrder(args["order"])

		cep.database.Add(&LimitOrder{
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
	case cep.orderBookABI.Events["OrderCancelled"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderCancelled", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderCancelled", "err", err)
			return
		}
		log.Info("HandleOrderBookEvent", "OrderCancelled args", args)
		signature := args["signature"].([]byte)

		cep.database.Delete(signature)
	case cep.orderBookABI.Events["OrdersMatched"].ID:
		log.Info("OrdersMatched event")
		err := cep.orderBookABI.UnpackIntoMap(args, "OrdersMatched", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrdersMatched", "err", err)
			return
		}
		log.Info("HandleOrderBookEvent", "OrdersMatched args", args)
		fmt.Println("xxxxx")
		signatures := args["signatures"].([2][]byte)
		fmt.Println("yyyy")
		fillAmount := args["fillAmount"].(*big.Int)
		cep.database.UpdateFilledBaseAssetQuantity(fillAmount, signatures[0])
		cep.database.UpdateFilledBaseAssetQuantity(fillAmount, signatures[1])
	case cep.orderBookABI.Events["LiquidationOrderMatched"].ID:
		log.Info("LiquidationOrderMatched event")
		err := cep.orderBookABI.UnpackIntoMap(args, "LiquidationOrderMatched", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "LiquidationOrderMatched", "err", err)
			return
		}
		log.Info("HandleOrderBookEvent", "LiquidationOrderMatched args", args)
		signature := args["signature"].([]byte)
		fillAmount := args["fillAmount"].(*big.Int)
		cep.database.UpdateFilledBaseAssetQuantity(fillAmount, signature)
	}
	log.Info("Log found", "log_.Address", event.Address.String(), "log_.BlockNumber", event.BlockNumber, "log_.Index", event.Index, "log_.TxHash", event.TxHash.String())

}

func (cep *ContractEventsProcessor) HandleMarginAccountEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.marginAccountABI.Events["MarginAdded"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginAdded", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginAdded", "err", err)
			return
		}
		collateral := event.Topics[2].Big().Int64()
		cep.database.UpdateMargin(getAddressFromTopicHash(event.Topics[1]), Collateral(collateral), args["amount"].(*big.Int))
	case cep.marginAccountABI.Events["MarginRemoved"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginRemoved", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginRemoved", "err", err)
			return
		}
		collateral := event.Topics[2].Big().Int64()
		cep.database.UpdateMargin(getAddressFromTopicHash(event.Topics[1]), Collateral(collateral), big.NewInt(0).Neg(args["amount"].(*big.Int)))
	case cep.marginAccountABI.Events["PnLRealized"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "PnLRealized", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "PnLRealized", "err", err)
			return
		}
		realisedPnL := args["realizedPnl"].(*big.Int)

		cep.database.UpdateMargin(getAddressFromTopicHash(event.Topics[1]), HUSD, realisedPnL)
	}
	log.Info("Log found", "log_.Address", event.Address.String(), "log_.BlockNumber", event.BlockNumber, "log_.Index", event.Index, "log_.TxHash", event.TxHash.String())
}

func (cep *ContractEventsProcessor) HandleClearingHouseEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.clearingHouseABI.Events["FundingRateUpdated"].ID:
		log.Info("FundingRateUpdated event")
		err := cep.clearingHouseABI.UnpackIntoMap(args, "FundingRateUpdated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingRateUpdated", "err", err)
			return
		}
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		nextFundingTime := args["nextFundingTime"].(*big.Int)
		market := Market(int(event.Topics[1].Big().Int64()))
		cep.database.UpdateUnrealisedFunding(Market(market), cumulativePremiumFraction)
		cep.database.UpdateNextFundingTime(nextFundingTime.Uint64())

	case cep.clearingHouseABI.Events["FundingPaid"].ID:
		log.Info("FundingPaid event")
		err := cep.clearingHouseABI.UnpackIntoMap(args, "FundingPaid", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingPaid", "err", err)
			return
		}
		market := Market(int(event.Topics[2].Big().Int64()))
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		cep.database.ResetUnrealisedFunding(Market(market), getAddressFromTopicHash(event.Topics[1]), cumulativePremiumFraction)

	// both PositionModified and PositionLiquidated have the exact same signature
	case cep.clearingHouseABI.Events["PositionModified"].ID:
		log.Info("PositionModified event")
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PositionModified", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionModified", "err", err)
			return
		}

		market := Market(int(event.Topics[2].Big().Int64()))
		baseAsset := args["baseAsset"].(*big.Int)
		quoteAsset := args["quoteAsset"].(*big.Int)
		lastPrice := big.NewInt(0).Div(big.NewInt(0).Mul(quoteAsset, big.NewInt(1e18)), baseAsset)
		cep.database.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		cep.database.UpdatePosition(getAddressFromTopicHash(event.Topics[1]), market, size, openNotional, false)
	case cep.clearingHouseABI.Events["PositionLiquidated"].ID:
		log.Info("PositionLiquidated event")
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PositionLiquidated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionLiquidated", "err", err)
			return
		}

		market := Market(int(event.Topics[2].Big().Int64()))
		baseAsset := args["baseAsset"].(*big.Int)
		quoteAsset := args["quoteAsset"].(*big.Int)
		lastPrice := big.NewInt(0).Div(big.NewInt(0).Mul(quoteAsset, big.NewInt(1e18)), baseAsset)
		cep.database.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		cep.database.UpdatePosition(getAddressFromTopicHash(event.Topics[1]), market, size, openNotional, true)
	}
}
