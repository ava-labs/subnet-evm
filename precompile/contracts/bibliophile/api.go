package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

type VariablesReadFromClearingHouseSlots struct {
	MaintenanceMargin  *big.Int         `json:"maintenance_margin"`
	MinAllowableMargin *big.Int         `json:"min_allowable_margin"`
	TakerFee           *big.Int         `json:"taker_fee"`
	Amms               []common.Address `json:"amms"`
	ActiveMarketsCount int64            `json:"active_markets_count"`
	NotionalPosition   *big.Int         `json:"notional_position"`
	Margin             *big.Int         `json:"margin"`
	TotalFunding       *big.Int         `json:"total_funding"`
	UnderlyingPrices   []*big.Int       `json:"underlying_prices"`
	PositionSizes      []*big.Int       `json:"position_sizes"`
}

func GetClearingHouseVariables(stateDB contract.StateDB, trader common.Address) VariablesReadFromClearingHouseSlots {
	maintenanceMargin := GetMaintenanceMargin(stateDB)
	minAllowableMargin := GetMinAllowableMargin(stateDB)
	takerFee := GetTakerFee(stateDB)
	amms := GetMarkets(stateDB)
	activeMarketsCount := GetActiveMarketsCount(stateDB)
	notionalPositionAndMargin := getNotionalPositionAndMargin(stateDB, &GetNotionalPositionAndMarginInput{
		Trader:                 trader,
		IncludeFundingPayments: false,
		Mode:                   0,
	}, 0 /* use new algorithm */)
	totalFunding := GetTotalFunding(stateDB, &trader)
	positionSizes := getPosSizes(stateDB, &trader)
	underlyingPrices := GetUnderlyingPrices(stateDB)

	return VariablesReadFromClearingHouseSlots{
		MaintenanceMargin:  maintenanceMargin,
		MinAllowableMargin: minAllowableMargin,
		TakerFee:           takerFee,
		Amms:               amms,
		ActiveMarketsCount: activeMarketsCount,
		NotionalPosition:   notionalPositionAndMargin.NotionalPosition,
		Margin:             notionalPositionAndMargin.Margin,
		TotalFunding:       totalFunding,
		PositionSizes:      positionSizes,
		UnderlyingPrices:   underlyingPrices,
	}
}

type VariablesReadFromMarginAccountSlots struct {
	Margin           *big.Int `json:"margin"`
	NormalizedMargin *big.Int `json:"normalized_margin"`
	ReservedMargin   *big.Int `json:"reserved_margin"`
}

func GetMarginAccountVariables(stateDB contract.StateDB, collateralIdx *big.Int, trader common.Address) VariablesReadFromMarginAccountSlots {
	margin := getMargin(stateDB, collateralIdx, trader)
	normalizedMargin := GetNormalizedMargin(stateDB, trader)
	reservedMargin := getReservedMargin(stateDB, trader)
	return VariablesReadFromMarginAccountSlots{
		Margin:           margin,
		NormalizedMargin: normalizedMargin,
		ReservedMargin:   reservedMargin,
	}
}

type VariablesReadFromAMMSlots struct {
	// positions, cumulativePremiumFraction, maxOracleSpreadRatio, maxLiquidationRatio, minSizeRequirement, oracle, underlyingAsset,
	// maxLiquidationPriceSpread, redStoneAdapter, redStoneFeedId, impactMarginNotional, lastTradePrice, bids, asks, bidsHead, asksHead
	LastPrice                 *big.Int       `json:"last_price"`
	CumulativePremiumFraction *big.Int       `json:"cumulative_premium_fraction"`
	MaxOracleSpreadRatio      *big.Int       `json:"max_oracle_spread_ratio"`
	OracleAddress             common.Address `json:"oracle_address"`
	MaxLiquidationRatio       *big.Int       `json:"max_liquidation_ratio"`
	MinSizeRequirement        *big.Int       `json:"min_size_requirement"`
	UnderlyingAssetAddress    common.Address `json:"underlying_asset_address"`
	UnderlyingPriceForMarket  *big.Int       `json:"underlying_price_for_market"`
	UnderlyingPrice           *big.Int       `json:"underlying_price"`
	MaxLiquidationPriceSpread *big.Int       `json:"max_liquidation_price_spread"`
	RedStoneAdapterAddress    common.Address `json:"red_stone_adapter_address"`
	RedStoneFeedId            common.Hash    `json:"red_stone_feed_id"`
	ImpactMarginNotional      *big.Int       `json:"impact_margin_notional"`
	Position                  Position       `json:"position"`
	BidsHead                  *big.Int       `json:"bids_head"`
	BidsHeadSize              *big.Int       `json:"bids_head_size"`
	AsksHead                  *big.Int       `json:"asks_head"`
	AsksHeadSize              *big.Int       `json:"asks_head_size"`
	UpperBound                *big.Int       `json:"upper_bound"`
	LowerBound                *big.Int       `json:"lower_bound"`
	MinAllowableMargin        *big.Int       `json:"min_allowable_margin"`
	TakerFee                  *big.Int       `json:"taker_fee"`
	TotalMargin               *big.Int       `json:"total_margin"`
	AvailableMargin           *big.Int       `json:"available_margin"`
	ReduceOnlyAmount          *big.Int       `json:"reduce_only_amount"`
	LongOpenOrders            *big.Int       `json:"long_open_orders"`
	ShortOpenOrders           *big.Int       `json:"short_open_orders"`
}

type Position struct {
	Size                 *big.Int `json:"size"`
	OpenNotional         *big.Int `json:"open_notional"`
	LastPremiumFraction  *big.Int `json:"last_premium_fraction"`
	LiquidationThreshold *big.Int `json:"liquidation_threshold"`
}

func GetAMMVariables(stateDB contract.StateDB, ammAddress common.Address, ammIndex int64, trader common.Address) VariablesReadFromAMMSlots {
	lastPrice := getLastPrice(stateDB, ammAddress)
	position := Position{
		Size:                getSize(stateDB, ammAddress, &trader),
		OpenNotional:        getOpenNotional(stateDB, ammAddress, &trader),
		LastPremiumFraction: GetLastPremiumFraction(stateDB, ammAddress, &trader),
	}
	cumulativePremiumFraction := GetCumulativePremiumFraction(stateDB, ammAddress)
	maxOracleSpreadRatio := GetMaxOraclePriceSpread(stateDB, ammIndex)
	maxLiquidationRatio := GetMaxLiquidationRatio(stateDB, ammIndex)
	maxLiquidationPriceSpread := GetMaxLiquidationPriceSpread(stateDB, ammIndex)
	minSizeRequirement := GetMinSizeRequirement(stateDB, ammIndex)
	oracleAddress := getOracleAddress(stateDB)
	underlyingAssetAddress := getUnderlyingAssetAddress(stateDB, ammAddress)
	underlyingPriceForMarket := getUnderlyingPriceForMarket(stateDB, ammIndex)
	underlyingPrice := getUnderlyingPrice(stateDB, ammAddress)
	redStoneAdapterAddress := getRedStoneAdapterAddress(stateDB, oracleAddress)
	redStoneFeedId := getRedStoneFeedId(stateDB, oracleAddress, underlyingAssetAddress)
	bidsHead := getBidsHead(stateDB, ammAddress)
	bidsHeadSize := getBidSize(stateDB, ammAddress, bidsHead)
	asksHead := getAsksHead(stateDB, ammAddress)
	asksHeadSize := getAskSize(stateDB, ammAddress, asksHead)
	upperBound, lowerBound := GetAcceptableBoundsForLiquidation(stateDB, ammIndex)
	minAllowableMargin := GetMinAllowableMargin(stateDB)
	takerFee := GetTakerFee(stateDB)
	totalMargin := GetNormalizedMargin(stateDB, trader)
	availableMargin := GetAvailableMargin(stateDB, trader, 0)
	reduceOnlyAmount := getReduceOnlyAmount(stateDB, trader, big.NewInt(ammIndex))
	longOpenOrdersAmount := getLongOpenOrdersAmount(stateDB, trader, big.NewInt(ammIndex))
	shortOpenOrdersAmount := getShortOpenOrdersAmount(stateDB, trader, big.NewInt(ammIndex))
	impactMarginNotional := getImpactMarginNotional(stateDB, ammAddress)
	return VariablesReadFromAMMSlots{
		LastPrice:                 lastPrice,
		CumulativePremiumFraction: cumulativePremiumFraction,
		MaxOracleSpreadRatio:      maxOracleSpreadRatio,
		OracleAddress:             oracleAddress,
		MaxLiquidationRatio:       maxLiquidationRatio,
		MinSizeRequirement:        minSizeRequirement,
		UnderlyingAssetAddress:    underlyingAssetAddress,
		UnderlyingPriceForMarket:  underlyingPriceForMarket,
		UnderlyingPrice:           underlyingPrice,
		MaxLiquidationPriceSpread: maxLiquidationPriceSpread,
		RedStoneAdapterAddress:    redStoneAdapterAddress,
		RedStoneFeedId:            redStoneFeedId,
		ImpactMarginNotional:      impactMarginNotional,
		Position:                  position,
		BidsHead:                  bidsHead,
		BidsHeadSize:              bidsHeadSize,
		AsksHead:                  asksHead,
		AsksHeadSize:              asksHeadSize,
		UpperBound:                upperBound,
		LowerBound:                lowerBound,
		MinAllowableMargin:        minAllowableMargin,
		TotalMargin:               totalMargin,
		AvailableMargin:           availableMargin,
		TakerFee:                  takerFee,
		ReduceOnlyAmount:          reduceOnlyAmount,
		LongOpenOrders:            longOpenOrdersAmount,
		ShortOpenOrders:           shortOpenOrdersAmount,
	}
}

type VariablesReadFromIOCOrdersSlots struct {
	OrderDetails     OrderDetails `json:"order_details"`
	IocExpirationCap *big.Int     `json:"ioc_expiration_cap"`
}

type OrderDetails struct {
	BlockPlaced  *big.Int `json:"block_placed"`
	FilledAmount *big.Int `json:"filled_amount"`
	OrderStatus  int64    `json:"order_status"`
}

func GetIOCOrdersVariables(stateDB contract.StateDB, orderHash common.Hash) VariablesReadFromIOCOrdersSlots {
	blockPlaced := iocGetBlockPlaced(stateDB, orderHash)
	filledAmount := iocGetOrderFilledAmount(stateDB, orderHash)
	orderStatus := iocGetOrderStatus(stateDB, orderHash)

	iocExpirationCap := iocGetExpirationCap(stateDB)
	return VariablesReadFromIOCOrdersSlots{
		OrderDetails: OrderDetails{
			BlockPlaced:  blockPlaced,
			FilledAmount: filledAmount,
			OrderStatus:  orderStatus,
		},
		IocExpirationCap: iocExpirationCap,
	}
}

type VariablesReadFromOrderbookSlots struct {
	OrderDetails      OrderDetails `json:"order_details"`
	IsTradingAuthoriy bool         `json:"is_trading_authority"`
}

func GetOrderBookVariables(stateDB contract.StateDB, traderAddress string, senderAddress string, orderHash common.Hash) VariablesReadFromOrderbookSlots {
	blockPlaced := getBlockPlaced(stateDB, orderHash)
	filledAmount := getOrderFilledAmount(stateDB, orderHash)
	orderStatus := getOrderStatus(stateDB, orderHash)
	isTradingAuthoriy := IsTradingAuthority(stateDB, common.HexToAddress(traderAddress), common.HexToAddress(senderAddress))
	return VariablesReadFromOrderbookSlots{
		OrderDetails: OrderDetails{
			BlockPlaced:  blockPlaced,
			FilledAmount: filledAmount,
			OrderStatus:  orderStatus,
		},
		IsTradingAuthoriy: isTradingAuthoriy,
	}
}
