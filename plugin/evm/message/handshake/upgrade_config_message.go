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

type networkUpgradeConfigMessage struct {
	OptionalNetworkUpgrades *params.OptionalNetworkUpgrades `serialize:"true,nullable"`
	// Config for modifying state as a network upgrade.
	StateUpgrades []params.StateUpgrade `serialize:"true"`
	// Config for enabling and disabling precompiles as network upgrades.
	PrecompileUpgrades []rawPrecompileUpgrade `serialize:"true"`
}

type UpgradeConfigMessage struct {
	bytes []byte
	hash  []byte
}

func (u *UpgradeConfigMessage) Bytes() []byte {
	return u.bytes
}

func (u *UpgradeConfigMessage) ID() []byte {
	return u.hash
}

// Attempts to parse a networkUpgradeConfigMessage from a []byte
//
// This function attempts to parse a stream of bytes as a
// networkUpgradeConfigMessage (as serialized from
// UpgradeConfigToNetworkMessage).
//
// The function returns a reference of *params.UpgradeConfig
func NewUpgradeConfigMessageFromBytes(bytes []byte) (*params.UpgradeConfig, error) {
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
		version, err := Codec.Unmarshal(precompileUpgrade.Bytes, preCompile)
		if version != Version {
			return nil, ErrInvalidVersion
		}
		if err != nil {
			return nil, err
		}
		PrecompileUpgrades = append(PrecompileUpgrades, params.PrecompileUpgrade{Config: preCompile})
	}

	return &params.UpgradeConfig{
		OptionalNetworkUpgrades: config.OptionalNetworkUpgrades,
		StateUpgrades:           config.StateUpgrades,
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
		bytes, err := Codec.Marshal(Version, precompileConfig.Config)
		if err != nil {
			return nil, err
		}
		PrecompileUpgrades = append(PrecompileUpgrades, rawPrecompileUpgrade{
			Key:   precompileConfig.Key(),
			Bytes: bytes,
		})
	}

	wrappedConfig := networkUpgradeConfigMessage{
		OptionalNetworkUpgrades: config.OptionalNetworkUpgrades,
		StateUpgrades:           config.StateUpgrades,
		PrecompileUpgrades:      PrecompileUpgrades,
	}

	bytes, err := Codec.Marshal(Version, wrappedConfig)
	if err != nil {
		return nil, err
	}

	hash := crypto.Keccak256(bytes)
	return &UpgradeConfigMessage{
		bytes: bytes,
		hash:  hash[:8],
	}, nil
}
