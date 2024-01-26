// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated
// This file is a generated precompile contract with stubbed abstract functions.

package rewardmanager

import (
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/docker/docker/pkg/units"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileconfig.Config = &Config{}

type InitialRewardConfig struct {
	AllowFeeRecipients bool           `json:"allowFeeRecipients"`
	RewardAddress      common.Address `json:"rewardAddress,omitempty"`
}

func (u *InitialRewardConfig) MarshalBinary() ([]byte, error) {
	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 1 * units.MiB,
	}
	p.PackBool(u.AllowFeeRecipients)
	if p.Err != nil {
		return nil, p.Err
	}
	p.PackFixedBytes(u.RewardAddress[:])
	return p.Bytes, p.Err
}

func (u *InitialRewardConfig) UnmarshalBinary(data []byte) error {
	p := wrappers.Packer{
		Bytes: data,
	}
	u.AllowFeeRecipients = p.UnpackBool()
	if p.Err != nil {
		return p.Err
	}
	u.RewardAddress = common.BytesToAddress(p.UnpackFixedBytes(common.AddressLength))
	if p.Err != nil {
		return p.Err
	}
	return nil
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
		StoreRewardAddress(state, i.RewardAddress)
	}
	return nil
}

// Config implements the StatefulPrecompileConfig interface while adding in the
// RewardManager specific precompile config.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
	InitialRewardConfig *InitialRewardConfig `json:"initialRewardConfig,omitempty"`
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// RewardManager with the given [admins], [enableds] and [managers] as members of the allowlist with [initialConfig] as initial rewards config if specified.
func NewConfig(blockTimestamp *uint64, admins []common.Address, enableds []common.Address, managers []common.Address, initialConfig *InitialRewardConfig) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
			ManagerAddresses: managers,
		},
		Upgrade:             precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		InitialRewardConfig: initialConfig,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables RewardManager.
func NewDisableConfig(blockTimestamp *uint64) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Key returns the key for the Contract precompileconfig.
// This should be the same key as used in the precompile module.
func (*Config) Key() string { return ConfigKey }

// Verify tries to verify Config and returns an error accordingly.
func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	if c.InitialRewardConfig != nil {
		if err := c.InitialRewardConfig.Verify(); err != nil {
			return err
		}
	}
	return c.AllowListConfig.Verify(chainConfig, c.Upgrade)
}

// Equal returns true if [cfg] is a [*RewardManagerConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}

	if c.InitialRewardConfig != nil {
		if other.InitialRewardConfig == nil {
			return false
		}
		if !c.InitialRewardConfig.Equal(other.InitialRewardConfig) {
			return false
		}
	}

	return c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
}

func (c *Config) MarshalBinary() ([]byte, error) {
	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 1 * units.MiB,
	}

	bytes, err := c.AllowListConfig.MarshalBinary()
	if err != nil {
		return nil, err
	}

	p.PackBytes(bytes)
	if p.Err != nil {
		return nil, p.Err
	}

	bytes, err = c.Upgrade.MarshalBinary()
	if err != nil {
		return nil, err
	}
	p.PackBytes(bytes)
	if p.Err != nil {
		return nil, p.Err
	}

	p.PackBool(c.InitialRewardConfig == nil)
	if p.Err != nil {
		return nil, p.Err
	}

	if c.InitialRewardConfig != nil {
		bytes, err := c.InitialRewardConfig.MarshalBinary()
		if err != nil {
			return nil, err
		}
		p.PackBytes(bytes)
	}

	return p.Bytes, p.Err
}

func (c *Config) UnmarshalBinary(bytes []byte) error {
	p := wrappers.Packer{
		Bytes: bytes,
	}
	allowList := p.UnpackBytes()
	if p.Err != nil {
		return p.Err
	}
	upgrade := p.UnpackBytes()
	if p.Err != nil {
		return p.Err
	}
	if err := c.AllowListConfig.UnmarshalBinary(allowList); err != nil {
		return err
	}
	if err := c.Upgrade.UnmarshalBinary(upgrade); err != nil {
		return err
	}

	isNil := p.UnpackBool()
	if p.Err == nil && !isNil {
		c.InitialRewardConfig = &InitialRewardConfig{}
		bytes := p.UnpackBytes()
		if p.Err != nil {
			return p.Err
		}
		return c.InitialRewardConfig.UnmarshalBinary(bytes)
	}

	return p.Err
}
