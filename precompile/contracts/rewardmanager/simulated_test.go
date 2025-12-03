// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/plugin/evm/customtypes"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/allowlist/allowlisttest"
	"github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/testutils"
	"github.com/ava-labs/subnet-evm/utils"

	sim "github.com/ava-labs/subnet-evm/ethclient/simulated"
	rewardmanagerbindings "github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager/rewardmanagertest/bindings"
)

var (
	adminKey, _        = crypto.GenerateKey()
	unprivilegedKey, _ = crypto.GenerateKey()

	adminAddress        = crypto.PubkeyToAddress(adminKey.PublicKey)
	unprivilegedAddress = crypto.PubkeyToAddress(unprivilegedKey.PublicKey)
)

func TestMain(m *testing.M) {
	// Ensure libevm extras are registered for tests.
	core.RegisterExtras()
	customtypes.Register()
	params.RegisterExtras()
	m.Run()
}

func newBackendWithRewardManager(t *testing.T) *sim.Backend {
	t.Helper()
	chainCfg := params.Copy(params.TestChainConfig)
	// Enable RewardManager at genesis with admin set to adminAddress.
	params.GetExtra(&chainCfg).GenesisPrecompiles = extras.Precompiles{
		rewardmanager.ConfigKey: rewardmanager.NewConfig(utils.NewUint64(0), []common.Address{adminAddress}, nil, nil, nil),
	}
	return sim.NewBackend(
		types.GenesisAlloc{
			adminAddress:        {Balance: big.NewInt(1000000000000000000)},
			unprivilegedAddress: {Balance: big.NewInt(1000000000000000000)},
		},
		sim.WithChainConfig(&chainCfg),
	)
}

// Helper functions

func deployRewardManagerTest(t *testing.T, b *sim.Backend, auth *bind.TransactOpts) (common.Address, *rewardmanagerbindings.RewardManagerTest) {
	t.Helper()
	addr, tx, contract, err := rewardmanagerbindings.DeployRewardManagerTest(auth, b.Client(), rewardmanager.ContractAddress)
	require.NoError(t, err)
	testutils.WaitReceiptSuccessful(t, b, tx)
	return addr, contract
}

func TestRewardManager(t *testing.T) {
	chainID := params.TestChainConfig.ChainID
	admin := testutils.NewAuth(t, adminKey, chainID)

	type testCase struct {
		name string
		test func(t *testing.T, backend *sim.Backend, rewardManagerIntf *rewardmanagerbindings.IRewardManager)
	}

	testCases := []testCase{
		{
			name: "should verify sender is admin",
			test: func(t *testing.T, _ *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				allowlisttest.VerifyRole(t, rewardManager, adminAddress, allowlist.AdminRole)
			},
		},
		{
			name: "should verify new address has no role",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				testContractAddr, _ := deployRewardManagerTest(t, backend, admin)
				allowlisttest.VerifyRole(t, rewardManager, testContractAddr, allowlist.NoRole)
			},
		},
		{
			name: "should not allow non-enabled to set reward address",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				testContractAddr, testContract := deployRewardManagerTest(t, backend, admin)

				allowlisttest.VerifyRole(t, rewardManager, testContractAddr, allowlist.NoRole)

				_, err := testContract.SetRewardAddress(admin, testContractAddr)
				require.ErrorContains(t, err, "execution reverted")
			},
		},
		{
			name: "should allow admin to enable contract",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				testContractAddr, _ := deployRewardManagerTest(t, backend, admin)

				allowlisttest.VerifyRole(t, rewardManager, testContractAddr, allowlist.NoRole)
				allowlisttest.SetAsEnabled(t, backend, rewardManager, admin, testContractAddr)
				allowlisttest.VerifyRole(t, rewardManager, testContractAddr, allowlist.EnabledRole)
			},
		},
		{
			name: "should allow enabled contract to set reward address",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				testContractAddr, testContract := deployRewardManagerTest(t, backend, admin)

				allowlisttest.SetAsEnabled(t, backend, rewardManager, admin, testContractAddr)

				tx, err := testContract.SetRewardAddress(admin, testContractAddr)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				currentAddr, err := testContract.CurrentRewardAddress(nil)
				require.NoError(t, err)
				require.Equal(t, testContractAddr, currentAddr)
			},
		},
		{
			name: "should return false for areFeeRecipientsAllowed by default",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				testContractAddr, testContract := deployRewardManagerTest(t, backend, admin)
				allowlisttest.SetAsEnabled(t, backend, rewardManager, admin, testContractAddr)

				isAllowed, err := testContract.AreFeeRecipientsAllowed(nil)
				require.NoError(t, err)
				require.False(t, isAllowed)
			},
		},
		{
			name: "should allow enabled contract to allow fee recipients",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				testContractAddr, testContract := deployRewardManagerTest(t, backend, admin)
				allowlisttest.SetAsEnabled(t, backend, rewardManager, admin, testContractAddr)

				tx, err := testContract.AllowFeeRecipients(admin)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				isAllowed, err := testContract.AreFeeRecipientsAllowed(nil)
				require.NoError(t, err)
				require.True(t, isAllowed)
			},
		},
		{
			name: "should allow enabled contract to disable rewards",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				testContractAddr, testContract := deployRewardManagerTest(t, backend, admin)
				allowlisttest.SetAsEnabled(t, backend, rewardManager, admin, testContractAddr)

				tx, err := testContract.SetRewardAddress(admin, testContractAddr)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				currentAddr, err := testContract.CurrentRewardAddress(nil)
				require.NoError(t, err)
				require.Equal(t, testContractAddr, currentAddr)

				tx, err = testContract.DisableRewards(admin)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				currentAddr, err = testContract.CurrentRewardAddress(nil)
				require.NoError(t, err)
				require.Equal(t, constants.BlackholeAddr, currentAddr)
			},
		},
		{
			name: "should return blackhole as default reward address",
			test: func(t *testing.T, _ *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				currentAddr, err := rewardManager.CurrentRewardAddress(nil)
				require.NoError(t, err)
				require.Equal(t, constants.BlackholeAddr, currentAddr)
			},
		},
		{
			name: "fees should go to blackhole by default",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				ctx := context.Background()
				client := backend.Client()

				initialBlackholeBalance, err := client.BalanceAt(ctx, constants.BlackholeAddr, nil)
				require.NoError(t, err)

				tx := testutils.SendSimpleTx(t, backend, adminKey)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				newBlackholeBalance, err := client.BalanceAt(ctx, constants.BlackholeAddr, nil)
				require.NoError(t, err)

				require.Greater(t, newBlackholeBalance.Cmp(initialBlackholeBalance), 0,
					"blackhole balance should have increased from fees")
			}},
		{
			name: "fees should go to configured reward address",
			test: func(t *testing.T, backend *sim.Backend, rewardManager *rewardmanagerbindings.IRewardManager) {
				ctx := context.Background()
				client := backend.Client()

				rewardRecipientAddr, _ := deployRewardManagerTest(t, backend, admin)

				initialRecipientBalance, err := client.BalanceAt(ctx, rewardRecipientAddr, nil)
				require.NoError(t, err)

				allowlisttest.SetAsEnabled(t, backend, rewardManager, admin, rewardRecipientAddr)

				tx, err := rewardManager.SetRewardAddress(admin, rewardRecipientAddr)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				currentAddr, err := rewardManager.CurrentRewardAddress(nil)
				require.NoError(t, err)
				require.Equal(t, rewardRecipientAddr, currentAddr)

				// Send a transaction to generate fees
				// The fees from THIS transaction should go to the reward address
				tx = testutils.SendSimpleTx(t, backend, adminKey)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				newRecipientBalance, err := client.BalanceAt(ctx, rewardRecipientAddr, nil)
				require.NoError(t, err)

				require.Greater(t, newRecipientBalance.Cmp(initialRecipientBalance), 0,
					"reward recipient balance should have increased from fees")
			}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			backend := newBackendWithRewardManager(t)
			defer backend.Close()

			rewardManager, err := rewardmanagerbindings.NewIRewardManager(rewardmanager.ContractAddress, backend.Client())
			require.NoError(t, err)

			tc.test(t, backend, rewardManager)
		})
	}
}

func TestIRewardManager_Events(t *testing.T) {
	chainID := params.TestChainConfig.ChainID
	admin := testutils.NewAuth(t, adminKey, chainID)
	testKey, _ := crypto.GenerateKey()
	testAddress := crypto.PubkeyToAddress(testKey.PublicKey)

	t.Run("should emit RewardAddressChanged event", func(t *testing.T) {
		backend := newBackendWithRewardManager(t)
		defer backend.Close()

		rewardManager, err := rewardmanagerbindings.NewIRewardManager(rewardmanager.ContractAddress, backend.Client())
		require.NoError(t, err)

		tx, err := rewardManager.SetRewardAddress(admin, testAddress)
		require.NoError(t, err)
		testutils.WaitReceiptSuccessful(t, backend, tx)

		iter, err := rewardManager.FilterRewardAddressChanged(nil, nil, nil, nil)
		require.NoError(t, err)
		defer iter.Close()

		require.True(t, iter.Next(), "expected to find RewardAddressChanged event")
		event := iter.Event
		require.Equal(t, adminAddress, event.Sender, "sender mismatch")
		require.Equal(t, constants.BlackholeAddr, event.OldRewardAddress, "old reward address mismatch")
		require.Equal(t, testAddress, event.NewRewardAddress, "new reward address mismatch")

		require.False(t, iter.Next(), "expected no more events")
		require.NoError(t, iter.Error())
	})

	t.Run("should emit FeeRecipientsAllowed event", func(t *testing.T) {
		backend := newBackendWithRewardManager(t)
		defer backend.Close()

		rewardManager, err := rewardmanagerbindings.NewIRewardManager(rewardmanager.ContractAddress, backend.Client())
		require.NoError(t, err)

		tx, err := rewardManager.AllowFeeRecipients(admin)
		require.NoError(t, err)
		testutils.WaitReceiptSuccessful(t, backend, tx)

		iter, err := rewardManager.FilterFeeRecipientsAllowed(nil, nil)
		require.NoError(t, err)
		defer iter.Close()

		require.True(t, iter.Next(), "expected to find FeeRecipientsAllowed event")
		event := iter.Event
		require.Equal(t, adminAddress, event.Sender, "sender mismatch")

		require.False(t, iter.Next(), "expected no more events")
		require.NoError(t, iter.Error())
	})

	t.Run("should emit RewardsDisabled event", func(t *testing.T) {
		backend := newBackendWithRewardManager(t)
		defer backend.Close()

		rewardManager, err := rewardmanagerbindings.NewIRewardManager(rewardmanager.ContractAddress, backend.Client())
		require.NoError(t, err)

		tx, err := rewardManager.DisableRewards(admin)
		require.NoError(t, err)
		testutils.WaitReceiptSuccessful(t, backend, tx)

		iter, err := rewardManager.FilterRewardsDisabled(nil, nil)
		require.NoError(t, err)
		defer iter.Close()

		require.True(t, iter.Next(), "expected to find RewardsDisabled event")
		event := iter.Event
		require.Equal(t, adminAddress, event.Sender, "sender mismatch")

		require.False(t, iter.Next(), "expected no more events")
		require.NoError(t, iter.Error())
	})
}

// TODO(jonathanoppenheimer): uncomment this once RunAllowListEventTests() is merged into main
// func TestIAllowList_Events(t *testing.T) {
// 	admin := testutils.NewAuth(t, adminKey, params.TestChainConfig.ChainID)
// 	allowlisttest.RunAllowListEventTests(t, newBackendWithRewardManager, rewardmanager.ContractAddress, admin, adminAddress)
// }
