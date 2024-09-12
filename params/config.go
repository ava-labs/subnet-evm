// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/upgrade"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	gethparams "github.com/ethereum/go-ethereum/params"
)

var (
	_ = do_init() // XXX: is (temporarily) here because type registeration must proceed the call to .Rules()
	// SubnetEVMDefaultConfig is the default configuration
	// without any network upgrades.
	SubnetEVMDefaultChainConfig = WithExtra(
		&ChainConfig{
			ChainID: DefaultChainID,

			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		},
		&ChainConfigExtra{
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades:    getDefaultNetworkUpgrades(upgrade.GetConfig(constants.MainnetID)), // This can be changed to correct network (local, test) via VM.
			GenesisPrecompiles: Precompiles{},
		},
	)

	TestChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
			ShanghaiTime:        utils.TimeToNewUint64(upgrade.GetConfig(constants.UnitTestID).DurangoTime),
			CancunTime:          utils.TimeToNewUint64(upgrade.GetConfig(constants.UnitTestID).EtnaTime),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades:    getDefaultNetworkUpgrades(upgrade.GetConfig(constants.UnitTestID)), // This can be changed to correct network (local, test) via VM.
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestPreSubnetEVMChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
				DurangoTimestamp:   utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
				EtnaTimestamp:      utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestSubnetEVMChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(0),
				DurangoTimestamp:   utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
				EtnaTimestamp:      utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestDurangoChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(0),
				DurangoTimestamp:   utils.NewUint64(0),
				EtnaTimestamp:      utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestEtnaChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(0),
				DurangoTimestamp:   utils.NewUint64(0),
				EtnaTimestamp:      utils.NewUint64(0),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)
	TestRules = TestChainConfig.Rules(new(big.Int), IsMergeTODO, 0)
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig = gethparams.ChainConfig

func GetExtra(c *ChainConfig) *ChainConfigExtra {
	if extra := FromChainConfig(c); extra != nil {
		return extra
	}
	return &ChainConfigExtra{}
}

func GetRulesExtra(r Rules) *RulesExtra {
	extra := FromRules(&r)
	return &extra
}

func Copy(c *ChainConfig) ChainConfig {
	cpy := *c
	extraCpy := *GetExtra(c)
	return *WithExtra(&cpy, &extraCpy)
}

type ChainConfigExtra struct {
	NetworkUpgrades // Config for timestamps that enable network upgrades. Skip encoding/decoding directly into ChainConfig.

	AvalancheContext `json:"-"` // Avalanche specific context set during VM initialization. Not serialized.

	FeeConfig          commontype.FeeConfig `json:"feeConfig"`                    // Set the configuration for the dynamic fee algorithm
	AllowFeeRecipients bool                 `json:"allowFeeRecipients,omitempty"` // Allows fees to be collected by block builders.

	GenesisPrecompiles Precompiles `json:"-"` // Config for enabling precompiles from genesis. JSON encode/decode will be handled by the custom marshaler/unmarshaler.
	UpgradeConfig      `json:"-"`  // Config specified in upgradeBytes (avalanche network upgrades or enable/disabling precompiles). Skip encoding/decoding directly into ChainConfig.
}

func (c *ChainConfigExtra) Description() string {
	// Handle nil receiver
	if c == nil {
		return ""
	}

	var banner string

	banner += "Avalanche Upgrades (timestamp based):\n"
	banner += c.NetworkUpgrades.Description()
	banner += "\n"

	precompileUpgradeBytes, err := json.Marshal(c.GenesisPrecompiles)
	if err != nil {
		precompileUpgradeBytes = []byte("cannot marshal PrecompileUpgrade")
	}
	banner += fmt.Sprintf("Precompile Upgrades: %s", string(precompileUpgradeBytes))
	banner += "\n"

	upgradeConfigBytes, err := json.Marshal(c.UpgradeConfig)
	if err != nil {
		upgradeConfigBytes = []byte("cannot marshal UpgradeConfig")
	}
	banner += fmt.Sprintf("Upgrade Config: %s", string(upgradeConfigBytes))
	banner += "\n"

	feeBytes, err := json.Marshal(c.FeeConfig)
	if err != nil {
		feeBytes = []byte("cannot marshal FeeConfig")
	}
	banner += fmt.Sprintf("Fee Config: %s", string(feeBytes))
	banner += "\n"

	banner += fmt.Sprintf("Allow Fee Recipients: %v", c.AllowFeeRecipients)
	banner += "\n"
	return banner
}

type fork struct {
	name      string
	block     *big.Int // some go-ethereum forks use block numbers
	timestamp *uint64  // Avalanche forks use timestamps
	optional  bool     // if true, the fork may be nil and next fork is still allowed
}

func (c *ChainConfigExtra) CheckConfigForkOrder() error {
	// Handle nil receiver
	if c == nil {
		return nil
	}

	// Note: In Avalanche, hard forks must take place via block timestamps instead
	// of block numbers since blocks are produced asynchronously. Therefore, we do not
	// check that the block timestamps in the same way as for
	// the block number forks since it would not be a meaningful comparison.
	// Instead, we check only that Phases are enabled in order.
	// Note: we do not add the optional stateful precompile configs in here because they are optional
	// and independent, such that the ordering they are enabled does not impact the correctness of the
	// chain config.
	if err := checkForks(c.forkOrder(), false); err != nil {
		return err
	}

	return nil
}

// checkForks checks that forks are enabled in order and returns an error if not
// [blockFork] is true if the fork is a block number fork, false if it is a timestamp fork
func checkForks(forks []fork, blockFork bool) error {
	lastFork := fork{}
	for _, cur := range forks {
		if blockFork && cur.block != nil && common.Big0.Cmp(cur.block) != 0 {
			return errNonGenesisForkByHeight
		}
		if lastFork.name != "" {
			switch {
			// Non-optional forks must all be present in the chain config up to the last defined fork
			case lastFork.block == nil && lastFork.timestamp == nil && (cur.block != nil || cur.timestamp != nil):
				if cur.block != nil {
					return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at block %v",
						lastFork.name, cur.name, cur.block)
				} else {
					return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at timestamp %v",
						lastFork.name, cur.name, *cur.timestamp)
				}

			// Fork (whether defined by block or timestamp) must follow the fork definition sequence
			case (lastFork.block != nil && cur.block != nil) || (lastFork.timestamp != nil && cur.timestamp != nil):
				if lastFork.block != nil && lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at block %v, but %v enabled at block %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				} else if lastFork.timestamp != nil && *lastFork.timestamp > *cur.timestamp {
					return fmt.Errorf("unsupported fork ordering: %v enabled at timestamp %v, but %v enabled at timestamp %v",
						lastFork.name, *lastFork.timestamp, cur.name, *cur.timestamp)
				}

				// Timestamp based forks can follow block based ones, but not the other way around
				if lastFork.timestamp != nil && cur.block != nil {
					return fmt.Errorf("unsupported fork ordering: %v used timestamp ordering, but %v reverted to block ordering",
						lastFork.name, cur.name)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || (cur.block != nil || cur.timestamp != nil) {
			lastFork = cur
		}
	}
	return nil
}

func (c *ChainConfigExtra) CheckCompatible(newcfg_ *ChainConfig, headNumber *big.Int, headTimestamp uint64) *ConfigCompatError {
	if c == nil {
		return nil
	}
	newcfg := GetExtra(newcfg_)

	// Check avalanche network upgrades
	if err := c.checkNetworkUpgradesCompatible(&newcfg.NetworkUpgrades, headTimestamp); err != nil {
		return err
	}

	// Check that the precompiles on the new config are compatible with the existing precompile config.
	if err := c.checkPrecompilesCompatible(newcfg.PrecompileUpgrades, headTimestamp); err != nil {
		return err
	}

	// Check that the state upgrades on the new config are compatible with the existing state upgrade config.
	if err := c.checkStateUpgradesCompatible(newcfg.StateUpgrades, headTimestamp); err != nil {
		return err
	}

	// TODO verify that the fee config is fully compatible between [c] and [newcfg].
	return nil
}

// isForkTimestampIncompatible returns true if a fork scheduled at timestamp s1
// cannot be rescheduled to timestamp s2 because head is already past the fork.
func isForkTimestampIncompatible(s1, s2 *uint64, head uint64) bool {
	return (isTimestampForked(s1, head) || isTimestampForked(s2, head)) && !configTimestampEqual(s1, s2)
}

// isTimestampForked returns whether a fork scheduled at timestamp s is active
// at the given head timestamp. Whilst this method is the same as isBlockForked,
// they are explicitly separate for clearer reading.
func isTimestampForked(s *uint64, head uint64) bool {
	if s == nil {
		return false
	}
	return *s <= head
}

func configTimestampEqual(x, y *uint64) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return *x == *y
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError = gethparams.ConfigCompatError

func newTimestampCompatError(what string, storedtime, newtime *uint64) *ConfigCompatError {
	var rew *uint64
	switch {
	case storedtime == nil:
		rew = newtime
	case newtime == nil || *storedtime < *newtime:
		rew = storedtime
	default:
		rew = newtime
	}
	err := &ConfigCompatError{
		What:         what,
		StoredTime:   storedtime,
		NewTime:      newtime,
		RewindToTime: 0,
	}
	if rew != nil && *rew > 0 {
		err.RewindToTime = *rew - 1
	}
	return err
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules = gethparams.Rules

type RulesExtra struct {
	gethrules gethparams.Rules

	// Rules for Avalanche releases
	AvalancheRules

	// ActivePrecompiles maps addresses to stateful precompiled contracts that are enabled
	// for this rule set.
	// Note: none of these addresses should conflict with the address space used by
	// any existing precompiles.
	ActivePrecompiles map[common.Address]precompileconfig.Config
	// Predicaters maps addresses to stateful precompile Predicaters
	// that are enabled for this rule set.
	Predicaters map[common.Address]precompileconfig.Predicater
	// AccepterPrecompiles map addresses to stateful precompile accepter functions
	// that are enabled for this rule set.
	AccepterPrecompiles map[common.Address]precompileconfig.Accepter

	gethparams.NOOPHooks // XXX: Embedded to ensure that Rules implements params.RulesHooks
}
