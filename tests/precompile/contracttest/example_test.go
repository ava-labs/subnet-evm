// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package contracttest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/contracts/bindings"
)

// TestExampleDeployment demonstrates the test infrastructure
// This is a simple example showing how to deploy and interact with contracts
func TestExampleDeployment(t *testing.T) {

	// Create test backend with funded accounts
	backend := NewTestBackend(t)
	defer backend.Close()

	// Deploy ExampleDeployerList contract using generated binding
	addr, tx, contract, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client())
	require.NoError(err, "bindings.DeployExampleDeployerList(...)")
	require.NotEqual(addr.Hex(), "0x0000000000000000000000000000000000000000", "contract address should not be zero")

	// Wait for deployment
	receipt := WaitForReceipt(t, backend, tx)
	require.Equalf(t, types.ReceiptStatusSuccessful, receipt.Status, "%T.Status", receipt)

	require.NotNil(contract, "contract should not be nil")

	// Example: Check owner (should be admin account)
	owner, err := contract.Owner(nil)
	require.NoError(err)
	require.Equal(backend.Admin.Address, owner, "owner should be admin account")
}

// TestAllowListRoles demonstrates role management
func TestAllowListRoles(t *testing.T) {
	// Create test backend
	backend := NewTestBackend(t)
	defer backend.Close()

	// Set admin role on the deployer allowlist
	SetupAllowListRole(t, backend, ContractDeployerAllowListAddress, backend.Admin.Address, RoleAdmin, backend.Admin)

	// Verify the role was set
	RequireRole(t, backend, ContractDeployerAllowListAddress,
		backend.Admin.Address, RoleAdmin)

	// Verify other address has no role
	RequireRole(t, backend, ContractDeployerAllowListAddress,
		backend.Unprivileged.Address, RoleNone)
}
