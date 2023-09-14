// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handshake

import (
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ethereum/go-ethereum/crypto"
)

type rawPrecompileUpgrade struct {
	Key   string `serialize:"true"`
	Bytes []byte `serialize:"true"`
}

type upgradeConfigMessage struct {
	OptionalNetworkUpgrades []params.Fork `serialize:"true"`

	// Config for modifying state as a network upgrade.
	StateUpgrades []params.StateUpgrade `serialize:"true"`

	// Config for enabling and disabling precompiles as network upgrades.
	PrecompileUpgrades []rawPrecompileUpgrade `serialize:"true"`
	config             params.UpgradeConfig
	bytes              []byte
}

func ParseUpgradeConfig(bytes []byte) (*upgradeConfigMessage, error) {
	var config upgradeConfigMessage
	version, err := Codec.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	if version != Version {
		return nil, ErrInvalidVersion
	}

	var PrecompileUpgrades []params.PrecompileUpgrade

	for _, precompileUpgrade := range config.PrecompileUpgrades {
		module, ok := modules.GetPrecompileModule(precompileUpgrade.Key)
		if !ok {
			return nil, ErrUnknowPrecompile
		}
		preCompile := module.MakeConfig()

		version, err := Codec.Unmarshal(precompileUpgrade.Bytes, preCompile)
		if version != Version {
			return nil, ErrInvalidVersion
		}
		if err != nil {
			return nil, err
		}
		PrecompileUpgrades = append(PrecompileUpgrades, params.PrecompileUpgrade{Config: preCompile})
	}

	config.config = params.UpgradeConfig{
		OptionalNetworkUpgrades: &params.OptionalNetworkUpgrades{Updates: config.OptionalNetworkUpgrades},
		StateUpgrades:           config.StateUpgrades,
		PrecompileUpgrades:      PrecompileUpgrades,
	}
	config.bytes = bytes

	return &config, nil
}

func NewUpgradeConfig(config params.UpgradeConfig) (*upgradeConfigMessage, error) {
	PrecompileUpgrades := make([]rawPrecompileUpgrade, 0)
	for _, precompileConfig := range config.PrecompileUpgrades {
		bytes, err := Codec.Marshal(Version, precompileConfig.Config)
		if err != nil {
			return nil, err
		}
		PrecompileUpgrades = append(PrecompileUpgrades, rawPrecompileUpgrade{
			Key:   precompileConfig.Key(),
			Bytes: bytes,
		})
	}

	optionalNetworkUpgrades := make([]params.Fork, 0)
	if config.OptionalNetworkUpgrades != nil {
		optionalNetworkUpgrades = config.OptionalNetworkUpgrades.Updates
	}

	wrappedConfig := upgradeConfigMessage{
		OptionalNetworkUpgrades: optionalNetworkUpgrades,
		StateUpgrades:           config.StateUpgrades,
		PrecompileUpgrades:      PrecompileUpgrades,
		config:                  config,
		bytes:                   make([]byte, 0),
	}
	bytes, err := Codec.Marshal(Version, wrappedConfig)
	if err != nil {
		return nil, err
	}
	wrappedConfig.bytes = bytes

	return &wrappedConfig, nil
}

func (r *upgradeConfigMessage) Config() params.UpgradeConfig {
	return r.config
}

func (r *upgradeConfigMessage) Bytes() []byte {
	return r.bytes
}

func (r *upgradeConfigMessage) Hash() [8]byte {
	hash := crypto.Keccak256(r.bytes)
	var firstBytes [8]byte
	copy(firstBytes[:], hash[:8])
	return firstBytes
}
