package limitorders

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contracts/hubblebibliophile"
)

type IConfigService interface {
	getSpreadRatioThreshold(market Market) *big.Int
	getMaxLiquidationRatio(market Market) *big.Int
	getMinAllowableMargin() *big.Int
	getMaintenanceMargin() *big.Int
	getMinSizeRequirement(market Market) *big.Int
	GetActiveMarketsCount() int64
}

type ConfigService struct {
	blockChain *core.BlockChain
}

func NewConfigService(blockChain *core.BlockChain) IConfigService {
	return &ConfigService{
		blockChain: blockChain,
	}
}

func (cs *ConfigService) getSpreadRatioThreshold(market Market) *big.Int {
	return hubblebibliophile.GetMaxOracleSpreadRatioForMarket(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) getMaxLiquidationRatio(market Market) *big.Int {
	return hubblebibliophile.GetMaxLiquidationRatioForMarket(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) getMinAllowableMargin() *big.Int {
	return hubblebibliophile.GetMinAllowableMargin(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMaintenanceMargin() *big.Int {
	return hubblebibliophile.GetMaintenanceMargin(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMinSizeRequirement(market Market) *big.Int {
	return hubblebibliophile.GetMinSizeRequirementForMarket(cs.getStateAtCurrentBlock(), int64(market))
}

func (cs *ConfigService) getStateAtCurrentBlock() *state.StateDB {
	stateDB, _ := cs.blockChain.StateAt(cs.blockChain.CurrentBlock().Root())
	return stateDB
}

func (cs *ConfigService) GetActiveMarketsCount() int64 {
	return hubblebibliophile.GetActiveMarketsCount(cs.getStateAtCurrentBlock())
}
