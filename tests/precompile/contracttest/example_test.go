// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package contracttest

import (
	"testing"

	"github.com/ava-labs/libevm/core/types"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/contracts/bindings"
)

// TestExampleDeployment demonstrates the test infrastructure
// This is a simple example showing how to deploy and interact with contracts
func TestExampleDeployment(t *testing.T) {
	// Create test backend with funded accounts
	backend := NewTestBackend(t)

	// Deploy ExampleDeployerList contract using generated binding
	addr, tx, contract, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client())
	require.NoError(t, err, "bindings.DeployExampleDeployerList(...)")
	require.NotZero(t, addr.Hex(), "contract address should not be zero")

	// Wait for deployment
	receipt := WaitForReceipt(t, backend, tx)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "%T.Status", receipt)

	// Example: Check owner (should be admin account)
	owner, err := contract.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, backend.Admin.Address, owner, "owner should be admin account")
}

// TestAllowListRoles demonstrates role management
func TestAllowListRoles(t *testing.T) {
	// Create test backend
	backend := NewTestBackend(t)
	defer backend.Close()

	// Set admin role on the deployer allowlist
	SetupAllowListRole(t, backend, ContractDeployerAllowListAddress, backend.Admin.Address, RoleAdmin, backend.Admin)

	// Verify the role was set
	RequireRole(t, backend, ContractDeployerAllowListAddress, backend.Admin.Address, RoleAdmin)

	// Verify other address has no role
	RequireRole(t, backend, ContractDeployerAllowListAddress, backend.Unprivileged.Address, RoleNone)
}
