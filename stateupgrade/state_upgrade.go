// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/common/hexutil"
	"github.com/ava-labs/libevm/common/math"
	"github.com/holiman/uint256"
)

// StateUpgrade describes the modifications to be made to the state during
// a state upgrade.
type StateUpgrade struct {
	BlockTimestamp *uint64 `json:"blockTimestamp,omitempty"`

	// map from account address to the modification to be made to the account.
	StateUpgradeAccounts map[common.Address]StateUpgradeAccount `json:"accounts"`
}

// StateUpgradeAccount describes the modifications to be made to an account during
// a state upgrade.
type StateUpgradeAccount struct {
	Code          hexutil.Bytes               `json:"code,omitempty"`
	Storage       map[common.Hash]common.Hash `json:"storage,omitempty"`
	BalanceChange *math.HexOrDecimal256       `json:"balanceChange,omitempty"`
}

func (s *StateUpgrade) Equal(other *StateUpgrade) bool {
	return reflect.DeepEqual(s, other)
}

// VerifyStateUpgrades checks [c.StateUpgrades] is well formed:
// - the specified blockTimestamps must monotonically increase
func VerifyStateUpgrades(upgrades []StateUpgrade) error {
	var previousUpgradeTimestamp *uint64
	for i, upgrade := range upgrades {
		upgradeTimestamp := upgrade.BlockTimestamp
		if upgradeTimestamp == nil {
			return fmt.Errorf("StateUpgrade[%d]: config block timestamp cannot be nil ", i)
		}
		// Verify the upgrade's timestamp is equal 0 (to avoid confusion with genesis).
		if *upgradeTimestamp == 0 {
			return fmt.Errorf("StateUpgrade[%d]: config block timestamp (%v) must be greater than 0", i, *upgradeTimestamp)
		}

		// Verify specified timestamps are strictly monotonically increasing.
		if previousUpgradeTimestamp != nil && *upgradeTimestamp <= *previousUpgradeTimestamp {
			return fmt.Errorf("StateUpgrade[%d]: config block timestamp (%v) <= previous timestamp (%v)", i, *upgradeTimestamp, *previousUpgradeTimestamp)
		}
		previousUpgradeTimestamp = upgradeTimestamp
	}
	return nil
}

// Configure applies the state upgrade to the state.
func Configure(stateUpgrade *StateUpgrade, chainConfig ChainContext, state StateDB, blockContext BlockContext) {
	isEIP158 := chainConfig.IsEIP158(blockContext.Number())
	for account, upgrade := range stateUpgrade.StateUpgradeAccounts {
		upgradeAccount(account, upgrade, state, isEIP158)
	}
}

// upgradeAccount applies the state upgrade to the given account.
func upgradeAccount(account common.Address, upgrade StateUpgradeAccount, state StateDB, isEIP158 bool) {
	// Create the account if it does not exist
	if !state.Exist(account) {
		state.CreateAccount(account)
	}

	if upgrade.BalanceChange != nil {
		balanceChange, _ := uint256.FromBig((*big.Int)(upgrade.BalanceChange))
		state.AddBalance(account, balanceChange)
	}
	if len(upgrade.Code) != 0 {
		// if the nonce is 0, set the nonce to 1 as we would when deploying a contract at
		// the address.
		if isEIP158 && state.GetNonce(account) == 0 {
			state.SetNonce(account, 1)
		}
		state.SetCode(account, upgrade.Code)
	}
	for key, value := range upgrade.Storage {
		state.SetState(account, key, value)
	}
}
