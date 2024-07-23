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
	"bytes"
	_ "embed"
	"math/big"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/dummy"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/precompile/allowlist"
	"github.com/ethereum/go-ethereum/precompile/contracts/deployerallowlist"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/ethereum/go-ethereum/triedb/pathdb"
	"github.com/ethereum/go-ethereum/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGenesisBlock(db ethdb.Database, triedb *triedb.Database, genesis *Genesis, lastAcceptedHash common.Hash) (*params.ChainConfig, common.Hash, error) {
	return SetupGenesisBlock(db, triedb, genesis, lastAcceptedHash, false)
}

func TestGenesisBlockForTesting(t *testing.T) {
	genesisBlockForTestingHash := common.HexToHash("0x114ce61b50051f70768f982f7b59e82dd73b7bbd768e310c9d9f508d492e687b")
	block := GenesisBlockForTesting(rawdb.NewMemoryDatabase(), common.Address{1}, big.NewInt(1))
	if block.Hash() != genesisBlockForTestingHash {
		t.Errorf("wrong testing genesis hash, got %v, want %v", block.Hash(), genesisBlockForTestingHash)
	}
}

func TestSetupGenesis(t *testing.T) {
	testSetupGenesis(t, rawdb.HashScheme)
	testSetupGenesis(t, rawdb.PathScheme)
}

func testSetupGenesis(t *testing.T, scheme string) {
	preSubnetConfig := *params.TestPreSubnetEVMConfig
	preSubnetConfig.SubnetEVMTimestamp = utils.NewUint64(100)
	var (
		customghash = common.HexToHash("0x4a12fe7bf8d40d152d7e9de22337b115186a4662aa3a97217b36146202bbfc66")
		customg     = Genesis{
			Config: &preSubnetConfig,
			Alloc: types.GenesisAlloc{
				{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
			},
			GasLimit: preSubnetConfig.FeeConfig.GasLimit.Uint64(),
		}
		oldcustomg = customg
	)

	rollbackpreSubnetConfig := preSubnetConfig
	rollbackpreSubnetConfig.SubnetEVMTimestamp = utils.NewUint64(90)
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
				return setupGenesisBlock(db, triedb.NewDatabase(db, newDbConfig(scheme)), new(Genesis), common.Hash{})
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: nil,
		},
		{
			name: "no block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				return setupGenesisBlock(db, triedb.NewDatabase(db, newDbConfig(scheme)), nil, common.Hash{})
			},
			wantErr:    ErrNoGenesis,
			wantConfig: nil,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				tdb := triedb.NewDatabase(db, newDbConfig(scheme))
				customg.Commit(db, tdb)
				return setupGenesisBlock(db, tdb, nil, common.Hash{})
			},
			wantErr:    ErrNoGenesis,
			wantConfig: nil,
		},
		{
			name: "compatible config in DB",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				tdb := triedb.NewDatabase(db, newDbConfig(scheme))
				oldcustomg.Commit(db, tdb)
				return setupGenesisBlock(db, tdb, &customg, customghash)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "incompatible config for avalanche fork in DB",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with SubnetEVM transition at 90.
				// Advance to block #4, past the SubnetEVM transition block of customg.
				tdb := triedb.NewDatabase(db, newDbConfig(scheme))
				genesis, err := oldcustomg.Commit(db, tdb)
				if err != nil {
					t.Fatal(err)
				}

				bc, _ := NewBlockChain(db, DefaultCacheConfigWithScheme(scheme), &oldcustomg, dummy.NewFullFaker(), vm.Config{}, genesis.Hash(), false)
				defer bc.Stop()

				_, blocks, _, err := GenerateChainWithGenesis(&oldcustomg, dummy.NewFullFaker(), 4, 25, nil)
				if err != nil {
					t.Fatal(err)
				}
				bc.InsertChain(blocks)

				for _, block := range blocks {
					if err := bc.Accept(block); err != nil {
						t.Fatal(err)
					}
				}

				// This should return a compatibility error.
				return setupGenesisBlock(db, tdb, &customg, bc.lastAccepted.Hash())
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
			wantErr: &params.ConfigCompatError{
				What:         "SubnetEVM fork block timestamp",
				StoredTime:   u64(90),
				NewTime:      u64(100),
				RewindToTime: 89,
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
					deployerallowlist.ConfigKey: deployerallowlist.NewConfig(utils.NewUint64(0), []common.Address{addr}, nil, nil),
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
				Alloc: types.GenesisAlloc{
					{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
				},
				GasLimit: config.FeeConfig.GasLimit.Uint64(),
			}

			db := rawdb.NewMemoryDatabase()

			genesisBlock := genesis.ToBlock()
			genesisRoot := genesisBlock.Root()

			_, _, err := setupGenesisBlock(db, triedb.NewDatabase(db, triedb.HashDefaults), genesis, genesisBlock.Hash())
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
		Config: params.TestChainConfig,
		Alloc: types.GenesisAlloc{
			{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
		},
		GasLimit: params.TestChainConfig.FeeConfig.GasLimit.Uint64(),
	}
	bc, _ := NewBlockChain(db, DefaultCacheConfig, &customg, dummy.NewFullFaker(), vm.Config{}, common.Hash{}, false)
	defer bc.Stop()

	// Advance header to block #4, past the ContractDeployerAllowListConfig.
	_, blocks, _, _ := GenerateChainWithGenesis(&customg, dummy.NewFullFaker(), 4, 25, nil)

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
	require.Greater(block.Time, bc.lastAccepted.Time())

	activatedGenesisConfig := *customg.Config
	contractDeployerConfig := deployerallowlist.NewConfig(utils.NewUint64(51), nil, nil, nil)
	activatedGenesisConfig.UpgradeConfig.PrecompileUpgrades = []params.PrecompileUpgrade{
		{
			Config: contractDeployerConfig,
		},
	}
	customg.Config = &activatedGenesisConfig

	// assert block is after the activation block
	require.Greater(block.Time, *contractDeployerConfig.Timestamp())
	// assert last accepted block is before the activation block
	require.Less(bc.lastAccepted.Time(), *contractDeployerConfig.Timestamp())

	// This should not return any error since the last accepted block is before the activation block.
	config, _, err := setupGenesisBlock(db, triedb.NewDatabase(db, nil), &customg, bc.lastAccepted.Hash())
	require.NoError(err)
	if !reflect.DeepEqual(config, customg.Config) {
		t.Errorf("returned %v\nwant     %v", config, customg.Config)
	}
}

func TestGenesisWriteUpgradesRegression(t *testing.T) {
	require := require.New(t)
	config := *params.TestChainConfig
	genesis := &Genesis{
		Config: &config,
		Alloc: types.GenesisAlloc{
			{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
		},
		GasLimit: config.FeeConfig.GasLimit.Uint64(),
	}

	db := rawdb.NewMemoryDatabase()
	trieDB := triedb.NewDatabase(db, triedb.HashDefaults)
	genesisBlock := genesis.MustCommit(db, trieDB)

	_, _, err := SetupGenesisBlock(db, trieDB, genesis, genesisBlock.Hash(), false)
	require.NoError(err)

	genesis.Config.UpgradeConfig.PrecompileUpgrades = []params.PrecompileUpgrade{
		{
			Config: deployerallowlist.NewConfig(utils.NewUint64(51), nil, nil, nil),
		},
	}
	_, _, err = SetupGenesisBlock(db, trieDB, genesis, genesisBlock.Hash(), false)
	require.NoError(err)

	timestamp := uint64(100)
	lastAcceptedBlock := types.NewBlock(&types.Header{
		ParentHash: common.Hash{1, 2, 3},
		Number:     big.NewInt(100),
		GasLimit:   8_000_000,
		Extra:      nil,
		Time:       timestamp,
	}, nil, nil, nil, trie.NewStackTrie(nil))
	rawdb.WriteBlock(db, lastAcceptedBlock)

	// Attempt restart after the chain has advanced past the activation of the precompile upgrade.
	// This tests a regression where the UpgradeConfig would not be written to disk correctly.
	_, _, err = SetupGenesisBlock(db, trieDB, genesis, lastAcceptedBlock.Hash(), false)
	require.NoError(err)
}

func newDbConfig(scheme string) *triedb.Config {
	if scheme == rawdb.HashScheme {
		return triedb.HashDefaults
	}
	return &triedb.Config{PathDB: pathdb.Defaults}
}

func TestVerkleGenesisCommit(t *testing.T) {
	var verkleTime uint64 = 0
	verkleConfig := &params.ChainConfig{
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
		CancunTime:          &verkleTime,
		VerkleTime:          &verkleTime,
	}

	genesis := &Genesis{
		BaseFee:    big.NewInt(params.TestInitialBaseFee),
		Config:     verkleConfig,
		Timestamp:  verkleTime,
		Difficulty: big.NewInt(0),
		Alloc: types.GenesisAlloc{
			{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
		},
	}

	expected := common.Hex2Bytes("14398d42be3394ff8d50681816a4b7bf8d8283306f577faba2d5bc57498de23b")
	got := genesis.ToBlock().Root().Bytes()
	if !bytes.Equal(got, expected) {
		t.Fatalf("invalid genesis state root, expected %x, got %x", expected, got)
	}

	db := rawdb.NewMemoryDatabase()
	triedb := triedb.NewDatabase(db, &triedb.Config{IsVerkle: true, PathDB: pathdb.Defaults})
	block := genesis.MustCommit(db, triedb)
	if !bytes.Equal(block.Root().Bytes(), expected) {
		t.Fatalf("invalid genesis state root, expected %x, got %x", expected, got)
	}

	// Test that the trie is verkle
	if !triedb.IsVerkle() {
		t.Fatalf("expected trie to be verkle")
	}

	if !rawdb.ExistsAccountTrieNode(db, nil) {
		t.Fatal("could not find node")
	}
}
