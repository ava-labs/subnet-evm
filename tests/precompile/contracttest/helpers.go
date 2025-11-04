// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package contracttest provides utilities for testing contracts in Go.
package contracttest

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/accounts/abi/bind"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/ava-labs/libevm/ethclient/simulated"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

const simulatedBackendChainID = 1337

var (
	// AdminAddress is the primary admin account
	AdminAddress = common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	// UnprivilegedAddress is an unprivileged account for testing (has no special roles)
	UnprivilegedAddress = common.HexToAddress("0x0Fa8EA536Be85F32724D57A37758761B86416123")
)

// Keys for test accounts
var testPrivateKeys = []string{
	"56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027",
	"7b4198529994b0dc604278c99d153cfd069d594753d471171a1d102a10438e07",
	"15614556be13730e9e8d6eacc1603143e7b96987429df8726384c2ec4502ef6e",
	"31b571bf6894a248831ff937bb49f7754509fe93bbd2517c9c73c4144c0e97dc",
	"6934bef917e01692b789da754a0eae31a8536eb465e7bff752ea291dad88c675",
	"e700bdbdbc279b808b1ec45f8c2370e4616d3a02c336e68d85d4668e08f53cff",
	"bbc2865b76ba28016bc2255c7504d000e046ae01934b04c694592a6276988630",
	"cdbfd34f687ced8c6968854f8a99ae47712c4f4183b78dcc4a903d1bfe8cbf60",
	"86f78c5416151fe3546dece84fda4b4b1e36089f2dbc48496faf3a950f16157c",
	"750839e9dbbd2a0910efe40f50b2f3b2f2f59f5580bb4b83bd8c1201cf9a010a",
}

// Account represents a funded test account
type Account struct {
	Key     *ecdsa.PrivateKey
	Address common.Address
	Auth    *bind.TransactOpts
}

// Accounts provides convenient access to test accounts.
// Admin and Unprivileged are convenience pointers to specific accounts in AllAccounts.
type Accounts struct {
	Admin        *Account   // Points to admin account in AllAccounts
	Unprivileged *Account   // Points to unprivileged test account in AllAccounts (has no special roles)
	AllAccounts  []*Account // All funded test accounts (includes Admin and Unprivileged)
}

type Backend struct {
	*simulated.Backend
	Accounts
}

// NewTestBackend creates a simulated backend with funded test accounts
func NewTestBackend(t testing.TB) *Backend {
	// Create test accounts
	var fundedAccounts Accounts
	fundedAccounts.AllAccounts = make([]*Account, len(testPrivateKeys))

	// Set up genesis allocation with funded accounts
	alloc := types.GenesisAlloc{}
	balance := new(uint256.Int).SetAllOne().ToBig()

	for i, keyHex := range testPrivateKeys {
		key, err := crypto.HexToECDSA(keyHex)
		require.NoError(t, err)

		addr := crypto.PubkeyToAddress(key.PublicKey)
		alloc[addr] = types.Account{Balance: balance}

		auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(simulatedBackendChainID))
		require.NoError(t, err)

		account := &Account{
			Key:     key,
			Address: addr,
			Auth:    auth,
		}
		fundedAccounts.AllAccounts[i] = account

		switch addr {
		case AdminAddress:
			fundedAccounts.Admin = account
		case UnprivilegedAddress:
			fundedAccounts.Unprivileged = account
		}
	}

	backend := &Backend{
		Backend:  simulated.NewBackend(alloc),
		Accounts: fundedAccounts,
	}

	t.Cleanup(func() { backend.Close() })

	return backend
}

// WaitForReceipt waits for a transaction receipt and commits a block if needed
func WaitForReceipt(t testing.TB, tb *Backend, tx *types.Transaction) *types.Receipt {
	t.Helper()
	tb.Commit()

	client := tb.Client()
	receipt, err := client.TransactionReceipt(t.Context(), tx.Hash())
	require.NoError(t, err, "failed to get transaction receipt")

	return receipt
}
