// (c) 2019-2021, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2017 The go-ethereum Authors
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

package core

import (
	_ "embed"
	"math/big"
	"reflect"
	"testing"

	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGenesisBlock(db ethdb.Database, genesis *Genesis, lastAcceptedHash common.Hash) (*params.ChainConfig, common.Hash, error) {
	conf, err := SetupGenesisBlock(db, genesis, lastAcceptedHash, false)
	stored := rawdb.ReadCanonicalHash(db, 0)
	return conf, stored, err
}

func TestGenesisBlockForTesting(t *testing.T) {
	genesisBlockForTestingHash := common.HexToHash("0x114ce61b50051f70768f982f7b59e82dd73b7bbd768e310c9d9f508d492e687b")
	block := GenesisBlockForTesting(rawdb.NewMemoryDatabase(), common.Address{1}, big.NewInt(1))
	if block.Hash() != genesisBlockForTestingHash {
		t.Errorf("wrong testing genesis hash, got %v, want %v", block.Hash(), genesisBlockForTestingHash)
	}
}

func TestSetupGenesis(t *testing.T) {
	preSubnetConfig := *params.TestPreSubnetEVMConfig
	preSubnetConfig.SubnetEVMTimestamp = big.NewInt(100)
	var (
		customghash = common.HexToHash("0x4a12fe7bf8d40d152d7e9de22337b115186a4662aa3a97217b36146202bbfc66")
		customg     = Genesis{
			Config: &preSubnetConfig,
			Alloc: GenesisAlloc{
				{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
			},
			GasLimit: preSubnetConfig.FeeConfig.GasLimit.Uint64(),
		}
		oldcustomg = customg
	)

	rollbackpreSubnetConfig := preSubnetConfig
	rollbackpreSubnetConfig.SubnetEVMTimestamp = big.NewInt(90)
	oldcustomg.Config = &rollbackpreSubnetConfig
	tests := []struct {
		name       string
		fn         func(ethdb.Database) (*params.ChainConfig, common.Hash, error)
		wantConfig *params.ChainConfig
		wantHash   common.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				return setupGenesisBlock(db, new(Genesis), common.Hash{})
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: nil,
		},
		{
			name: "no block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				return setupGenesisBlock(db, nil, common.Hash{})
			},
			wantErr:    ErrNoGenesis,
			wantConfig: nil,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				customg.MustCommit(db)
				return setupGenesisBlock(db, nil, common.Hash{})
			},
			wantErr:    ErrNoGenesis,
			wantHash:   customghash,
			wantConfig: nil,
		},
		{
			name: "compatible config in DB",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				oldcustomg.MustCommit(db)
				return setupGenesisBlock(db, &customg, customghash)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "incompatible config for avalanche fork in DB",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with SubnetEVM transition at 90.
				// Advance to block #4, past the SubnetEVM transition block of customg.
				genesis := oldcustomg.MustCommit(db)

				bc, _ := NewBlockChain(db, DefaultCacheConfig, oldcustomg.Config, dummy.NewFullFaker(), vm.Config{}, common.Hash{})
				defer bc.Stop()

				blocks, _, _ := GenerateChain(oldcustomg.Config, genesis, dummy.NewFullFaker(), db, 4, 25, nil)
				bc.InsertChain(blocks)

				for _, block := range blocks {
					if err := bc.Accept(block); err != nil {
						t.Fatal(err)
					}
				}

				// This should return a compatibility error.
				return setupGenesisBlock(db, &customg, bc.lastAccepted.Hash())
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
			wantErr: &params.ConfigCompatError{
				What:         "SubnetEVM fork block timestamp",
				StoredConfig: big.NewInt(90),
				NewConfig:    big.NewInt(100),
				RewindTo:     89,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			config, hash, err := test.fn(db)
			// Check the return values.
			if !reflect.DeepEqual(err, test.wantErr) {
				spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
				t.Errorf("returned error %#v, want %#v", spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
			}
			if !reflect.DeepEqual(config, test.wantConfig) {
				t.Errorf("returned %v\nwant     %v", config, test.wantConfig)
			}
			if hash != test.wantHash {
				t.Errorf("returned hash %s, want %s", hash.Hex(), test.wantHash.Hex())
			} else if err == nil {
				// Check database content.
				stored := rawdb.ReadBlock(db, test.wantHash, 0)
				if stored.Hash() != test.wantHash {
					t.Errorf("block in DB has hash %s, want %s", stored.Hash(), test.wantHash)
				}
			}
		})
	}
}

func TestStatefulPrecompilesConfigure(t *testing.T) {
	type test struct {
		getConfig   func() *params.ChainConfig             // Return the config that enables the stateful precompile at the genesis for the test
		assertState func(t *testing.T, sdb *state.StateDB) // Check that the stateful precompiles were configured correctly
	}

	addr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")

	// Test suite to ensure that stateful precompiles are configured correctly in the genesis.
	for name, test := range map[string]test{
		"allow list enabled in genesis": {
			getConfig: func() *params.ChainConfig {
				config := *params.TestChainConfig
				config.GenesisPrecompiles = params.Precompiles{
					deployerallowlist.ConfigKey: deployerallowlist.NewConfig(big.NewInt(0), []common.Address{addr}, nil),
				}
				return &config
			},
			assertState: func(t *testing.T, sdb *state.StateDB) {
				assert.Equal(t, allowlist.AdminRole, deployerallowlist.GetContractDeployerAllowListStatus(sdb, addr), "unexpected allow list status for modified address")
				assert.Equal(t, uint64(1), sdb.GetNonce(deployerallowlist.ContractAddress))
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			config := test.getConfig()

			genesis := &Genesis{
				Config: config,
				Alloc: GenesisAlloc{
					{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
				},
				GasLimit: config.FeeConfig.GasLimit.Uint64(),
			}

			db := rawdb.NewMemoryDatabase()

			genesisBlock := genesis.ToBlock(nil)
			genesisRoot := genesisBlock.Root()

			_, err := SetupGenesisBlock(db, genesis, genesisBlock.Hash(), false)
			if err != nil {
				t.Fatal(err)
			}

			statedb, err := state.New(genesisRoot, state.NewDatabase(db), nil)
			if err != nil {
				t.Fatal(err)
			}

			if test.assertState != nil {
				test.assertState(t, statedb)
			}
		})
	}
}

// regression test for precompile activation after header block
func TestPrecompileActivationAfterHeaderBlock(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	customg := Genesis{
		Config: params.SubnetEVMDefaultChainConfig,
		Alloc: GenesisAlloc{
			{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
		},
		GasLimit: params.SubnetEVMDefaultChainConfig.FeeConfig.GasLimit.Uint64(),
	}
	genesis := customg.MustCommit(db)
	bc, _ := NewBlockChain(db, DefaultCacheConfig, customg.Config, dummy.NewFullFaker(), vm.Config{}, common.Hash{})
	defer bc.Stop()

	// Advance header to block #4, past the ContractDeployerAllowListConfig.
	blocks, _, _ := GenerateChain(customg.Config, genesis, dummy.NewFullFaker(), db, 4, 25, nil)

	require := require.New(t)
	_, err := bc.InsertChain(blocks)
	require.NoError(err)

	// accept up to block #2
	for _, block := range blocks[:2] {
		require.NoError(bc.Accept(block))
	}
	block := bc.CurrentBlock()

	require.Equal(blocks[1].Hash(), bc.lastAccepted.Hash())
	// header must be bigger than last accepted
	require.Greater(block.Time(), bc.lastAccepted.Time())

	activatedGenesis := customg
	contractDeployerConfig := deployerallowlist.NewConfig(big.NewInt(51), nil, nil)
	activatedGenesis.Config.UpgradeConfig.PrecompileUpgrades = []params.PrecompileUpgrade{
		{
			Config: contractDeployerConfig,
		},
	}

	// assert block is after the activation block
	require.Greater(block.Time(), contractDeployerConfig.Timestamp().Uint64())
	// assert last accepted block is before the activation block
	require.Less(bc.lastAccepted.Time(), contractDeployerConfig.Timestamp().Uint64())

	// This should not return any error since the last accepted block is before the activation block.
	config, _, err := setupGenesisBlock(db, &activatedGenesis, bc.lastAccepted.Hash())
	require.NoError(err)
	if !reflect.DeepEqual(config, activatedGenesis.Config) {
		t.Errorf("returned %v\nwant     %v", config, activatedGenesis.Config)
	}
}
