package acp224feemanager

import (
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/subnet-evm/commontype"
)

var _ precompileconfig.Config = (*Config)(nil)

// Config implements the precompileconfig.Config interface and
// adds specific configuration for ACP224FeeManager.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
	InitialFeeConfig *commontype.ACP224FeeConfig `json:"initialFeeConfig,omitempty"` // initial fee config to be immediately activated
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ACP224FeeManager with the given [admins], [enableds] and [managers] members of the allowlist .
func NewConfig(
	blockTimestamp *uint64,
	admins []common.Address,
	enableds []common.Address,
	managers []common.Address,
	initialConfig *commontype.ACP224FeeConfig,
) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
			ManagerAddresses: managers,
		},
		Upgrade:          precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		InitialFeeConfig: initialConfig,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables ACP224FeeManager.
func NewDisableConfig(blockTimestamp *uint64) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Key returns the key for the ACP224FeeManager precompileconfig.
// This should be the same key as used in the precompile module.
func (*Config) Key() string { return ConfigKey }

// Equal returns true if [s] is a [*ACP224FeeManagerConfig] and it has been configured identical to [c].
func (c *Config) Equal(s precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (s).(*Config)
	if !ok {
		return false
	}

	eq := c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
	if !eq {
		return false
	}

	if c.InitialFeeConfig == nil {
		return other.InitialFeeConfig == nil
	}

	return c.InitialFeeConfig.Equal(other.InitialFeeConfig)
}

// Verify tries to verify Config and returns an error accordingly.
func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	if err := c.AllowListConfig.Verify(chainConfig, c.Upgrade); err != nil {
		return err
	}

	if c.InitialFeeConfig == nil {
		return nil
	}

	return c.InitialFeeConfig.Verify()
}
