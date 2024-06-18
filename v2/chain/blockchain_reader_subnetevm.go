package chain

import (
	"math/big"

	"github.com/ava-labs/coreth/commontype"
	"github.com/ava-labs/coreth/constants"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/precompile/contracts/feemanager"
	"github.com/ava-labs/coreth/precompile/contracts/rewardmanager"
	"github.com/ethereum/go-ethereum/common"
)

// GetFeeConfigAt returns the fee configuration and the last changed block number at [parent].
// If FeeManager is activated at [parent], returns the fee config in the precompile contract state.
// Otherwise returns the fee config in the chain config.
// Assumes that a valid configuration is stored when the precompile is activated.
func (bc *blockChain) GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error) {
	config := bc.Config()
	if !config.IsPrecompileEnabled(feemanager.ContractAddress, parent.Time) {
		return config.FeeConfig, common.Big0, nil
	}

	// try to return it from the cache
	if cached, hit := bc.feeConfigCache.Get(parent.Root); hit {
		return cached.feeConfig, cached.lastChangedAt, nil
	}

	stateDB, err := bc.StateAt(parent.Root)
	if err != nil {
		return commontype.EmptyFeeConfig, nil, err
	}

	storedFeeConfig := feemanager.GetStoredFeeConfig(stateDB)
	// this should not return an invalid fee config since it's assumed that
	// StoreFeeConfig returns an error when an invalid fee config is attempted to be stored.
	// However an external stateDB call can modify the contract state.
	// This check is added to add a defense in-depth.
	if err := storedFeeConfig.Verify(); err != nil {
		return commontype.EmptyFeeConfig, nil, err
	}
	lastChangedAt := feemanager.GetFeeConfigLastChangedAt(stateDB)
	cacheable := &cacheableFeeConfig{feeConfig: storedFeeConfig, lastChangedAt: lastChangedAt}
	// add it to the cache
	bc.feeConfigCache.Add(parent.Root, cacheable)
	return storedFeeConfig, lastChangedAt, nil
}

// GetCoinbaseAt returns the configured coinbase address at [parent].
// If RewardManager is activated at [parent], returns the reward manager config in the precompile contract state.
// If fee recipients are allowed, returns true in the second return value.
func (bc *blockChain) GetCoinbaseAt(parent *types.Header) (common.Address, bool, error) {
	config := bc.Config()
	if !config.IsSubnetEVM(parent.Time) {
		return constants.BlackholeAddr, false, nil
	}

	if !config.IsPrecompileEnabled(rewardmanager.ContractAddress, parent.Time) {
		if bc.config.AllowFeeRecipients {
			return common.Address{}, true, nil
		} else {
			return constants.BlackholeAddr, false, nil
		}
	}

	// try to return it from the cache
	if cached, hit := bc.coinbaseConfigCache.Get(parent.Root); hit {
		return cached.coinbaseAddress, cached.allowFeeRecipients, nil
	}

	stateDB, err := bc.StateAt(parent.Root)
	if err != nil {
		return common.Address{}, false, err
	}
	rewardAddress, feeRecipients := rewardmanager.GetStoredRewardAddress(stateDB)

	cacheable := &cacheableCoinbaseConfig{coinbaseAddress: rewardAddress, allowFeeRecipients: feeRecipients}
	bc.coinbaseConfigCache.Add(parent.Root, cacheable)
	return rewardAddress, feeRecipients, nil
}
