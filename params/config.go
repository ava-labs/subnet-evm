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
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

var (
	errNonGenesisForkByHeight = errors.New("subnet-evm only supports forking by height at the genesis block")

	SubnetEVMChainID = big.NewInt(43214)

	// For legacy tests
	MinGasPrice        int64 = 225_000_000_000
	TestInitialBaseFee int64 = 225_000_000_000
	TestMaxBaseFee           = big.NewInt(225_000_000_000)

	ExtraDataSize        = 80
	RollupWindow  uint64 = 10

	DefaultFeeConfig = commontype.FeeConfig{
		GasLimit:        big.NewInt(8_000_000),
		TargetBlockRate: 2, // in seconds

		MinBaseFee:               big.NewInt(25_000_000_000),
		TargetGas:                big.NewInt(15_000_000),
		BaseFeeChangeDenominator: big.NewInt(36),

		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		BlockGasCostStep: big.NewInt(200_000),
	}
)

var (
	// SubnetEVMDefaultChainConfig is the default configuration
	SubnetEVMDefaultChainConfig = &ChainConfig{
		ChainID:            SubnetEVMChainID,
		FeeConfig:          DefaultFeeConfig,
		AllowFeeRecipients: false,

		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),

		NetworkUpgrades: NetworkUpgrades{
			SubnetEVMTimestamp: big.NewInt(0),
		},
	}

	TestChainConfig = &ChainConfig{
		AvalancheContext:    AvalancheContext{snow.DefaultContextTest()},
		ChainID:             big.NewInt(1),
		FeeConfig:           DefaultFeeConfig,
		AllowFeeRecipients:  false,
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.Hash{},
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		NetworkUpgrades:     NetworkUpgrades{big.NewInt(0)},
		PrecompileUpgrade:   PrecompileUpgrade{},
		UpgradeConfig:       UpgradeConfig{},
	}

	TestPreSubnetEVMConfig = &ChainConfig{
		AvalancheContext:    AvalancheContext{snow.DefaultContextTest()},
		ChainID:             big.NewInt(1),
		FeeConfig:           DefaultFeeConfig,
		AllowFeeRecipients:  false,
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.Hash{},
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		NetworkUpgrades:     NetworkUpgrades{},
		PrecompileUpgrade:   PrecompileUpgrade{},
		UpgradeConfig:       UpgradeConfig{},
	}
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	AvalancheContext `json:"-"` // Avalanche specific context set during VM initialization. Not serialized.

	ChainID            *big.Int             `json:"chainId"`                      // chainId identifies the current chain and is used for replay protection
	FeeConfig          commontype.FeeConfig `json:"feeConfig"`                    // Set the configuration for the dynamic fee algorithm
	AllowFeeRecipients bool                 `json:"allowFeeRecipients,omitempty"` // Allows fees to be collected by block builders.

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)

	NetworkUpgrades              // Config for timestamps that enable avalanche network upgrades
	PrecompileUpgrade            // Config for enabling precompiles from genesis
	UpgradeConfig     `json:"-"` // Config specified in upgradeBytes (avalanche network upgrades or enable/disabling precompiles). Skip encoding/decoding directly into ChainConfig.
}

// UpgradeConfig includes the following configs that may be specified in upgradeBytes:
// - Timestamps that enable avalanche network upgrades,
// - Enabling or disabling precompiles as network upgrades.
type UpgradeConfig struct {
	// Config for blocks/timestamps that enable network upgrades.
	// Note: if NetworkUpgrades is specified in the JSON all previously activated
	// forks must be present or upgradeBytes will be rejected.
	NetworkUpgrades *NetworkUpgrades `json:"networkUpgrades,omitempty"`

	// Config for enabling and disabling precompiles as network upgrades.
	PrecompileUpgrades []PrecompileUpgrade `json:"precompileUpgrades,omitempty"`
}

// AvalancheContext provides Avalanche specific context directly into the EVM.
type AvalancheContext struct {
	SnowCtx *snow.Context
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	// convert nested data structures to json
	feeBytes, err := json.Marshal(c.FeeConfig)
	if err != nil {
		feeBytes = []byte("cannot marshal FeeConfig")
	}
	networkUpgradesBytes, err := json.Marshal(c.NetworkUpgrades)
	if err != nil {
		networkUpgradesBytes = []byte("cannot marshal NetworkUpgrades")
	}
	precompileUpgradeBytes, err := json.Marshal(c.PrecompileUpgrade)
	if err != nil {
		precompileUpgradeBytes = []byte("cannot marshal PrecompileUpgrade")
	}
	upgradeConfigBytes, err := json.Marshal(c.UpgradeConfig)
	if err != nil {
		upgradeConfigBytes = []byte("cannot marshal UpgradeConfig")
	}

	return fmt.Sprintf("{ChainID: %v Homestead: %v EIP150: %v EIP155: %v EIP158: %v Byzantium: %v Constantinople: %v Petersburg: %v Istanbul: %v, Muir Glacier: %v, Subnet EVM: %v, FeeConfig: %v, AllowFeeRecipients: %v, NetworkUpgrades: %v, PrecompileUpgrade: %v, UpgradeConfig: %v, Engine: Dummy Consensus Engine}",
		c.ChainID,
		c.HomesteadBlock,
		c.EIP150Block,
		c.EIP155Block,
		c.EIP158Block,
		c.ByzantiumBlock,
		c.ConstantinopleBlock,
		c.PetersburgBlock,
		c.IstanbulBlock,
		c.MuirGlacierBlock,
		c.SubnetEVMTimestamp,
		string(feeBytes),
		c.AllowFeeRecipients,
		string(networkUpgradesBytes),
		string(precompileUpgradeBytes),
		string(upgradeConfigBytes),
	)
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *ChainConfig) IsHomestead(num *big.Int) bool {
	return utils.IsForked(c.HomesteadBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return utils.IsForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return utils.IsForked(c.EIP155Block, num)
}

// IsEIP158 returns whether num is either equal to the EIP158 fork block or greater.
func (c *ChainConfig) IsEIP158(num *big.Int) bool {
	return utils.IsForked(c.EIP158Block, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *ChainConfig) IsByzantium(num *big.Int) bool {
	return utils.IsForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *ChainConfig) IsConstantinople(num *big.Int) bool {
	return utils.IsForked(c.ConstantinopleBlock, num)
}

// IsMuirGlacier returns whether num is either equal to the Muir Glacier (EIP-2384) fork block or greater.
func (c *ChainConfig) IsMuirGlacier(num *big.Int) bool {
	return utils.IsForked(c.MuirGlacierBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return utils.IsForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && utils.IsForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return utils.IsForked(c.IstanbulBlock, num)
}

// IsSubnetEVM returns whether [blockTimestamp] is either equal to the SubnetEVM fork block timestamp or greater.
func (c *ChainConfig) IsSubnetEVM(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.getNetworkUpgrades().SubnetEVMTimestamp, blockTimestamp)
}

// PRECOMPILE UPGRADES START HERE

// IsContractDeployerAllowList returns whether [blockTimestamp] is either equal to the ContractDeployerAllowList fork block timestamp or greater.
func (c *ChainConfig) IsContractDeployerAllowList(blockTimestamp *big.Int) bool {
	config := c.GetContractDeployerAllowListConfig(blockTimestamp)
	return config != nil && !config.Disable
}

// IsContractNativeMinter returns whether [blockTimestamp] is either equal to the NativeMinter fork block timestamp or greater.
func (c *ChainConfig) IsContractNativeMinter(blockTimestamp *big.Int) bool {
	config := c.GetContractNativeMinterConfig(blockTimestamp)
	return config != nil && !config.Disable
}

// IsTxAllowList returns whether [blockTimestamp] is either equal to the TxAllowList fork block timestamp or greater.
func (c *ChainConfig) IsTxAllowList(blockTimestamp *big.Int) bool {
	config := c.GetTxAllowListConfig(blockTimestamp)
	return config != nil && !config.Disable
}

// IsFeeConfigManager returns whether [blockTimestamp] is either equal to the FeeConfigManager fork block timestamp or greater.
func (c *ChainConfig) IsFeeConfigManager(blockTimestamp *big.Int) bool {
	config := c.GetFeeConfigManagerConfig(blockTimestamp)
	return config != nil && !config.Disable
}

// IsRewardManager returns whether [blockTimestamp] is either equal to the RewardManager fork block timestamp or greater.
func (c *ChainConfig) IsRewardManager(blockTimestamp *big.Int) bool {
	config := c.GetRewardManagerConfig(blockTimestamp)
	return config != nil && !config.Disable
}

// ADD YOUR PRECOMPILE HERE
/*
func (c *ChainConfig) Is{YourPrecompile}(blockTimestamp *big.Int) bool {
	config := c.Get{YourPrecompile}Config(blockTimestamp)
	return config != nil && !config.Disable
}
*/

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64, timestamp uint64) *ConfigCompatError {
	bNumber := new(big.Int).SetUint64(height)
	bTimestamp := new(big.Int).SetUint64(timestamp)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bNumber, bTimestamp)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bNumber.SetUint64(err.RewindTo)
	}
	return lasterr
}

// Verify verifies chain config and returns error
func (c *ChainConfig) Verify() error {
	if err := c.FeeConfig.Verify(); err != nil {
		return err
	}

	// Verify the precompile upgrades are internally consistent given the existing chainConfig.
	if err := c.verifyPrecompileUpgrades(); err != nil {
		return err
	}

	return nil
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name     string
		block    *big.Int
		optional bool // if true, the fork may be nil and next fork is still allowed
	}
	var lastFork fork
	for _, cur := range []fork{
		{name: "homesteadBlock", block: c.HomesteadBlock},
		{name: "eip150Block", block: c.EIP150Block},
		{name: "eip155Block", block: c.EIP155Block},
		{name: "eip158Block", block: c.EIP158Block},
		{name: "byzantiumBlock", block: c.ByzantiumBlock},
		{name: "constantinopleBlock", block: c.ConstantinopleBlock},
		{name: "petersburgBlock", block: c.PetersburgBlock},
		{name: "istanbulBlock", block: c.IstanbulBlock},
		{name: "muirGlacierBlock", block: c.MuirGlacierBlock, optional: true},
	} {
		if cur.block != nil && common.Big0.Cmp(cur.block) != 0 {
			return errNonGenesisForkByHeight
		}
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}

	// Note: In Avalanche, hard forks must take place via block timestamps instead
	// of block numbers since blocks are produced asynchronously. Therefore, we do not
	// check that the block timestamps in the same way as for
	// the block number forks since it would not be a meaningful comparison.
	// Instead, we check only that Phases are enabled in order.
	// Note: we do not add the optional stateful precompile configs in here because they are optional
	// and independent, such that the ordering they are enabled does not impact the correctness of the
	// chain config.
	lastFork = fork{}
	for _, cur := range []fork{
		{name: "subnetEVMTimestamp", block: c.SubnetEVMTimestamp},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}
	return nil
}

// checkCompatible confirms that [newcfg] is backwards compatible with [c] to upgrade with the given head block height and timestamp.
// This confirms that all Ethereum and Avalanche upgrades are backwards compatible as well as that the precompile config is backwards
// compatible.
func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, lastHeight *big.Int, lastTimestamp *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.HomesteadBlock, newcfg.HomesteadBlock, lastHeight) {
		return newCompatError("Homestead fork block", c.HomesteadBlock, newcfg.HomesteadBlock)
	}
	if isForkIncompatible(c.EIP150Block, newcfg.EIP150Block, lastHeight) {
		return newCompatError("EIP150 fork block", c.EIP150Block, newcfg.EIP150Block)
	}
	if isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, lastHeight) {
		return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isForkIncompatible(c.EIP158Block, newcfg.EIP158Block, lastHeight) {
		return newCompatError("EIP158 fork block", c.EIP158Block, newcfg.EIP158Block)
	}
	if c.IsEIP158(lastHeight) && !utils.BigNumEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP158 chain ID", c.EIP158Block, newcfg.EIP158Block)
	}
	if isForkIncompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, lastHeight) {
		return newCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	}
	if isForkIncompatible(c.ConstantinopleBlock, newcfg.ConstantinopleBlock, lastHeight) {
		return newCompatError("Constantinople fork block", c.ConstantinopleBlock, newcfg.ConstantinopleBlock)
	}
	if isForkIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, lastHeight) {
		// the only case where we allow Petersburg to be set in the past is if it is equal to Constantinople
		// mainly to satisfy fork ordering requirements which state that Petersburg fork be set if Constantinople fork is set
		if isForkIncompatible(c.ConstantinopleBlock, newcfg.PetersburgBlock, lastHeight) {
			return newCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
		}
	}
	if isForkIncompatible(c.IstanbulBlock, newcfg.IstanbulBlock, lastHeight) {
		return newCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if isForkIncompatible(c.MuirGlacierBlock, newcfg.MuirGlacierBlock, lastHeight) {
		return newCompatError("Muir Glacier fork block", c.MuirGlacierBlock, newcfg.MuirGlacierBlock)
	}

	// Check subnet-evm specific activations
	newNetworkUpgrades := newcfg.getNetworkUpgrades()
	if c.UpgradeConfig.NetworkUpgrades != nil && newcfg.UpgradeConfig.NetworkUpgrades == nil {
		// Note: if the current NetworkUpgrades are set via UpgradeConfig, then a new config
		// without NetworkUpgrades will be treated as having specified an empty set of network
		// upgrades (ie., treated as the user intends to cancel scheduled forks)
		newNetworkUpgrades = &NetworkUpgrades{}
	}
	if err := c.getNetworkUpgrades().CheckCompatible(newNetworkUpgrades, lastTimestamp); err != nil {
		return err
	}

	// Check that the precompiles on the new config are compatible with the existing precompile config.
	if err := c.CheckPrecompilesCompatible(newcfg.PrecompileUpgrades, lastTimestamp); err != nil {
		return err
	}

	// TODO verify that the fee config is fully compatible between [c] and [newcfg].
	return nil
}

// getNetworkUpgrades returns NetworkUpgrades from upgrade config if set there,
// otherwise it falls back to the genesis chain config.
func (c *ChainConfig) getNetworkUpgrades() *NetworkUpgrades {
	if upgradeConfigOverride := c.UpgradeConfig.NetworkUpgrades; upgradeConfigOverride != nil {
		return upgradeConfigOverride
	}
	return &c.NetworkUpgrades
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (utils.IsForked(s1, head) || utils.IsForked(s2, head)) && !utils.BigNumEqual(s1, s2)
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool

	// Rules for Avalanche releases
	IsSubnetEVM bool

	// Optional stateful precompile rules
	IsContractDeployerAllowListEnabled bool
	IsContractNativeMinterEnabled      bool
	IsTxAllowListEnabled               bool
	IsFeeConfigManagerEnabled          bool
	IsRewardManagerEnabled             bool
	// ADD YOUR PRECOMPILE HERE
	// Is{YourPrecompile}Enabled         bool

	// Precompiles maps addresses to stateful precompiled contracts that are enabled
	// for this rule set.
	// Note: none of these addresses should conflict with the address space used by
	// any existing precompiles.
	Precompiles map[common.Address]precompile.StatefulPrecompiledContract
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:          new(big.Int).Set(chainID),
		IsHomestead:      c.IsHomestead(num),
		IsEIP150:         c.IsEIP150(num),
		IsEIP155:         c.IsEIP155(num),
		IsEIP158:         c.IsEIP158(num),
		IsByzantium:      c.IsByzantium(num),
		IsConstantinople: c.IsConstantinople(num),
		IsPetersburg:     c.IsPetersburg(num),
		IsIstanbul:       c.IsIstanbul(num),
	}
}

// AvalancheRules returns the Avalanche modified rules to support Avalanche
// network upgrades
func (c *ChainConfig) AvalancheRules(blockNum, blockTimestamp *big.Int) Rules {
	rules := c.rules(blockNum)

	rules.IsSubnetEVM = c.IsSubnetEVM(blockTimestamp)
	rules.IsContractDeployerAllowListEnabled = c.IsContractDeployerAllowList(blockTimestamp)
	rules.IsContractNativeMinterEnabled = c.IsContractNativeMinter(blockTimestamp)
	rules.IsTxAllowListEnabled = c.IsTxAllowList(blockTimestamp)
	rules.IsFeeConfigManagerEnabled = c.IsFeeConfigManager(blockTimestamp)
	rules.IsRewardManagerEnabled = c.IsRewardManager(blockTimestamp)
	// ADD YOUR PRECOMPILE HERE
	// rules.Is{YourPrecompile}Enabled = c.{IsYourPrecompile}(blockTimestamp)

	// Initialize the stateful precompiles that should be enabled at [blockTimestamp].
	rules.Precompiles = make(map[common.Address]precompile.StatefulPrecompiledContract)
	for _, config := range c.EnabledStatefulPrecompiles(blockTimestamp) {
		if config.IsDisabled() {
			continue
		}
		rules.Precompiles[config.Address()] = config.Contract()
	}

	return rules
}

// GetFeeConfig returns the original FeeConfig contained in the genesis ChainConfig.
// Implements precompile.ChainConfig interface.
func (c *ChainConfig) GetFeeConfig() commontype.FeeConfig {
	return c.FeeConfig
}

// AllowedFeeRecipients returns the original AllowedFeeRecipients parameter contained in the genesis ChainConfig.
// Implements precompile.ChainConfig interface.
func (c *ChainConfig) AllowedFeeRecipients() bool {
	return c.AllowFeeRecipients
}
