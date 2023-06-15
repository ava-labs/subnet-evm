package limitorders

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contracts/hubblebibliophile"
	"github.com/ethereum/go-ethereum/common"
)

type IConfigService interface {
	getMaxLiquidationRatio(market Market) *big.Int
	getLiquidationSpreadThreshold(market Market) *big.Int
	getMinAllowableMargin() *big.Int
	getMaintenanceMargin() *big.Int
	getMinSizeRequirement(market Market) *big.Int
	GetActiveMarketsCount() int64
	GetUnderlyingPrices() []*big.Int
	GetLastPremiumFraction(market Market, trader *common.Address) *big.Int
	GetCumulativePremiumFraction(market Market) *big.Int
	GetAcceptableBounds(market Market) (*big.Int, *big.Int)
	GetAcceptableBoundsForLiquidation(market Market) (*big.Int, *big.Int)
}

type ConfigService struct {
	blockChain *core.BlockChain
}

func NewConfigService(blockChain *core.BlockChain) IConfigService {
	return &ConfigService{
		blockChain: blockChain,
	}
}

func (cs *ConfigService) GetAcceptableBounds(market Market) (*big.Int, *big.Int) {
	return hubblebibliophile.GetAcceptableBounds(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) GetAcceptableBoundsForLiquidation(market Market) (*big.Int, *big.Int) {
	return hubblebibliophile.GetAcceptableBoundsForLiquidation(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) getLiquidationSpreadThreshold(market Market) *big.Int {
	return hubblebibliophile.GetMaxLiquidationPriceSpread(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) getMaxLiquidationRatio(market Market) *big.Int {
	return hubblebibliophile.GetMaxLiquidationRatio(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) getMinAllowableMargin() *big.Int {
	return hubblebibliophile.GetMinAllowableMargin(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMaintenanceMargin() *big.Int {
	return hubblebibliophile.GetMaintenanceMargin(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMinSizeRequirement(market Market) *big.Int {
	return hubblebibliophile.GetMinSizeRequirement(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) getStateAtCurrentBlock() *state.StateDB {
	stateDB, _ := cs.blockChain.StateAt(cs.blockChain.CurrentBlock().Root())
	return stateDB
}

func (cs *ConfigService) GetActiveMarketsCount() int64 {
	return hubblebibliophile.GetActiveMarketsCount(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) GetUnderlyingPrices() []*big.Int {
	return hubblebibliophile.GetUnderlyingPrices(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) GetLastPremiumFraction(market Market, trader *common.Address) *big.Int {
	markets := hubblebibliophile.GetMarkets(cs.getStateAtCurrentBlock())
	return hubblebibliophile.GetLastPremiumFraction(cs.getStateAtCurrentBlock(), markets[market], trader)
}

func (cs *ConfigService) GetCumulativePremiumFraction(market Market) *big.Int {
	markets := hubblebibliophile.GetMarkets(cs.getStateAtCurrentBlock())
	return hubblebibliophile.GetCumulativePremiumFraction(cs.getStateAtCurrentBlock(), markets[market])
}
