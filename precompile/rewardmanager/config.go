// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated
// This file is a generated precompile contract with stubbed abstract functions.

package rewardmanager

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompile.StatefulPrecompileConfig = &RewardManagerConfig{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "rewardManagerConfig"

type InitialRewardConfig struct {
	AllowFeeRecipients bool           `json:"allowFeeRecipients"`
	RewardAddress      common.Address `json:"rewardAddress,omitempty"`
}

func (i *InitialRewardConfig) Verify() error {
	switch {
	case i.AllowFeeRecipients && i.RewardAddress != (common.Address{}):
		return ErrCannotEnableBothRewards
	default:
		return nil
	}
}

func (i *InitialRewardConfig) Equal(other *InitialRewardConfig) bool {
	if other == nil {
		return false
	}

	return i.AllowFeeRecipients == other.AllowFeeRecipients && i.RewardAddress == other.RewardAddress
}

func (i *InitialRewardConfig) Configure(state precompile.StateDB) error {
	// enable allow fee recipients
	if i.AllowFeeRecipients {
		EnableAllowFeeRecipients(state)
	} else if i.RewardAddress == (common.Address{}) {
		// if reward address is empty and allow fee recipients is false
		// then disable rewards
		DisableFeeRewards(state)
	} else {
		// set reward address
		return StoreRewardAddress(state, i.RewardAddress)
	}
	return nil
}

// RewardManagerConfig implements the StatefulPrecompileConfig
// interface while adding in the RewardManager specific precompile config.
type RewardManagerConfig struct {
	allowlist.AllowListConfig
	precompile.UpgradeableConfig
	InitialRewardConfig *InitialRewardConfig `json:"initialRewardConfig,omitempty"`
}

// NewModule returns a new module for RewardManager.
func NewModule() precompile.StatefulPrecompileModule {
	return &RewardManagerConfig{}
}

// NewRewardManagerConfig returns a config for a network upgrade at [blockTimestamp] that enables
// RewardManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial rewards config if specified.
func NewRewardManagerConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *InitialRewardConfig) *RewardManagerConfig {
	return &RewardManagerConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig:   precompile.UpgradeableConfig{BlockTimestamp: blockTimestamp},
		InitialRewardConfig: initialConfig,
	}
}

// NewDisableRewardManagerConfig returns config for a network upgrade at [blockTimestamp]
// that disables RewardManager.
func NewDisableRewardManagerConfig(blockTimestamp *big.Int) *RewardManagerConfig {
	return &RewardManagerConfig{
		UpgradeableConfig: precompile.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Address returns the address of the RewardManager. Addresses reside under the precompile/params.go
// Select a non-conflicting address and set it in the params.go.
func (*RewardManagerConfig) Address() common.Address {
	return ContractAddress
}

// Contract returns the singleton stateful precompiled contract to be used for RewardManager.
func (*RewardManagerConfig) Contract() precompile.StatefulPrecompiledContract {
	return RewardManagerPrecompile
}

// Key returns the key used in json config files to specify this precompile config.
func (*RewardManagerConfig) Key() string {
	return ConfigKey
}

// NewConfig returns a new instance of RewardManagerConfig.
func (*RewardManagerConfig) NewConfig() precompile.StatefulPrecompileConfig {
	return new(RewardManagerConfig)
}

// Equal returns true if [s] is a [*RewardManagerConfig] and it has been configured identical to [c].
func (c *RewardManagerConfig) Equal(s precompile.StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*RewardManagerConfig)
	if !ok {
		return false
	}
	// modify this boolean accordingly with your custom RewardManagerConfig, to check if [other] and the current [c] are equal
	// if RewardManagerConfig contains only UpgradeableConfig and precompile.AllowListConfig you can skip modifying it.
	equals := c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
	if !equals {
		return false
	}

	if c.InitialRewardConfig == nil {
		return other.InitialRewardConfig == nil
	}

	return c.InitialRewardConfig.Equal(other.InitialRewardConfig)
}

// Configure configures [state] with the initial configuration.
func (c *RewardManagerConfig) Configure(chainConfig precompile.ChainConfig, state precompile.StateDB, _ precompile.BlockContext) error {
	c.AllowListConfig.Configure(state, ContractAddress)
	// configure the RewardManager with the given initial configuration
	if c.InitialRewardConfig != nil {
		return c.InitialRewardConfig.Configure(state)
	} else if chainConfig.AllowedFeeRecipients() {
		// configure the RewardManager according to chainConfig
		EnableAllowFeeRecipients(state)
	} else {
		// chainConfig does not have any reward address
		// if chainConfig does not enable fee recipients
		// default to disabling rewards
		DisableFeeRewards(state)
	}
	return nil
}

func (c *RewardManagerConfig) Verify() error {
	if err := c.AllowListConfig.Verify(); err != nil {
		return err
	}
	if c.InitialRewardConfig != nil {
		return c.InitialRewardConfig.Verify()
	}
	return nil
}
