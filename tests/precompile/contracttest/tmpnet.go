// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package contracttest

import (
	"context"
	"testing"

	"github.com/ava-labs/libevm/accounts/abi/bind"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/ava-labs/libevm/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/contracts/bindings"
)

// TmpnetBackend wraps an ethclient connection to a tmpnet subnet
// This is used for E2E testing against real subnets with precompiles enabled
type TmpnetBackend struct {
	Client *ethclient.Client
	Accounts
}

// NewTmpnetBackendSimple creates a test backend connected to a tmpnet subnet via RPC
// This version doesn't require a testing.TB interface, suitable for use in Ginkgo tests
func NewTmpnetBackendSimple(rpcURL string) *TmpnetBackend {
	// Connect to the subnet
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		panic("failed to connect to tmpnet subnet: " + err.Error())
	}

	// Create test accounts (same keys as simulated backend)
	var fundedAccounts Accounts
	fundedAccounts.FundedKeys = make([]*TestAccount, len(testPrivateKeys))

	// Get chainID from the connected chain
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic("failed to get chain ID: " + err.Error())
	}

	for i, keyHex := range testPrivateKeys {
		key, err := crypto.HexToECDSA(keyHex)
		if err != nil {
			panic("failed to parse private key: " + err.Error())
		}

		addr := crypto.PubkeyToAddress(key.PublicKey)

		auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
		if err != nil {
			panic("failed to create transactor: " + err.Error())
		}

		// Leave Nonce as nil so the binding fetches it automatically for each transaction
		// This is essential for parallel tests to avoid nonce conflicts

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

	return &TmpnetBackend{
		Client:   client,
		Accounts: fundedAccounts,
	}
}

// NewTmpnetBackend creates a test backend connected to a tmpnet subnet via RPC
func NewTmpnetBackend(t testing.TB, rpcURL string) *TmpnetBackend {
	require := require.New(t)

	// Connect to the subnet
	client, err := ethclient.Dial(rpcURL)
	require.NoError(err, "failed to connect to tmpnet subnet")

	// Create test accounts (same keys as simulated backend)
	var fundedAccounts Accounts
	fundedAccounts.FundedKeys = make([]*TestAccount, len(testPrivateKeys))

	// Get chainID from the connected chain
	chainID, err := client.ChainID(context.Background())
	require.NoError(err, "failed to get chain ID")

	for i, keyHex := range testPrivateKeys {
		key, err := crypto.HexToECDSA(keyHex)
		require.NoError(err)

		addr := crypto.PubkeyToAddress(key.PublicKey)

		auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
		require.NoError(err)

		// Leave Nonce as nil so the binding fetches it automatically for each transaction
		// This is essential for parallel tests to avoid nonce conflicts

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

	return &TmpnetBackend{
		Client:   client,
		Accounts: fundedAccounts,
	}
}

// Close closes the client connection
func (tb *TmpnetBackend) Close() {
	if tb.Client != nil {
		tb.Client.Close()
	}
}

// GetAccount returns a TestAccount by address, or nil if not found
func (tb *TmpnetBackend) GetAccount(addr common.Address) *TestAccount {
	for _, account := range tb.FundedKeys {
		if account.Address == addr {
			return account
		}
	}
	return nil
}

// WaitForTxReceipt waits for a transaction receipt from the RPC endpoint
func WaitForTxReceipt(t testing.TB, client *ethclient.Client, tx *types.Transaction) *types.Receipt {
	require := require.New(t)

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	require.NoError(err, "failed to wait for transaction")
	require.NotNil(receipt, "receipt is nil")

	return receipt
}

// DeployContractToTmpnet deploys a contract to a tmpnet subnet
// Example usage:
//
//	addr, tx, contract, err := bindings.DeployExampleDeployerList(auth, client)
//	receipt := WaitForTxReceipt(t, client, tx)
func DeployContractToTmpnet(
	t testing.TB,
	client *ethclient.Client,
	tx *types.Transaction,
) *types.Receipt {
	return WaitForTxReceipt(t, client, tx)
}

// SetupAllowListRoleOnTmpnet configures an address with a specific role on a tmpnet subnet
func SetupAllowListRoleOnTmpnet(
	t testing.TB,
	client *ethclient.Client,
	allowListAddress common.Address,
	targetAddress common.Address,
	role uint8,
	fromAccount *TestAccount,
) {
	require := require.New(t)

	// Get the IAllowList interface at the precompile address
	allowList, err := bindings.NewIAllowList(allowListAddress, client)
	require.NoError(err, "failed to create allowlist interface")

	var tx *types.Transaction
	switch role {
	case RoleAdmin:
		tx, err = allowList.SetAdmin(fromAccount.Auth, targetAddress)
	case RoleManager:
		tx, err = allowList.SetManager(fromAccount.Auth, targetAddress)
	case RoleEnabled:
		tx, err = allowList.SetEnabled(fromAccount.Auth, targetAddress)
	case RoleNone:
		tx, err = allowList.SetNone(fromAccount.Auth, targetAddress)
	default:
		require.Fail("invalid role")
	}

	require.NoError(err, "failed to set role")

	// Wait for transaction and verify
	receipt := WaitForTxReceipt(t, client, tx)
	RequireSuccessReceipt(t, receipt)
}

// GetAllowListRoleOnTmpnet returns the role of an address on a tmpnet subnet
func GetAllowListRoleOnTmpnet(
	t testing.TB,
	client *ethclient.Client,
	allowListAddress common.Address,
	targetAddress common.Address,
) uint8 {
	require := require.New(t)

	// Get the IAllowList interface at the precompile address
	allowList, err := bindings.NewIAllowList(allowListAddress, client)
	require.NoError(err, "failed to create allowlist interface")

	role, err := allowList.ReadAllowList(&bind.CallOpts{}, targetAddress)
	require.NoError(err, "failed to read allowlist role")

	return uint8(role.Uint64())
}
