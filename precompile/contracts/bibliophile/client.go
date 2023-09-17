package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

type BibliophileClient interface {
	//margin account
	GetAvailableMargin(trader common.Address) *big.Int
	//clearing house
	GetMarketAddressFromMarketID(marketId int64) common.Address
	GetMinAllowableMargin() *big.Int
	GetTakerFee() *big.Int
	//orderbook
	GetSize(market common.Address, trader *common.Address) *big.Int
	GetLongOpenOrdersAmount(trader common.Address, ammIndex *big.Int) *big.Int
	GetShortOpenOrdersAmount(trader common.Address, ammIndex *big.Int) *big.Int
	GetReduceOnlyAmount(trader common.Address, ammIndex *big.Int) *big.Int
	IsTradingAuthority(trader, senderOrSigner common.Address) bool
	IsValidator(senderOrSigner common.Address) bool
	// Limit Order
	GetBlockPlaced(orderHash [32]byte) *big.Int
	GetOrderFilledAmount(orderHash [32]byte) *big.Int
	GetOrderStatus(orderHash [32]byte) int64
	// IOC Order
	IOC_GetBlockPlaced(orderHash [32]byte) *big.Int
	IOC_GetOrderFilledAmount(orderHash [32]byte) *big.Int
	IOC_GetOrderStatus(orderHash [32]byte) int64
	IOC_GetExpirationCap() *big.Int

	// AMM
	GetMinSizeRequirement(marketId int64) *big.Int
	GetLastPrice(ammAddress common.Address) *big.Int
	GetBidSize(ammAddress common.Address, price *big.Int) *big.Int
	GetAskSize(ammAddress common.Address, price *big.Int) *big.Int
	GetNextBidPrice(ammAddress common.Address, price *big.Int) *big.Int
	GetNextAskPrice(ammAddress common.Address, price *big.Int) *big.Int
	GetImpactMarginNotional(ammAddress common.Address) *big.Int
	GetBidsHead(market common.Address) *big.Int
	GetAsksHead(market common.Address) *big.Int
	GetUpperAndLowerBoundForMarket(marketId int64) (*big.Int, *big.Int)
	GetAcceptableBoundsForLiquidation(marketId int64) (*big.Int, *big.Int)

	GetAccessibleState() contract.AccessibleState
	GetNotionalPositionAndMargin(trader common.Address, includeFundingPayments bool, mode uint8) (*big.Int, *big.Int)
}

// Define a structure that will implement the Bibliophile interface
type bibliophileClient struct {
	accessibleState contract.AccessibleState
}

func NewBibliophileClient(accessibleState contract.AccessibleState) BibliophileClient {
	return &bibliophileClient{
		accessibleState: accessibleState,
	}
}

func (b *bibliophileClient) GetAccessibleState() contract.AccessibleState {
	return b.accessibleState
}

func (b *bibliophileClient) GetSize(market common.Address, trader *common.Address) *big.Int {
	return getSize(b.accessibleState.GetStateDB(), market, trader)
}

func (b *bibliophileClient) GetMinSizeRequirement(marketId int64) *big.Int {
	return GetMinSizeRequirement(b.accessibleState.GetStateDB(), marketId)
}

func (b *bibliophileClient) GetMinAllowableMargin() *big.Int {
	return GetMinAllowableMargin(b.accessibleState.GetStateDB())
}

func (b *bibliophileClient) GetTakerFee() *big.Int {
	return GetTakerFee(b.accessibleState.GetStateDB())
}

func (b *bibliophileClient) GetMarketAddressFromMarketID(marketID int64) common.Address {
	return getMarketAddressFromMarketID(marketID, b.accessibleState.GetStateDB())
}

func (b *bibliophileClient) GetBlockPlaced(orderHash [32]byte) *big.Int {
	return getBlockPlaced(b.accessibleState.GetStateDB(), orderHash)
}

func (b *bibliophileClient) GetOrderFilledAmount(orderHash [32]byte) *big.Int {
	return getOrderFilledAmount(b.accessibleState.GetStateDB(), orderHash)
}

func (b *bibliophileClient) GetOrderStatus(orderHash [32]byte) int64 {
	return getOrderStatus(b.accessibleState.GetStateDB(), orderHash)
}

func (b *bibliophileClient) IOC_GetBlockPlaced(orderHash [32]byte) *big.Int {
	return iocGetBlockPlaced(b.accessibleState.GetStateDB(), orderHash)
}

func (b *bibliophileClient) IOC_GetOrderFilledAmount(orderHash [32]byte) *big.Int {
	return iocGetOrderFilledAmount(b.accessibleState.GetStateDB(), orderHash)
}

func (b *bibliophileClient) IOC_GetOrderStatus(orderHash [32]byte) int64 {
	return iocGetOrderStatus(b.accessibleState.GetStateDB(), orderHash)
}

func (b *bibliophileClient) IsTradingAuthority(trader, senderOrSigner common.Address) bool {
	return IsTradingAuthority(b.accessibleState.GetStateDB(), trader, senderOrSigner)
}

func (b *bibliophileClient) IsValidator(senderOrSigner common.Address) bool {
	return IsValidator(b.accessibleState.GetStateDB(), senderOrSigner)
}

func (b *bibliophileClient) IOC_GetExpirationCap() *big.Int {
	return iocGetExpirationCap(b.accessibleState.GetStateDB())
}

func (b *bibliophileClient) GetLastPrice(ammAddress common.Address) *big.Int {
	return getLastPrice(b.accessibleState.GetStateDB(), ammAddress)
}

func (b *bibliophileClient) GetBidSize(ammAddress common.Address, price *big.Int) *big.Int {
	return getBidSize(b.accessibleState.GetStateDB(), ammAddress, price)
}

func (b *bibliophileClient) GetAskSize(ammAddress common.Address, price *big.Int) *big.Int {
	return getAskSize(b.accessibleState.GetStateDB(), ammAddress, price)
}

func (b *bibliophileClient) GetNextBidPrice(ammAddress common.Address, price *big.Int) *big.Int {
	return getNextBid(b.accessibleState.GetStateDB(), ammAddress, price)
}

func (b *bibliophileClient) GetNextAskPrice(ammAddress common.Address, price *big.Int) *big.Int {
	return getNextAsk(b.accessibleState.GetStateDB(), ammAddress, price)
}

func (b *bibliophileClient) GetImpactMarginNotional(ammAddress common.Address) *big.Int {
	return getImpactMarginNotional(b.accessibleState.GetStateDB(), ammAddress)
}

func (b *bibliophileClient) GetUpperAndLowerBoundForMarket(marketId int64) (*big.Int, *big.Int) {
	return GetAcceptableBounds(b.accessibleState.GetStateDB(), marketId)
}

func (b *bibliophileClient) GetAcceptableBoundsForLiquidation(marketId int64) (*big.Int, *big.Int) {
	return GetAcceptableBoundsForLiquidation(b.accessibleState.GetStateDB(), marketId)
}

func (b *bibliophileClient) GetBidsHead(market common.Address) *big.Int {
	return getBidsHead(b.accessibleState.GetStateDB(), market)
}

func (b *bibliophileClient) GetAsksHead(market common.Address) *big.Int {
	return getAsksHead(b.accessibleState.GetStateDB(), market)
}

func (b *bibliophileClient) GetLongOpenOrdersAmount(trader common.Address, ammIndex *big.Int) *big.Int {
	return getLongOpenOrdersAmount(b.accessibleState.GetStateDB(), trader, ammIndex)
}

func (b *bibliophileClient) GetShortOpenOrdersAmount(trader common.Address, ammIndex *big.Int) *big.Int {
	return getShortOpenOrdersAmount(b.accessibleState.GetStateDB(), trader, ammIndex)
}

func (b *bibliophileClient) GetReduceOnlyAmount(trader common.Address, ammIndex *big.Int) *big.Int {
	return getReduceOnlyAmount(b.accessibleState.GetStateDB(), trader, ammIndex)
}

func (b *bibliophileClient) GetAvailableMargin(trader common.Address) *big.Int {
	return GetAvailableMargin(b.accessibleState.GetStateDB(), trader)
}

func (b *bibliophileClient) GetNotionalPositionAndMargin(trader common.Address, includeFundingPayments bool, mode uint8) (*big.Int, *big.Int) {
	output := getNotionalPositionAndMargin(b.accessibleState.GetStateDB(), &GetNotionalPositionAndMarginInput{Trader: trader, IncludeFundingPayments: includeFundingPayments, Mode: mode})
	return output.NotionalPosition, output.Margin
}
