// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package contracttest

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/accounts/abi/bind"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/ava-labs/libevm/ethclient"
	"github.com/stretchr/testify/require"
)

// TmpnetBackend wraps an ethclient connection to a tmpnet subnet
// This is used for E2E testing against real subnets with precompiles enabled
type TmpnetBackend struct {
	Client *ethclient.Client
	Accounts
}

// NewTmpnetBackend creates a test backend connected to a tmpnet subnet via RPC
func NewTmpnetBackend(t testing.TB, rpcURL string) *TmpnetBackend {
	// Connect to the subnet
	client, err := ethclient.Dial(rpcURL)
	require.NoError(t, err, "failed to connect to tmpnet subnet")

	// Create test accounts (same keys as simulated backend)
	var fundedAccounts Accounts
	fundedAccounts.AllAccounts = make([]*Account, len(testPrivateKeys))

	// Get chainID from the connected chain
	chainID, err := client.ChainID(context.Background())
	require.NoError(t, err, "failed to get chain ID")

	for i, keyHex := range testPrivateKeys {
		key, err := crypto.HexToECDSA(keyHex)
		require.NoError(t, err, "failed to parse private key at index %d", i)

		addr := crypto.PubkeyToAddress(key.PublicKey)

		auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
		require.NoError(t, err, "failed to create transactor for address %s", addr.Hex())

		account := &Account{
			Key:     key,
			Address: addr,
			Auth:    auth,
		}
		fundedAccounts.AllAccounts[i] = account

		// Set up special accounts as convenience pointers
		switch addr {
		case AdminAddress:
			fundedAccounts.Admin = account
		case UnprivilegedAddress:
			fundedAccounts.Unprivileged = account
		}
	}

	// Validate that required accounts were found
	require.NotNil(t, fundedAccounts.Admin, "Admin account not found - AdminAddress must match one of testPrivateKeys")
	require.NotNil(t, fundedAccounts.Unprivileged, "Unprivileged account not found - UnprivilegedAddress must match one of testPrivateKeys")

	// Verify that accounts are funded in genesis (at least Admin account needs balance for transactions)
	// This prevents silent failures later when trying to send transactions
	adminBalance, err := client.BalanceAt(t.Context(), fundedAccounts.Admin.Address, nil)
	require.NoError(t, err, "failed to check Admin account balance")
	require.NotNil(t, adminBalance, "Admin account balance is nil")
	require.Greater(t, adminBalance.Cmp(big.NewInt(0)), 0, "Admin account (%s) has zero balance - account must be funded in genesis file", fundedAccounts.Admin.Address.Hex())

	unprivilegedBalance, err := client.BalanceAt(t.Context(), fundedAccounts.Unprivileged.Address, nil)
	require.NoError(t, err, "failed to check Unprivileged account balance")
	require.NotNil(t, unprivilegedBalance, "Unprivileged account balance is nil")
	require.Greater(t, unprivilegedBalance.Cmp(big.NewInt(0)), 0, "Unprivileged account (%s) has zero balance - account must be funded in genesis file", fundedAccounts.Unprivileged.Address.Hex())

	backend := &TmpnetBackend{
		Client:   client,
		Accounts: fundedAccounts,
	}

	t.Cleanup(backend.Close)

	return backend
}

// Close closes the client connection
func (tb *TmpnetBackend) Close() {
	if tb.Client != nil {
		tb.Client.Close()
	}
}

// WaitForReceipt waits for a transaction receipt from the RPC endpoint
func (be *TmpnetBackend) WaitForReceipt(t testing.TB, tx *types.Transaction) *types.Receipt {
	t.Helper()

	receipt, err := bind.WaitMined(context.Background(), be.Client, tx)
	require.NoError(t, err, "failed to wait for transaction")
	require.NotNil(t, receipt, "receipt is nil")

	return receipt
}
