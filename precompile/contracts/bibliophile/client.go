package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

type BibliophileClient interface {
	GetSize(market common.Address, trader *common.Address) *big.Int
	GetMinSizeRequirement(marketId int64) *big.Int
	GetMarketAddressFromMarketID(marketId int64) common.Address
	DetermineFillPrice(marketId int64, longOrderPrice, shortOrderPrice, blockPlaced0, blockPlaced1 *big.Int) (*ValidateOrdersAndDetermineFillPriceOutput, error)
	DetermineLiquidationFillPrice(marketId int64, baseAssetQuantity, price *big.Int) (*big.Int, error)

	// Misc
	IsTradingAuthority(senderOrSigner, trader common.Address) bool

	// Limit Order
	GetBlockPlaced(orderHash [32]byte) *big.Int
	GetOrderFilledAmount(orderHash [32]byte) *big.Int
	GetOrderStatus(orderHash [32]byte) int64

	// IOC Order
	IOC_GetBlockPlaced(orderHash [32]byte) *big.Int
	IOC_GetOrderFilledAmount(orderHash [32]byte) *big.Int
	IOC_GetOrderStatus(orderHash [32]byte) int64
	IOC_GetExpirationCap() *big.Int

	GetAccessibleState() contract.AccessibleState
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

func (b *bibliophileClient) GetMarketAddressFromMarketID(marketID int64) common.Address {
	return getMarketAddressFromMarketID(marketID, b.accessibleState.GetStateDB())
}

func (b *bibliophileClient) DetermineFillPrice(marketId int64, longOrderPrice, shortOrderPrice, blockPlaced0, blockPlaced1 *big.Int) (*ValidateOrdersAndDetermineFillPriceOutput, error) {
	return DetermineFillPrice(b.accessibleState.GetStateDB(), marketId, longOrderPrice, shortOrderPrice, blockPlaced0, blockPlaced1)
}

func (b *bibliophileClient) DetermineLiquidationFillPrice(marketId int64, baseAssetQuantity, price *big.Int) (*big.Int, error) {
	return DetermineLiquidationFillPrice(b.accessibleState.GetStateDB(), marketId, baseAssetQuantity, price)
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

func (b *bibliophileClient) IOC_GetExpirationCap() *big.Int {
	return iocGetExpirationCap(b.accessibleState.GetStateDB())
}
