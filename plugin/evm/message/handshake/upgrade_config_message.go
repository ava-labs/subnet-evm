// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handshake

import (
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type rawPrecompileUpgrade struct {
	Key   string `serialize:"true"`
	Bytes []byte `serialize:"true"`
}

type networkUpgradeConfigMessage struct {
	OptionalNetworkUpgrades *params.OptionalNetworkUpgrades
	// Config for modifying state as a network upgrade.
	StateUpgrades [][]byte `serialize:"true"`
	// Config for enabling and disabling precompiles as network upgrades.
	PrecompileUpgrades []rawPrecompileUpgrade `serialize:"true"`
}

type UpgradeConfigMessage struct {
	bytes []byte
	hash  common.Hash
}

func (u *UpgradeConfigMessage) Bytes() []byte {
	return u.bytes
}

func (u *UpgradeConfigMessage) ID() common.Hash {
	return u.hash
}

// Attempts to parse a `*params.UpgradeConfig` from a []byte
func NewUpgradeConfigFromBytes(bytes []byte) (*params.UpgradeConfig, error) {
	var config networkUpgradeConfigMessage
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
		err := preCompile.FromBytes(precompileUpgrade.Bytes)
		if err != nil {
			return nil, err
		}
		PrecompileUpgrades = append(PrecompileUpgrades, params.PrecompileUpgrade{Config: preCompile})
	}

	var stateUpgrades []params.StateUpgrade

	for _, bytes := range config.StateUpgrades {
		stateUpgrade := params.StateUpgrade{}
		if err := stateUpgrade.FromBytes(bytes); err != nil {
			return nil, err
		}
		stateUpgrades = append(stateUpgrades, stateUpgrade)
	}

	return &params.UpgradeConfig{
		OptionalNetworkUpgrades: config.OptionalNetworkUpgrades,
		StateUpgrades:           stateUpgrades,
		PrecompileUpgrades:      PrecompileUpgrades,
	}, nil
}

// Wraps an instance of *params.UpgradeConfig
//
// This function returns the serialized UpgradeConfig, ready to be send over to
// other peers. The struct also includes a hash of the content, ready to be used
// as part of the handshake protocol.
//
// Since params.UpgradeConfig should never change without a node reloading, it
// is safe to call this function once and store its output globally to re-use
// multiple times
func NewUpgradeConfigMessage(config *params.UpgradeConfig) (*UpgradeConfigMessage, error) {
	PrecompileUpgrades := make([]rawPrecompileUpgrade, 0)
	for _, precompileConfig := range config.PrecompileUpgrades {
		bytes, err := precompileConfig.Config.ToBytes()
		if err != nil {
			return nil, err
		}
		PrecompileUpgrades = append(PrecompileUpgrades, rawPrecompileUpgrade{
			Key:   precompileConfig.Key(),
			Bytes: bytes,
		})
	}

	stateUpgrades := make([][]byte, 0)

	for _, config := range config.StateUpgrades {
		bytes, err := config.ToBytes()
		if err != nil {
			return nil, err
		}
		stateUpgrades = append(stateUpgrades, bytes)
	}

	wrappedConfig := networkUpgradeConfigMessage{
		OptionalNetworkUpgrades: config.OptionalNetworkUpgrades,
		StateUpgrades:           stateUpgrades,
		PrecompileUpgrades:      PrecompileUpgrades,
	}

	bytes, err := Codec.Marshal(Version, wrappedConfig)
	if err != nil {
		return nil, err
	}

	hash := crypto.Keccak256Hash(bytes)
	return &UpgradeConfigMessage{
		bytes: bytes,
		hash:  hash,
	}, nil
}
