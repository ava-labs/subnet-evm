// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var _ precompileconfig.Config = &Config{}

// Config implements the StatefulPrecompileConfig interface while adding in the
// ContractNativeMinter specific precompile config.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
	InitialMint map[common.Address]*math.HexOrDecimal256 `json:"initialMint,omitempty"`
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractNativeMinter with the given [admins], [enableds] and [managers] as members of the allowlist.
// Also mints balances according to [initialMint] when the upgrade activates.
func NewConfig(blockTimestamp *uint64, admins []common.Address, enableds []common.Address, managers []common.Address, initialMint map[common.Address]*math.HexOrDecimal256) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
			ManagerAddresses: managers,
		},
		Upgrade:     precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		InitialMint: initialMint,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractNativeMinter.
func NewDisableConfig(blockTimestamp *uint64) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}
func (*Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*ContractNativeMinterConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	eq := c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
	if !eq {
		return false
	}

	if len(c.InitialMint) != len(other.InitialMint) {
		return false
	}

	for address, amount := range c.InitialMint {
		val, ok := other.InitialMint[address]
		if !ok {
			return false
		}
		bigIntAmount := (*big.Int)(amount)
		bigIntVal := (*big.Int)(val)
		if !utils.BigNumEqual(bigIntAmount, bigIntVal) {
			return false
		}
	}

	return true
}

func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	// ensure that all of the initial mint values in the map are non-nil positive values
	for addr, amount := range c.InitialMint {
		if amount == nil {
			return fmt.Errorf("initial mint cannot contain nil amount for address %s", addr)
		}
		bigIntAmount := (*big.Int)(amount)
		if bigIntAmount.Sign() < 1 {
			return fmt.Errorf("initial mint cannot contain invalid amount %v for address %s", bigIntAmount, addr)
		}
	}
	return c.AllowListConfig.Verify(chainConfig, c.Upgrade)
}

func (c *Config) MarshalBinary() ([]byte, error) {
	keys := make([]common.Address, 0)
	for key := range c.InitialMint {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i][:], keys[j][:]) < 0
	})

	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 32 * 1024,
	}

	if err := c.AllowListConfig.ToBytesWithPacker(&p); err != nil {
		return nil, err
	}

	if err := c.Upgrade.ToBytesWithPacker(&p); err != nil {
		return nil, err
	}

	p.PackInt(uint32(len(keys)))
	if p.Err != nil {
		return nil, p.Err
	}

	for _, key := range keys {
		p.PackBytes(key[:])
		if p.Err != nil {
			return nil, p.Err
		}
		p.PackBytes((*big.Int)(c.InitialMint[key]).Bytes())
		if p.Err != nil {
			return nil, p.Err
		}
	}

	return p.Bytes, nil
}

func (c *Config) UnmarshalBinary(bytes []byte) error {
	p := wrappers.Packer{
		Bytes: bytes,
	}
	if err := c.AllowListConfig.FromBytesWithPacker(&p); err != nil {
		return err
	}
	if err := c.Upgrade.FromBytesWithPacker(&p); err != nil {
		return err
	}
	len := p.UnpackInt()
	c.InitialMint = make(map[common.Address]*math.HexOrDecimal256, len)

	for i := uint32(0); i < len; i++ {
		key := common.BytesToAddress(p.UnpackBytes())
		if p.Err != nil {
			return p.Err
		}
		value := p.UnpackBytes()
		if p.Err != nil {
			return p.Err
		}
		c.InitialMint[key] = (*math.HexOrDecimal256)(big.NewInt(0).SetBytes(value))
	}

	return nil
}
