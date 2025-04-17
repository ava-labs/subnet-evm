// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"math/big"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/holiman/uint256"
)

// Configure applies the state upgrade to the state.
func Configure(stateUpgrade *extras.StateUpgrade, chainConfig ChainContext, state StateDB, blockContext BlockContext) error {
	isEIP158 := chainConfig.IsEIP158(blockContext.Number())
	for account, upgrade := range stateUpgrade.StateUpgradeAccounts {
		if err := upgradeAccount(account, upgrade, state, isEIP158); err != nil {
			return err
		}
	}
	return nil
}

// upgradeAccount applies the state upgrade to the given account.
func upgradeAccount(account common.Address, upgrade extras.StateUpgradeAccount, state StateDB, isEIP158 bool) error {
	// Create the account if it does not exist
	if !state.Exist(account) {
		state.CreateAccount(account)
	}

	if upgrade.BalanceChange != nil {
		balanceChange, _ := uint256.FromBig((*big.Int)(upgrade.BalanceChange))
		state.AddBalance(account, balanceChange)
	}

	// Set nonce if explicitly provided
	if upgrade.Nonce != nil {
		state.SetNonce(account, *upgrade.Nonce)
	} else if len(upgrade.Code) != 0 {
		// If no explicit nonce is provided but code is being set, set the nonce to
		// 1 as we would when deploying a contract at the address.
		if isEIP158 && state.GetNonce(account) == 0 {
			state.SetNonce(account, 1)
		}
	}

	if len(upgrade.Code) != 0 {
		state.SetCode(account, upgrade.Code)
	}
	for key, value := range upgrade.Storage {
		state.SetState(account, key, value)
	}
	return nil
}
