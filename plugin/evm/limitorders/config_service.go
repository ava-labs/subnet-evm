package limitorders

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contracts/hubbleconfigmanager"
)

type IConfigService interface {
	getSpreadRatioThreshold() *big.Int
	getMaxLiquidationRatio() *big.Int
	getMinAllowableMargin() *big.Int
	getMaintenanceMargin() *big.Int
	getMinSizeRequirement() *big.Int
}

type ConfigService struct {
	blockChain *core.BlockChain
}

func NewConfigService(blockChain *core.BlockChain) IConfigService {
	return &ConfigService{
		blockChain: blockChain,
	}
}

func (cs *ConfigService) getSpreadRatioThreshold() *big.Int {
	return hubbleconfigmanager.GetSpreadRatioThreshold(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMaxLiquidationRatio() *big.Int {
	return hubbleconfigmanager.GetMaxLiquidationRatio(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMinAllowableMargin() *big.Int {
	return hubbleconfigmanager.GetMinAllowableMargin(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMaintenanceMargin() *big.Int {
	return hubbleconfigmanager.GetMaintenanceMargin(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getMinSizeRequirement() *big.Int {
	return hubbleconfigmanager.GetMinSizeRequirement(cs.getStateAtCurrentBlock())
}

func (cs *ConfigService) getStateAtCurrentBlock() *state.StateDB {
	stateDB, _ := cs.blockChain.StateAt(cs.blockChain.CurrentBlock().Root())
	return stateDB
}
