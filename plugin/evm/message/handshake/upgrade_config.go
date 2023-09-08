package handshake

import (
	"crypto/sha256"
	"fmt"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/modules"
)

type PrecompileUpgrade struct {
	StructName string `serialize:"true"`
	Bytes      []byte `serialize:"true"`
}

type UpgradeConfig struct {
	OptionalNetworkUpgrades []params.Fork `serialize:"true"`

	// Config for modifying state as a network upgrade.
	StateUpgrades []params.StateUpgrade `serialize:"true"`

	// Config for enabling and disabling precompiles as network upgrades.
	PrecompileUpgrades []PrecompileUpgrade `serialize:"true"`
	config             params.UpgradeConfig
	bytes              []byte
}

func ParseUpgradeConfig(bytes []byte) (*UpgradeConfig, error) {
	var config UpgradeConfig
	version, err := Codec.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	if version != Version {
		return nil, fmt.Errorf("Invalid version")
	}

	var PrecompileUpgrades []params.PrecompileUpgrade

	for _, precompileUpgrade := range config.PrecompileUpgrades {
		module, ok := modules.GetPrecompileModule(precompileUpgrade.StructName)
		if !ok {
			return nil, fmt.Errorf("unknown precompile config: %s", precompileUpgrade.StructName)
		}
		preCompile := module.MakeConfig()

		version, err := Codec.Unmarshal(precompileUpgrade.Bytes, preCompile)
		if version != Version {
			return nil, fmt.Errorf("Invalid version")
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

func NewUpgradeConfig(config params.UpgradeConfig) (*UpgradeConfig, error) {
	PrecompileUpgrades := make([]PrecompileUpgrade, 0)
	for _, precompileConfig := range config.PrecompileUpgrades {
		bytes, err := Codec.Marshal(Version, precompileConfig.Config)
		if err != nil {
			return nil, err
		}
		PrecompileUpgrades = append(PrecompileUpgrades, PrecompileUpgrade{
			StructName: precompileConfig.Key(),
			Bytes:      bytes,
		})
	}

	optionalNetworkUpgrades := make([]params.Fork, 0)
	if config.OptionalNetworkUpgrades != nil {
		optionalNetworkUpgrades = config.OptionalNetworkUpgrades.Updates
	}

	wrappedConfig := UpgradeConfig{
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

func (r *UpgradeConfig) Config() params.UpgradeConfig {
	return r.config
}

func (r *UpgradeConfig) Bytes() []byte {
	return r.bytes
}

func (r *UpgradeConfig) Hash() [8]byte {
	hash := sha256.Sum256(r.bytes)
	var firstBytes [8]byte
	copy(firstBytes[:], hash[:8])
	return firstBytes
}
