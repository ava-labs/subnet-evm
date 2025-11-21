// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativemintertest

import (
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/libevm/crypto"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/plugin/evm/customtypes"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/contracts/testutils"
	"github.com/ava-labs/subnet-evm/utils"

	sim "github.com/ava-labs/subnet-evm/ethclient/simulated"
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

func newBackendWithNativeMinter(t *testing.T) *sim.Backend {
	t.Helper()
	chainCfg := params.Copy(params.TestChainConfig)
	// Match the simulated backend chain ID used for signing (1337).
	chainCfg.ChainID = big.NewInt(1337)
	// Enable ContractNativeMinter at genesis with admin set to adminAddress.
	params.GetExtra(&chainCfg).GenesisPrecompiles = extras.Precompiles{
		nativeminter.ConfigKey: nativeminter.NewConfig(utils.NewUint64(0), []common.Address{adminAddress}, nil, nil, nil),
	}
	return sim.NewBackend(
		types.GenesisAlloc{
			adminAddress:        {Balance: big.NewInt(1000000000000000000)},
			unprivilegedAddress: {Balance: big.NewInt(1000000000000000000)},
		},
		sim.WithChainConfig(&chainCfg),
	)
}

// Helper functions to reduce test boilerplate

func deployERC20NativeMinterTest(t *testing.T, b *sim.Backend, auth *bind.TransactOpts, initSupply *big.Int) (common.Address, *ERC20NativeMinterTest) {
	t.Helper()
	addr, tx, contract, err := DeployERC20NativeMinterTest(auth, b.Client(), nativeminter.ContractAddress, initSupply)
	require.NoError(t, err)
	testutils.WaitReceiptSuccessful(t, b, tx)
	return addr, contract
}

func deployMinter(t *testing.T, b *sim.Backend, auth *bind.TransactOpts, tokenAddress common.Address) (common.Address, *Minter) {
	t.Helper()
	addr, tx, contract, err := DeployMinter(auth, b.Client(), tokenAddress)
	require.NoError(t, err)
	testutils.WaitReceiptSuccessful(t, b, tx)
	return addr, contract
}

func verifyRole(t *testing.T, nativeMinter *INativeMinter, address common.Address, expectedRole allowlist.Role) {
	t.Helper()
	role, err := nativeMinter.ReadAllowList(nil, address)
	require.NoError(t, err)
	require.Equal(t, expectedRole.Big(), role)
}

func setAsEnabled(t *testing.T, b *sim.Backend, nativeMinter *INativeMinter, auth *bind.TransactOpts, address common.Address) {
	t.Helper()
	tx, err := nativeMinter.SetEnabled(auth, address)
	require.NoError(t, err)
	testutils.WaitReceiptSuccessful(t, b, tx)
}

func TestNativeMinter(t *testing.T) {
	chainID := big.NewInt(1337)
	admin := testutils.NewAuth(t, adminKey, chainID)

	initSupply := big.NewInt(1000)
	amount := big.NewInt(100)

	type testCase struct {
		name string
		test func(t *testing.T, backend *sim.Backend, nativeMinterIntf *INativeMinter)
	}

	testCases := []testCase{
		{
			name: "contract should not be able to mintdraw",
			test: func(t *testing.T, backend *sim.Backend, nativeMinter *INativeMinter) {
				tokenAddr, token := deployERC20NativeMinterTest(t, backend, admin, initSupply)

				verifyRole(t, nativeMinter, tokenAddr, allowlist.NoRole)

				// Try to mintdraw - should fail because token contract is not enabled
				_, err := token.Mintdraw(admin, amount)
				require.ErrorContains(t, err, vm.ErrExecutionReverted.Error())
			},
		},
		{
			name: "should be added to minter list",
			test: func(t *testing.T, backend *sim.Backend, nativeMinter *INativeMinter) {
				tokenAddr, _ := deployERC20NativeMinterTest(t, backend, admin, initSupply)

				verifyRole(t, nativeMinter, tokenAddr, allowlist.NoRole)

				setAsEnabled(t, backend, nativeMinter, admin, tokenAddr)

				verifyRole(t, nativeMinter, tokenAddr, allowlist.EnabledRole)
			},
		},
		{
			name: "admin should mintdraw",
			test: func(t *testing.T, backend *sim.Backend, nativeMinter *INativeMinter) {
				tokenAddr, token := deployERC20NativeMinterTest(t, backend, admin, initSupply)

				setAsEnabled(t, backend, nativeMinter, admin, tokenAddr)

				// Get initial token balance (admin has the tokens)
				initialTokenBalance, err := token.BalanceOf(nil, adminAddress)
				require.NoError(t, err)

				// Perform mintdraw - burns admin's ERC20 tokens and mints native coins to admin
				tx, err := token.Mintdraw(admin, amount)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				// Verify ERC20 token balance decreased by exactly the amount
				finalTokenBalance, err := token.BalanceOf(nil, adminAddress)
				require.NoError(t, err)
				expectedTokenBalance := new(big.Int).Sub(initialTokenBalance, amount)
				require.Zero(t, expectedTokenBalance.Cmp(finalTokenBalance), "ERC20 balance should have decreased by amount")
			},
		},
		{
			name: "minter should not mintdraw without tokens",
			test: func(t *testing.T, backend *sim.Backend, nativeMinter *INativeMinter) {
				tokenAddr, token := deployERC20NativeMinterTest(t, backend, admin, initSupply)
				minterAddr, minter := deployMinter(t, backend, admin, tokenAddr)

				setAsEnabled(t, backend, nativeMinter, admin, tokenAddr)

				// Verify minter has no ERC20 tokens
				initialTokenBalance, err := token.BalanceOf(nil, minterAddr)
				require.NoError(t, err)
				require.Zero(t, initialTokenBalance.Cmp(big.NewInt(0)), "minter should have no ERC20 tokens")

				// Try to mintdraw - should fail because minter has no ERC20 tokens to burn
				_, err = minter.Mintdraw(admin, amount)
				require.Error(t, err)
			},
		},
		{
			name: "should deposit for minter",
			test: func(t *testing.T, backend *sim.Backend, nativeMinter *INativeMinter) {
				tokenAddr, token := deployERC20NativeMinterTest(t, backend, admin, initSupply)
				minterAddr, minter := deployMinter(t, backend, admin, tokenAddr)

				setAsEnabled(t, backend, nativeMinter, admin, tokenAddr)

				// Mint native coins to minter address
				tx, err := nativeMinter.MintNativeCoin(admin, minterAddr, amount)
				require.NoError(t, err)
				testutils.WaitReceipt(t, backend, tx)

				// Get initial balances
				initialTokenBalance, err := token.BalanceOf(nil, minterAddr)
				require.NoError(t, err)
				initialNativeBalance, err := backend.Client().BalanceAt(t.Context(), minterAddr, nil)
				require.NoError(t, err)

				// Deposit (convert native coin to ERC20 token)
				tx, err = minter.Deposit(admin, amount)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				// Verify balances changed correctly
				finalTokenBalance, err := token.BalanceOf(nil, minterAddr)
				require.NoError(t, err)
				expectedTokenBalance := new(big.Int).Add(initialTokenBalance, amount)
				require.Equal(t, expectedTokenBalance, finalTokenBalance)

				finalNativeBalance, err := backend.Client().BalanceAt(t.Context(), minterAddr, nil)
				require.NoError(t, err)
				expectedNativeBalance := new(big.Int).Sub(initialNativeBalance, amount)
				require.Equal(t, expectedNativeBalance, finalNativeBalance)
			},
		},
		{
			name: "minter should mintdraw",
			test: func(t *testing.T, backend *sim.Backend, nativeMinter *INativeMinter) {
				tokenAddr, token := deployERC20NativeMinterTest(t, backend, admin, initSupply)
				minterAddr, minter := deployMinter(t, backend, admin, tokenAddr)

				setAsEnabled(t, backend, nativeMinter, admin, tokenAddr)

				// Verify minter starts with no native balance
				initialNativeBalance, err := backend.Client().BalanceAt(t.Context(), minterAddr, nil)
				require.NoError(t, err)
				require.Zero(t, initialNativeBalance.Cmp(big.NewInt(0)), "minter should start with no native balance")

				// Mint ERC20 tokens to minter
				tx, err := token.Mint(admin, minterAddr, amount)
				require.NoError(t, err)
				testutils.WaitReceipt(t, backend, tx)

				// Verify minter received tokens
				initialTokenBalance, err := token.BalanceOf(nil, minterAddr)
				require.NoError(t, err)
				require.Zero(t, amount.Cmp(initialTokenBalance), "minter should have received ERC20 tokens")

				// Mintdraw (convert ERC20 token to native coin)
				tx, err = minter.Mintdraw(admin, amount)
				require.NoError(t, err)
				testutils.WaitReceiptSuccessful(t, backend, tx)

				// Verify final balances
				finalTokenBalance, err := token.BalanceOf(nil, minterAddr)
				require.NoError(t, err)
				require.Zero(t, finalTokenBalance.Cmp(big.NewInt(0)), "minter should have no ERC20 tokens left")

				finalNativeBalance, err := backend.Client().BalanceAt(t.Context(), minterAddr, nil)
				require.NoError(t, err)
				require.Zero(t, amount.Cmp(finalNativeBalance), "minter should have received native coins")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			backend := newBackendWithNativeMinter(t)
			defer backend.Close()

			nativeMinter, err := NewINativeMinter(nativeminter.ContractAddress, backend.Client())
			require.NoError(t, err)

			tc.test(t, backend, nativeMinter)
		})
	}
}

func TestINativeMinter_Events(t *testing.T) {
	chainID := big.NewInt(1337)
	admin := testutils.NewAuth(t, adminKey, chainID)
	testKey, _ := crypto.GenerateKey()
	testAddress := crypto.PubkeyToAddress(testKey.PublicKey)

	backend := newBackendWithNativeMinter(t)
	defer backend.Close()

	nativeMinter, err := NewINativeMinter(nativeminter.ContractAddress, backend.Client())
	require.NoError(t, err)

	t.Run("should emit NativeCoinMinted event", func(t *testing.T) {
		require := require.New(t)

		amount := big.NewInt(1000)

		tx, err := nativeMinter.MintNativeCoin(admin, testAddress, amount)
		require.NoError(err)
		testutils.WaitReceiptSuccessful(t, backend, tx)

		// Filter for NativeCoinMinted events
		iter, err := nativeMinter.FilterNativeCoinMinted(
			nil,
			[]common.Address{adminAddress},
			[]common.Address{testAddress},
		)
		require.NoError(err)
		defer iter.Close()

		// Verify event fields match expected values
		require.True(iter.Next(), "expected to find NativeCoinMinted event")
		event := iter.Event
		require.Equal(adminAddress, event.Sender)
		require.Equal(testAddress, event.Recipient)
		require.Zero(amount.Cmp(event.Amount), "amount mismatch")

		// Verify there are no more events
		require.False(iter.Next(), "expected no more NativeCoinMinted events")
		require.NoError(iter.Error())
	})
}
