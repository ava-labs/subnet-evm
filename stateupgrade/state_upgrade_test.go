// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/common/math"
	"github.com/ava-labs/libevm/core/rawdb"
	"github.com/ava-labs/libevm/core/state"
	"github.com/ava-labs/libevm/core/types"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/core/extstate"
	"github.com/ava-labs/subnet-evm/params/extras"
)

func TestUpgradeAccount_BalanceChanges(t *testing.T) {
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	tests := []struct {
		name           string
		initialBalance *uint256.Int
		balanceChange  *math.HexOrDecimal256
		accountExists  bool
		wantBalance    *uint256.Int
		wantError      string
	}{
		{
			name:           "positive balance change on existing account",
			initialBalance: uint256.NewInt(100),
			balanceChange:  hexOrDecimal256FromInt64(50),
			accountExists:  true,
			wantBalance:    uint256.NewInt(150),
		},
		{
			name:           "negative balance change with sufficient funds",
			initialBalance: uint256.NewInt(100),
			balanceChange:  hexOrDecimal256FromInt64(-30),
			accountExists:  true,
			wantBalance:    uint256.NewInt(70),
		},
		{
			name:           "negative balance change with insufficient funds",
			initialBalance: uint256.NewInt(50),
			balanceChange:  hexOrDecimal256FromInt64(-100),
			accountExists:  true,
			wantBalance:    uint256.NewInt(50), // unchanged
			wantError:      "insufficient balance for subtraction",
		},
		{
			name:           "zero balance change",
			initialBalance: uint256.NewInt(100),
			balanceChange:  hexOrDecimal256FromInt64(0),
			accountExists:  true,
			wantBalance:    uint256.NewInt(100),
		},
		{
			name:           "nil balance change",
			initialBalance: uint256.NewInt(100),
			balanceChange:  nil,
			accountExists:  true,
			wantBalance:    uint256.NewInt(100),
		},
		{
			name:           "new account with positive balance",
			initialBalance: uint256.NewInt(0),
			balanceChange:  hexOrDecimal256FromInt64(1000),
			accountExists:  false,
			wantBalance:    uint256.NewInt(1000),
		},
		{
			name:           "new account with negative balance",
			initialBalance: uint256.NewInt(0),
			balanceChange:  hexOrDecimal256FromInt64(-100),
			accountExists:  false,
			wantBalance:    uint256.NewInt(0), // unchanged
			wantError:      "insufficient balance for subtraction",
		},
		{
			name:           "exact balance subtraction",
			initialBalance: uint256.NewInt(100),
			balanceChange:  hexOrDecimal256FromInt64(-100),
			accountExists:  true,
			wantBalance:    uint256.NewInt(0),
		},
		{
			name:           "off-by-one underflow",
			initialBalance: uint256.NewInt(100),
			balanceChange:  hexOrDecimal256FromInt64(-101),
			accountExists:  true,
			wantBalance:    uint256.NewInt(100), // unchanged
			wantError:      "insufficient balance for subtraction",
		},
		{
			name:           "large positive balance change",
			initialBalance: uint256.MustFromBig(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1000))),
			balanceChange:  hexOrDecimal256FromInt64(500),
			accountExists:  true,
			wantBalance:    uint256.MustFromBig(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(500))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh state database for each test
			statedb := createTestStateDB(t)

			// Set up initial account state
			if tt.accountExists {
				setAccountBalance(t, statedb, testAddr, tt.initialBalance)
			}

			// Create the upgrade configuration
			upgrade := extras.StateUpgradeAccount{
				BalanceChange: tt.balanceChange,
			}

			// Execute the upgrade
			err := upgradeAccount(testAddr, upgrade, statedb, false)

			// Check error expectations
			if tt.wantError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantError)
				require.Contains(t, err.Error(), testAddr.Hex())
			} else {
				require.NoError(t, err)
			}

			// Check final balance
			actualBalance := statedb.GetBalance(testAddr)
			require.Equal(t, tt.wantBalance, actualBalance, "final balance mismatch")

			// Verify account was created if it didn't exist
			require.True(t, statedb.Exist(testAddr), "account should exist after upgrade")
		})
	}
}

func TestUpgradeAccount_ErrorMessageFormat(t *testing.T) {
	statedb := createTestStateDB(t)
	addr := common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12")

	// Set initial balance to 50
	setAccountBalance(t, statedb, addr, uint256.NewInt(50))

	upgrade := extras.StateUpgradeAccount{
		BalanceChange: hexOrDecimal256FromInt64(-75),
	}

	err := upgradeAccount(addr, upgrade, statedb, false)
	require.Error(t, err)

	// Check specific error message components
	errorMsg := err.Error()
	require.Contains(t, errorMsg, "insufficient balance for subtraction")
	require.Contains(t, errorMsg, addr.Hex()) // EIP-55 checksum format
	require.Contains(t, errorMsg, "has 50")
	require.Contains(t, errorMsg, "trying to subtract 75")
}

func TestUpgradeAccount_CompleteUpgrade(t *testing.T) {
	statedb := createTestStateDB(t)
	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Create a complete state upgrade with balance change, code, and storage
	code := []byte{0xde, 0xad, 0xbe, 0xef}
	storageKey := common.HexToHash("0x1234")
	storageValue := common.HexToHash("0x5678")

	upgrade := extras.StateUpgradeAccount{
		BalanceChange: hexOrDecimal256FromInt64(1000),
		Code:          code,
		Storage:       map[common.Hash]common.Hash{storageKey: storageValue},
	}

	// Execute the upgrade
	err := upgradeAccount(addr, upgrade, statedb, true) // Test with EIP158 = true
	require.NoError(t, err)

	// Verify all changes were applied
	require.Equal(t, uint256.NewInt(1000), statedb.GetBalance(addr))
	require.True(t, statedb.Exist(addr))
	require.Equal(t, uint64(1), statedb.GetNonce(addr)) // Should be set to 1 due to EIP158 and code deployment

	// Cast to access code and storage (these aren't in our StateDB interface)
	extStateDB := statedb.(*extstate.StateDB)
	require.Equal(t, code, extStateDB.GetCode(addr))
	require.Equal(t, storageValue, extStateDB.GetState(addr, storageKey))
}

// createTestStateDB creates a real StateDB instance for testing
func createTestStateDB(t *testing.T) StateDB {
	t.Helper()

	// Create an in-memory database
	db := rawdb.NewMemoryDatabase()

	// Create a new state database
	statedb, err := state.New(types.EmptyRootHash, state.NewDatabase(db), nil)
	require.NoError(t, err)

	// Wrap it with extstate for predicate support (matching the interface)
	return extstate.New(statedb)
}

// setAccountBalance is a helper to set an account's initial balance for testing
func setAccountBalance(t *testing.T, statedb StateDB, addr common.Address, balance *uint256.Int) {
	t.Helper()

	// Cast to extstate.StateDB to access the underlying state
	extStateDB := statedb.(*extstate.StateDB)
	extStateDB.CreateAccount(addr)
	extStateDB.SetBalance(addr, balance)
}

// Helper function to create HexOrDecimal256 from int64
func hexOrDecimal256FromInt64(i int64) *math.HexOrDecimal256 {
	return (*math.HexOrDecimal256)(big.NewInt(i))
}
