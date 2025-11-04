// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package contracttest

import (
	"testing"

	"github.com/ava-labs/libevm/accounts/abi/bind"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/contracts/bindings"
)

// AllowList roles (matching TypeScript Roles in utils.ts)
const (
	RoleNone    = uint8(0)
	RoleEnabled = uint8(1)
	RoleAdmin   = uint8(2)
	RoleManager = uint8(3)
)

// Precompile addresses
var (
	ContractDeployerAllowListAddress = common.HexToAddress("0x0200000000000000000000000000000000000000")
	TxAllowListAddress               = common.HexToAddress("0x0200000000000000000000000000000000000002")
	NativeMinterAddress              = common.HexToAddress("0x0200000000000000000000000000000000000001")
	FeeManagerAddress                = common.HexToAddress("0x0200000000000000000000000000000000000003")
	RewardManagerAddress             = common.HexToAddress("0x0200000000000000000000000000000000000004")
	WarpAddress                      = common.HexToAddress("0x0200000000000000000000000000000000000005")
)

// SetupAllowListRole configures an address with a specific role on an allowlist precompile
func SetupAllowListRole(
	t testing.TB,
	backend *Backend,
	allowListAddress common.Address,
	targetAddress common.Address,
	role uint8,
	fromAccount *Account,
) {
	require := require.New(t)

	// Get the IAllowList interface at the precompile address
	allowList, err := bindings.NewIAllowList(allowListAddress, backend.Client())
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

	// Commit and verify
	receipt := WaitForReceipt(t, backend, tx)
	RequireSuccessReceipt(t, receipt)
}

// GetAllowListRole returns the role of an address on an allowlist precompile
func GetAllowListRole(
	t testing.TB,
	backend *Backend,
	allowListAddress common.Address,
	targetAddress common.Address,
) uint8 {
	require := require.New(t)

	// Get the IAllowList interface at the precompile address
	allowList, err := bindings.NewIAllowList(allowListAddress, backend.Client())
	require.NoError(err, "failed to create allowlist interface")

	role, err := allowList.ReadAllowList(&bind.CallOpts{}, targetAddress)
	require.NoError(err, "failed to read allowlist role")

	return uint8(role.Uint64())
}

// RequireRole asserts that an address has the expected role
func RequireRole(
	t testing.TB,
	backend *Backend,
	allowListAddress common.Address,
	targetAddress common.Address,
	expectedRole uint8,
) {
	actualRole := GetAllowListRole(t, backend, allowListAddress, targetAddress)
	require.Equal(t, expectedRole, actualRole, "role mismatch for address %s", targetAddress.Hex())
}
