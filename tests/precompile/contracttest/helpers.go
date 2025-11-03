// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package contracttest provides utilities for testing contracts in Go.
package contracttest

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/accounts/abi/bind"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/ava-labs/libevm/ethclient/simulated"
	"github.com/stretchr/testify/require"
)

var _ bind.ContractBackend = (simulated.Client)(nil)

// TestBackend wraps a simulated backend with common test utilities
type TestBackend struct {
	*simulated.Backend
	Accounts
}

// Common test accounts matching the TypeScript test setup
type Accounts struct {
	Admin      *TestAccount
	OtherAddr  *TestAccount
	FundedKeys []*TestAccount
}

// TestAccount represents a funded test account
type TestAccount struct {
	Key     *ecdsa.PrivateKey
	Address common.Address
	Auth    *bind.TransactOpts
}

// Common addresses from TypeScript tests
var (
	// AdminAddress is the primary admin account
	AdminAddress = common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	// OtherAddress is a secondary account for testing
	OtherAddress = common.HexToAddress("0x0Fa8EA536Be85F32724D57A37758761B86416123")
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

// NewTestBackend creates a simulated backend with funded test accounts
func NewTestBackend(t testing.TB) *TestBackend {
	require := require.New(t)

	// Create test accounts
	var fundedAccounts Accounts
	fundedAccounts.FundedKeys = make([]*TestAccount, len(testPrivateKeys))

	// Set up genesis allocation with funded accounts
	alloc := types.GenesisAlloc{}
	balance := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18)) // 1000 ETH each

	for i, keyHex := range testPrivateKeys {
		key, err := crypto.HexToECDSA(keyHex)
		require.NoError(err)

		addr := crypto.PubkeyToAddress(key.PublicKey)
		alloc[addr] = types.Account{Balance: balance}

		chainID := big.NewInt(1337) // simulated backend uses chainID 1337
		auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
		require.NoError(err)

		account := &TestAccount{
			Key:     key,
			Address: addr,
			Auth:    auth,
		}
		fundedAccounts.FundedKeys[i] = account

		// Set up special accounts
		switch addr {
		case AdminAddress:
			fundedAccounts.Admin = account
		case OtherAddress:
			fundedAccounts.OtherAddr = account
		}
	}

	// Create simulated backend from libevm (uses AllDevChainProtocolChanges by default)
	backend := simulated.NewBackend(alloc)

	return &TestBackend{
		Backend:  backend,
		Accounts: fundedAccounts,
	}
}

// GetAccount returns a TestAccount by address, or nil if not found
func (tb *TestBackend) GetAccount(addr common.Address) *TestAccount {
	for _, account := range tb.FundedKeys {
		if account.Address == addr {
			return account
		}
	}
	return nil
}

// WaitForReceipt waits for a transaction receipt and commits a block if needed
func WaitForReceipt(t testing.TB, tb *TestBackend, tx *types.Transaction) *types.Receipt {
	require := require.New(t)

	// Commit the transaction (mines a block)
	tb.Commit()

	// Get the receipt
	client := tb.Client()
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(err, "failed to get transaction receipt")
	require.NotNil(receipt, "receipt is nil")

	return receipt
}

// RequireSuccessReceipt asserts that a transaction was successful
func RequireSuccessReceipt(t testing.TB, receipt *types.Receipt) {
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "transaction failed")
}
