// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// This test suite migrates the Hardhat tests for the ContractDeployerAllowList
// to Go using the simulated backend and the generated bindings in this package.
package deployerallowlisttest

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/plugin/evm/customtypes"

	sim "github.com/ava-labs/subnet-evm/ethclient/simulated"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/stretchr/testify/require"
)

// Test keys matching the Hardhat suite identities.
const ()

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

func newAuth(t *testing.T, key *ecdsa.PrivateKey, chainID *big.Int) *bind.TransactOpts {
	t.Helper()
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	require.NoError(t, err)
	return auth
}

func newBackendWithDeployerAllowList(t *testing.T) *sim.Backend {
	t.Helper()
	chainCfg := params.Copy(params.TestChainConfig)
	// Match the simulated backend chain ID used for signing (1337).
	chainCfg.ChainID = big.NewInt(1337)
	// Enable ContractDeployerAllowList at genesis with admin set to adminAddress.
	params.GetExtra(&chainCfg).GenesisPrecompiles = extras.Precompiles{
		deployerallowlist.ConfigKey: deployerallowlist.NewConfig(utils.NewUint64(0), []common.Address{adminAddress}, nil, nil),
	}
	return sim.NewBackend(
		types.GenesisAlloc{
			adminAddress:        {Balance: big.NewInt(1000000000000000000)},
			unprivilegedAddress: {Balance: big.NewInt(1000000000000000000)},
		},
		sim.WithChainConfig(&chainCfg),
	)
}

func waitReceipt(t *testing.T, b *sim.Backend, tx *types.Transaction) *types.Receipt {
	t.Helper()
	b.Commit(true)
	receipt, err := b.Client().TransactionReceipt(t.Context(), tx.Hash())
	require.NoError(t, err)
	return receipt
}

func TestDeployerAllowList_Steps(t *testing.T) {
	chainID := big.NewInt(1337)
	admin := newAuth(t, adminKey, chainID)

	backend := newBackendWithDeployerAllowList(t)
	defer backend.Close()

	type testCase struct {
		name           string
		runStep        func(*DeployerListTest, *bind.TransactOpts) (*types.Transaction, error)
		expectedStatus uint64
	}

	testCases := []testCase{
		{
			name:           "step_verifySenderIsAdmin",
			runStep:        (*DeployerListTest).StepVerifySenderIsAdmin,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
		{
			name:           "step_newAddressHasNoRole",
			runStep:        (*DeployerListTest).StepNewAddressHasNoRole,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
		{
			name:           "step_noRoleIsNotAdmin",
			runStep:        (*DeployerListTest).StepNoRoleIsNotAdmin,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
		{
			name:           "step_noRoleCannotDeploy",
			runStep:        (*DeployerListTest).StepNoRoleCannotDeploy,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
		{
			name:           "step_adminAddContractAsAdmin",
			runStep:        (*DeployerListTest).StepAdminAddContractAsAdmin,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
		{
			name:           "step_addDeployerThroughContract",
			runStep:        (*DeployerListTest).StepAddDeployerThroughContract,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
		{
			name:           "step_deployerCanDeploy",
			runStep:        (*DeployerListTest).StepDeployerCanDeploy,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
		{
			name:           "step_adminCanRevokeDeployer",
			runStep:        (*DeployerListTest).StepAdminCanRevokeDeployer,
			expectedStatus: types.ReceiptStatusSuccessful,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			testContractAddr, tx, testContract, err := DeployDeployerListTest(admin, backend.Client())
			require.NoError(err)
			require.Equal(types.ReceiptStatusSuccessful, waitReceipt(t, backend, tx).Status)

			allowList, err := allowlist.NewIAllowList(deployerallowlist.ContractAddress, backend.Client())
			require.NoError(err)
			// Set the contract address as admin in the deployer allow list precompile to enable the contract to deploy and
			// modify the allow list.
			tx, err = allowList.SetAdmin(admin, testContractAddr)
			require.NoError(err)
			require.Equal(types.ReceiptStatusSuccessful, waitReceipt(t, backend, tx).Status)

			// Run the setup method to initialize the test contract.
			tx, err = testContract.SetUp(admin)
			require.NoError(err)
			require.Equal(types.ReceiptStatusSuccessful, waitReceipt(t, backend, tx).Status)

			auth := newAuth(t, adminKey, chainID)
			tx, err = tc.runStep(testContract, auth)
			require.NoError(err)
			require.Equal(tc.expectedStatus, waitReceipt(t, backend, tx).Status)
		})
	}
}

func TestIAllowList_Events(t *testing.T) {
	chainID := big.NewInt(1337)
	admin := newAuth(t, adminKey, chainID)
	testKey, _ := crypto.GenerateKey()
	testAddress := crypto.PubkeyToAddress(testKey.PublicKey)

	type testCase struct {
		name           string
		setup          func(*allowlist.IAllowList, *bind.TransactOpts, *sim.Backend, *testing.T, common.Address) error
		runMethod      func(*allowlist.IAllowList, *bind.TransactOpts, common.Address) (*types.Transaction, error)
		expectedEvents []allowlist.IAllowListRoleSet
	}

	testCases := []testCase{
		{
			name: "should emit event after set admin",
			runMethod: func(allowList *allowlist.IAllowList, auth *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
				return allowList.SetAdmin(auth, addr)
			},
			expectedEvents: []allowlist.IAllowListRoleSet{
				{
					Role:    allowlist.AdminRole.Big(),
					Account: testAddress,
					Sender:  adminAddress,
					OldRole: allowlist.NoRole.Big(),
				},
			},
		},
		{
			name: "should emit event after set manager",
			runMethod: func(allowList *allowlist.IAllowList, auth *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
				return allowList.SetManager(auth, addr)
			},
			expectedEvents: []allowlist.IAllowListRoleSet{
				{
					Role:    allowlist.ManagerRole.Big(),
					Account: testAddress,
					Sender:  adminAddress,
					OldRole: allowlist.NoRole.Big(),
				},
			},
		},
		{
			name: "should emit event after set enabled",
			runMethod: func(allowList *allowlist.IAllowList, auth *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
				return allowList.SetEnabled(auth, addr)
			},
			expectedEvents: []allowlist.IAllowListRoleSet{
				{
					Role:    allowlist.EnabledRole.Big(),
					Account: testAddress,
					Sender:  adminAddress,
					OldRole: allowlist.NoRole.Big(),
				},
			},
		},
		{
			name: "should emit event after set none",
			setup: func(allowList *allowlist.IAllowList, auth *bind.TransactOpts, backend *sim.Backend, t *testing.T, addr common.Address) error {
				// First set the address to Enabled so we can test setting it to None
				tx, err := allowList.SetEnabled(auth, addr)
				if err != nil {
					return err
				}
				waitReceipt(t, backend, tx)
				return nil
			},
			runMethod: func(allowList *allowlist.IAllowList, auth *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
				return allowList.SetNone(auth, addr)
			},
			expectedEvents: []allowlist.IAllowListRoleSet{
				{
					Role:    allowlist.EnabledRole.Big(),
					Account: testAddress,
					Sender:  adminAddress,
					OldRole: allowlist.NoRole.Big(),
				},
				{
					Role:    allowlist.NoRole.Big(),
					Account: testAddress,
					Sender:  adminAddress,
					OldRole: allowlist.EnabledRole.Big(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			backend := newBackendWithDeployerAllowList(t)
			defer backend.Close()

			allowList, err := allowlist.NewIAllowList(deployerallowlist.ContractAddress, backend.Client())
			require.NoError(err)

			if tc.setup != nil {
				err := tc.setup(allowList, admin, backend, t, testAddress)
				require.NoError(err)
			}

			tx, err := tc.runMethod(allowList, admin, testAddress)
			require.NoError(err)
			receipt := waitReceipt(t, backend, tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			// Filter for RoleSet events using FilterRoleSet
			// This will filter for all RoleSet events.
			iter, err := allowList.FilterRoleSet(
				nil,
				nil,
				nil,
				nil,
			)
			require.NoError(err)
			defer iter.Close()

			// Verify event fields match expected values
			for _, expectedEvent := range tc.expectedEvents {
				require.True(iter.Next(), "expected to find RoleSet event")
				event := iter.Event
				require.Equal(0, expectedEvent.Role.Cmp(event.Role), "role mismatch")
				require.Equal(expectedEvent.Account, event.Account, "account mismatch")
				require.Equal(expectedEvent.Sender, event.Sender, "sender mismatch")
				require.Equal(0, expectedEvent.OldRole.Cmp(event.OldRole), "oldRole mismatch")
			}

			// Verify there are no more events
			require.False(iter.Next(), "expected no more RoleSet events")
			require.NoError(iter.Error())
		})
	}
}
