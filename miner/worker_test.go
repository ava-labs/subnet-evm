// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package miner

import (
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils/predicate"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mockConfigurator struct {
	contract.Configurator
}

func (m *mockConfigurator) Configure(precompileconfig.ChainConfig, precompileconfig.Config, contract.StateDB, contract.ConfigurationBlockContext) error {
	return nil
}

func TestEnvionmentPredicateResults(t *testing.T) {
	var (
		require        = require.New(t)
		ctrl           = gomock.NewController(t)
		db             = rawdb.NewMemoryDatabase()
		key1, _        = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1          = crypto.PubkeyToAddress(key1.PublicKey)
		genesisBalance = big.NewInt(1 * params.Ether)
		config         = *params.TestChainConfig
	)

	// Setup a pretend precompile module and its config
	type mockPrecompileConfig struct {
		precompileconfig.Predicater
		precompileconfig.Config
	}
	mockConfig := precompileconfig.NewMockConfig(ctrl)
	mockConfig.EXPECT().Timestamp().Return(new(uint64)).AnyTimes()
	mockConfig.EXPECT().IsDisabled().Return(false).AnyTimes()
	mockPredicater := precompileconfig.NewMockPredicater(ctrl)
	module := modules.Module{
		ConfigKey:    "mock",
		Address:      common.HexToAddress("0x02000000000000000000000000000000000000ff"),
		Configurator: &mockConfigurator{},
	}
	require.NoError(modules.RegisterModule(module))

	// Add the precompile to the genesis config
	config.GenesisPrecompiles = map[string]precompileconfig.Config{
		module.ConfigKey: &mockPrecompileConfig{
			Predicater: mockPredicater,
			Config:     mockConfig,
		},
	}

	// Setup genesis and the chain
	genesis := &core.Genesis{
		Config: &config,
		Alloc:  core.GenesisAlloc{addr1: {Balance: genesisBalance}},
	}
	chain, err := core.NewBlockChain(
		db,
		core.DefaultCacheConfig,
		genesis,
		dummy.NewCoinbaseFaker(),
		vm.Config{},
		common.Hash{},
		false,
	)
	require.NoError(err)
	defer chain.Stop()

	parent := chain.CurrentBlock()
	predicateResults := []byte{1, 2, 3}

	// Test that the predicate results are stored in the environment,
	// but not if an error occurs in [core.ApplyTransaction]. Here we use
	// a gas limit of 0 to force an error in one test case, and a normal
	// gas limit in the other to test the default behavior.
	type test struct {
		gasLimit                 uint64
		expectedErr              error
		expectedPredicateResults []byte
		expectedTxCount          int
	}
	for _, tt := range []test{
		{gasLimit: 0, expectedErr: core.ErrGasLimitReached},
		{gasLimit: parent.GasLimit, expectedPredicateResults: predicateResults, expectedTxCount: 1},
	} {
		header := &types.Header{
			ParentHash: parent.Hash(),
			Number:     new(big.Int).Add(parent.Number, common.Big1),
			GasLimit:   tt.gasLimit,
			Extra:      nil,
			Time:       parent.Time + 1,
			Difficulty: big.NewInt(1),
			BaseFee:    params.DefaultFeeConfig.MinBaseFee,
		}
		chainID := chain.Config().ChainID

		to := common.Address{}
		tx, err := types.SignTx(
			predicate.NewPredicateTx(
				chainID,
				0,
				&to,
				1_000_000,
				big.NewInt(225*params.GWei),
				big.NewInt(params.GWei),
				common.Big0,
				nil,
				types.AccessList{},
				module.Address,
				nil,
			),
			types.LatestSignerForChainID(chainID),
			key1,
		)
		require.NoError(err)

		// Setup the mock predicate for gas and verification results
		mockPredicater.EXPECT().PredicateGas(gomock.Any()).Return(uint64(0), nil).AnyTimes()
		mockPredicater.EXPECT().VerifyPredicate(gomock.Any(), gomock.Any()).Return(predicateResults)

		w := &worker{
			chain:       chain,
			chainConfig: genesis.Config,
		}
		env, err := w.createCurrentEnvironment(nil, parent, header, time.Now())
		require.NoError(err)
		require.NotNil(env)
		_, err = w.commitTransaction(env, tx, w.coinbase)
		require.ErrorIs(err, tt.expectedErr)
		require.Len(env.txs, tt.expectedTxCount)
		require.Equal(
			tt.expectedPredicateResults,
			env.predicateResults.GetPredicateResults(tx.Hash(), module.Address),
		)
	}
}
