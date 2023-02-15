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

var _ config.Config = &Config{}

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

// Config implements the StatefulPrecompileConfig
// interface while adding in the RewardManager specific precompile config.
type Config struct {
	allowlist.Config
	config.Uprade
	InitialRewardConfig *InitialRewardConfig `json:"initialRewardConfig,omitempty"`
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// RewardManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial rewards config if specified.
func NewConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *InitialRewardConfig) *Config {
	return &Config{
		Config: allowlist.Config{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Uprade:              config.Uprade{BlockTimestamp: blockTimestamp},
		InitialRewardConfig: initialConfig,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables RewardManager.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Uprade: config.Uprade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (Config) Key() string { return ConfigKey }

func (c *Config) Verify() error {
	if err := c.Config.Verify(); err != nil {
		return err
	}
	if c.InitialRewardConfig != nil {
		return c.InitialRewardConfig.Verify()
	}
	return nil
}

// Equal returns true if [cfg] is a [*RewardManagerConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}

	equals := c.Uprade.Equal(&other.Uprade) && c.Config.Equal(&other.Config)
	if !equals {
		return false
	}

	if c.InitialRewardConfig == nil {
		return other.InitialRewardConfig == nil
	}

	return c.InitialRewardConfig.Equal(other.InitialRewardConfig)
}
