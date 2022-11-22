package abi

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"
)

const ammJson = `[{"inputs":[{"internalType":"address","name":"_clearingHouse","type":"address"},{"internalType":"uint256","name":"_unbondRoundOff","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"trader","type":"address"},{"indexed":false,"internalType":"int256","name":"takerFundingPayment","type":"int256"},{"indexed":false,"internalType":"int256","name":"makerFundingPayment","type":"int256"}],"name":"FundingPaid","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"int256","name":"premiumFraction","type":"int256"},{"indexed":false,"internalType":"uint256","name":"underlyingPrice","type":"uint256"},{"indexed":false,"internalType":"int256","name":"cumulativePremiumFraction","type":"int256"},{"indexed":false,"internalType":"int256","name":"cumulativePremiumPerDtoken","type":"int256"},{"indexed":false,"internalType":"int256","name":"posAccumulator","type":"int256"},{"indexed":false,"internalType":"uint256","name":"nextFundingTime","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"FundingRateUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"trader","type":"address"},{"indexed":false,"internalType":"uint256","name":"dToken","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"baseAsset","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"quoteAsset","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"LiquidityAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"trader","type":"address"},{"indexed":false,"internalType":"uint256","name":"dToken","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"baseAsset","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"quoteAsset","type":"uint256"},{"indexed":false,"internalType":"int256","name":"realizedPnl","type":"int256"},{"indexed":false,"internalType":"bool","name":"isLiquidation","type":"bool"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"LiquidityRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"trader","type":"address"},{"components":[{"internalType":"uint256","name":"vUSD","type":"uint256"},{"internalType":"uint256","name":"vAsset","type":"uint256"},{"internalType":"uint256","name":"dToken","type":"uint256"},{"internalType":"int256","name":"pos","type":"int256"},{"internalType":"int256","name":"posAccumulator","type":"int256"},{"internalType":"int256","name":"lastPremiumFraction","type":"int256"},{"internalType":"int256","name":"lastPremiumPerDtoken","type":"int256"},{"internalType":"uint256","name":"unbondTime","type":"uint256"},{"internalType":"uint256","name":"unbondAmount","type":"uint256"},{"internalType":"uint256","name":"ignition","type":"uint256"}],"indexed":false,"internalType":"struct IAMM.Maker","name":"maker","type":"tuple"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"MakerPositionChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"trader","type":"address"},{"indexed":false,"internalType":"int256","name":"size","type":"int256"},{"indexed":false,"internalType":"uint256","name":"openNotional","type":"uint256"},{"indexed":false,"internalType":"int256","name":"realizedPnl","type":"int256"}],"name":"PositionChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"lastPrice","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"openInterestNotional","type":"uint256"}],"name":"Swap","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"trader","type":"address"},{"indexed":false,"internalType":"uint256","name":"unbondAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"unbondTime","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"Unbonded","type":"event"},{"inputs":[{"internalType":"address","name":"maker","type":"address"},{"internalType":"uint256","name":"baseAssetQuantity","type":"uint256"},{"internalType":"uint256","name":"minDToken","type":"uint256"}],"name":"addLiquidity","outputs":[{"internalType":"uint256","name":"dToken","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"ammState","outputs":[{"internalType":"enum IAMM.AMMState","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_oracle","type":"address"}],"name":"changeOracle","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"clearingHouse","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"maker","type":"address"},{"internalType":"uint256","name":"quoteAsset","type":"uint256"}],"name":"commitLiquidity","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"cumulativePremiumFraction","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"cumulativePremiumPerDtoken","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"maker","type":"address"}],"name":"forceRemoveLiquidity","outputs":[{"internalType":"int256","name":"realizedPnl","type":"int256"},{"internalType":"uint256","name":"makerOpenNotional","type":"uint256"},{"internalType":"int256","name":"makerPosition","type":"int256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"fundingBufferPeriod","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"fundingPeriod","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"int256","name":"baseAssetQuantity","type":"int256"}],"name":"getCloseQuote","outputs":[{"internalType":"uint256","name":"quoteAssetQuantity","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"vUSD","type":"uint256"}],"name":"getIgnitionShare","outputs":[{"internalType":"uint256","name":"vAsset","type":"uint256"},{"internalType":"uint256","name":"dToken","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"trader","type":"address"}],"name":"getNotionalPositionAndUnrealizedPnl","outputs":[{"internalType":"uint256","name":"notionalPosition","type":"uint256"},{"internalType":"int256","name":"unrealizedPnl","type":"int256"},{"internalType":"int256","name":"size","type":"int256"},{"internalType":"uint256","name":"openNotional","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"int256","name":"positionSize","type":"int256"},{"internalType":"uint256","name":"newNotionalPosition","type":"uint256"},{"internalType":"int256","name":"unrealizedPnl","type":"int256"},{"internalType":"int256","name":"baseAssetQuantity","type":"int256"}],"name":"getOpenNotionalWhileReducingPosition","outputs":[{"internalType":"uint256","name":"remainOpenNotional","type":"uint256"},{"internalType":"int256","name":"realizedPnl","type":"int256"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"trader","type":"address"},{"internalType":"int256","name":"margin","type":"int256"},{"internalType":"enum IClearingHouse.Mode","name":"mode","type":"uint8"}],"name":"getOracleBasedPnl","outputs":[{"internalType":"uint256","name":"notionalPosition","type":"uint256"},{"internalType":"int256","name":"unrealizedPnl","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"trader","type":"address"}],"name":"getPendingFundingPayment","outputs":[{"internalType":"int256","name":"takerFundingPayment","type":"int256"},{"internalType":"int256","name":"makerFundingPayment","type":"int256"},{"internalType":"int256","name":"latestCumulativePremiumFraction","type":"int256"},{"internalType":"int256","name":"latestPremiumPerDtoken","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getSnapshotLen","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"trader","type":"address"}],"name":"getTakerNotionalPositionAndUnrealizedPnl","outputs":[{"internalType":"uint256","name":"takerNotionalPosition","type":"uint256"},{"internalType":"int256","name":"unrealizedPnl","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_intervalInSeconds","type":"uint256"}],"name":"getTwapPrice","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_intervalInSeconds","type":"uint256"}],"name":"getUnderlyingTwapPrice","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"governance","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"ignition","outputs":[{"internalType":"uint256","name":"quoteAsset","type":"uint256"},{"internalType":"uint256","name":"baseAsset","type":"uint256"},{"internalType":"uint256","name":"dToken","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"},{"internalType":"address","name":"_underlyingAsset","type":"address"},{"internalType":"address","name":"_oracle","type":"address"},{"internalType":"uint256","name":"_minSizeRequirement","type":"uint256"},{"internalType":"address","name":"_vamm","type":"address"},{"internalType":"address","name":"_governance","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"isOverSpreadLimit","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"lastPrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"liftOff","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"trader","type":"address"}],"name":"liquidatePosition","outputs":[{"internalType":"int256","name":"realizedPnl","type":"int256"},{"internalType":"int256","name":"baseAsset","type":"int256"},{"internalType":"uint256","name":"quoteAsset","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"longOpenInterestNotional","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"maker","type":"address"}],"name":"makers","outputs":[{"components":[{"internalType":"uint256","name":"vUSD","type":"uint256"},{"internalType":"uint256","name":"vAsset","type":"uint256"},{"internalType":"uint256","name":"dToken","type":"uint256"},{"internalType":"int256","name":"pos","type":"int256"},{"internalType":"int256","name":"posAccumulator","type":"int256"},{"internalType":"int256","name":"lastPremiumFraction","type":"int256"},{"internalType":"int256","name":"lastPremiumPerDtoken","type":"int256"},{"internalType":"uint256","name":"unbondTime","type":"uint256"},{"internalType":"uint256","name":"unbondAmount","type":"uint256"},{"internalType":"uint256","name":"ignition","type":"uint256"}],"internalType":"struct IAMM.Maker","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxFundingRate","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxLiquidationPriceSpread","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxLiquidationRatio","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxOracleSpreadRatio","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxPriceSpreadPerBlock","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"minSizeRequirement","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"nextFundingTime","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"openInterestNotional","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"trader","type":"address"},{"internalType":"int256","name":"baseAssetQuantity","type":"int256"},{"internalType":"uint256","name":"quoteAssetLimit","type":"uint256"}],"name":"openPosition","outputs":[{"internalType":"int256","name":"realizedPnl","type":"int256"},{"internalType":"uint256","name":"quoteAsset","type":"uint256"},{"internalType":"bool","name":"isPositionIncreased","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"oracle","outputs":[{"internalType":"contract IOracle","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"posAccumulator","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"positions","outputs":[{"internalType":"int256","name":"size","type":"int256"},{"internalType":"uint256","name":"openNotional","type":"uint256"},{"internalType":"int256","name":"lastPremiumFraction","type":"int256"},{"internalType":"uint256","name":"liquidationThreshold","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"putAmmInIgnition","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"maker","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"minQuote","type":"uint256"},{"internalType":"uint256","name":"minBase","type":"uint256"}],"name":"removeLiquidity","outputs":[{"internalType":"int256","name":"realizedPnl","type":"int256"},{"internalType":"uint256","name":"makerOpenNotional","type":"uint256"},{"internalType":"int256","name":"makerPosition","type":"int256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"reserveSnapshots","outputs":[{"internalType":"uint256","name":"lastPrice","type":"uint256"},{"internalType":"uint256","name":"timestamp","type":"uint256"},{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bool","name":"isLiquidation","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_fundingBufferPeriod","type":"uint256"}],"name":"setFundingBufferPeriod","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_governance","type":"address"}],"name":"setGovernace","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_maxLiquidationRatio","type":"uint256"},{"internalType":"uint256","name":"_maxLiquidationPriceSpread","type":"uint256"}],"name":"setLiquidationParams","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"int256","name":"_maxFundingRate","type":"int256"}],"name":"setMaxFundingRate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_minSizeRequirement","type":"uint256"}],"name":"setMinSizeRequirement","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_maxOracleSpreadRatio","type":"uint256"},{"internalType":"uint256","name":"_maxPriceSpreadPerBlock","type":"uint256"}],"name":"setPriceSpreadParams","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_unbondPeriod","type":"uint256"}],"name":"setUnbondPeriod","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"settleFunding","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"shortOpenInterestNotional","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"spotPriceTwapInterval","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"dToken","type":"uint256"}],"name":"unbondLiquidity","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unbondPeriod","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unbondRoundOff","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"underlyingAsset","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"trader","type":"address"}],"name":"updatePosition","outputs":[{"internalType":"int256","name":"fundingPayment","type":"int256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"vamm","outputs":[{"internalType":"contract IVAMM","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"withdrawPeriod","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

type SolidityJSON struct {
	ContractName string `json:"contractName"`
	SourceName   string `json:"sourceName"`
	Abi          []Abi  `json:"abi"`
}

type Abi struct {
	Inputs          []Input `json:"inputs"`
	StateMutability string  `json:"stateMutability,omitempty"`
	Type            string  `json:"type"`
	Anonymous       bool    `json:"anonymous,omitempty"`
	Name            string  `json:"name,omitempty"`
	Outputs         []Input `json:"outputs,omitempty"`
}

type Input struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Indexed      bool   `json:"indexed"`
}

func getFunctionType(type_ string) FunctionType {
	typeMap := map[string]FunctionType{
		"function":    Function,
		"receive":     Receive,
		"fallback":    Fallback,
		"constructor": Constructor,
	}
	return typeMap[type_]
}

func FromSolidityJson(abiJsonInput string) (ABI, error) {
	solidityJson := SolidityJSON{}
	err := json.Unmarshal([]byte(abiJsonInput), &solidityJson)
	if err != nil {
		log.Error("Error in decoding ABI json", "err", err)
		return ABI{}, err
	}

	var constructor Method
	var receive Method
	var fallback Method

	methods := map[string]Method{}
	events := map[string]Event{}
	errors := map[string]Error{}

	for _, method := range solidityJson.Abi {
		inputs := []Argument{}
		for _, input := range method.Inputs {
			type_, _ := NewType(input.Type, input.InternalType, nil)
			inputs = append(inputs, Argument{
				Name:    input.Name,
				Type:    type_,
				Indexed: input.Indexed,
			})
		}

		if method.Type == "event" {
			abiEvent := NewEvent(method.Name, method.Name, method.Anonymous, inputs)
			events[method.Name] = abiEvent
			continue
		}

		if method.Type == "error" {
			abiError := NewError(method.Name, inputs)
			errors[method.Name] = abiError
			continue
		}

		outputs := []Argument{}
		for _, output := range method.Outputs {
			type_, _ := NewType(output.Type, output.InternalType, nil)
			inputs = append(inputs, Argument{
				Name:    output.Name,
				Type:    type_,
				Indexed: output.Indexed,
			})
		}

		methodType := getFunctionType(method.Type)
		abiMethod := NewMethod(method.Name, method.Name, methodType, method.StateMutability, false, method.StateMutability == "payable", inputs, outputs)

		// don't include the method in the list if it's a constructor
		if methodType == Constructor {
			constructor = abiMethod
			continue
		}
		if methodType == Fallback {
			fallback = abiMethod
			continue
		}
		if methodType == Receive {
			receive = abiMethod
			continue
		}
		methods[method.Name] = abiMethod
	}

	Abi := ABI{
		Constructor: constructor,
		Methods:     methods,
		Events:      events,
		Errors:      errors,

		Receive:  receive,
		Fallback: fallback,
	}

	return Abi, nil
}
