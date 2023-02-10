// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated
// This file is a generated precompile contract with stubbed abstract functions.

package rewardmanager

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

var _ config.Config = &RewardManagerConfig{}

type InitialRewardConfig struct {
	AllowFeeRecipients bool           `json:"allowFeeRecipients"`
	RewardAddress      common.Address `json:"rewardAddress,omitempty"`
}

func (i *InitialRewardConfig) Equal(other *InitialRewardConfig) bool {
	if other == nil {
		return false
	}

	return i.AllowFeeRecipients == other.AllowFeeRecipients && i.RewardAddress == other.RewardAddress
}

func (i *InitialRewardConfig) Verify() error {
	switch {
	case i.AllowFeeRecipients && i.RewardAddress != (common.Address{}):
		return ErrCannotEnableBothRewards
	default:
		return nil
	}
}

func (i *InitialRewardConfig) Configure(state contract.StateDB) error {
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
	config.UpgradeableConfig
	InitialRewardConfig *InitialRewardConfig `json:"initialRewardConfig,omitempty"`
}

// NewRewardManagerConfig returns a config for a network upgrade at [blockTimestamp] that enables
// RewardManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial rewards config if specified.
func NewRewardManagerConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *InitialRewardConfig) *RewardManagerConfig {
	return &RewardManagerConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig:   config.UpgradeableConfig{BlockTimestamp: blockTimestamp},
		InitialRewardConfig: initialConfig,
	}
}

// NewDisableRewardManagerConfig returns config for a network upgrade at [blockTimestamp]
// that disables RewardManager.
func NewDisableRewardManagerConfig(blockTimestamp *big.Int) *RewardManagerConfig {
	return &RewardManagerConfig{
		UpgradeableConfig: config.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (RewardManagerConfig) Address() common.Address { return ContractAddress }

func (RewardManagerConfig) Key() string { return ConfigKey }

func (c *RewardManagerConfig) Verify() error {
	if err := c.AllowListConfig.Verify(); err != nil {
		return err
	}
	if c.InitialRewardConfig != nil {
		return c.InitialRewardConfig.Verify()
	}
	return nil
}

// Equal returns true if [cfg] is a [*RewardManagerConfig] and it has been configured identical to [c].
func (c *RewardManagerConfig) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*RewardManagerConfig)
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
