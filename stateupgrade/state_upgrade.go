// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
)

// Configure applies the state upgrade to the state.
func Configure(stateUpgrade *params.StateUpgrade, state StateDB) error {
	for account, upgrade := range stateUpgrade.StateUpgradeAccounts {
		if err := upgradeAccount(account, upgrade, state); err != nil {
			return err
		}
	}
	return nil
}

// upgradeAccount applies the state upgrade to the given account.
func upgradeAccount(account common.Address, upgrade params.StateUpgradeAccount, state StateDB) error {
	// Create the account if it does not exist
	if !state.Exist(account) {
		state.CreateAccount(account)
	}

	if upgrade.BalanceChange != nil {
		state.AddBalance(account, (*big.Int)(upgrade.BalanceChange))
	}
	if len(upgrade.Code) != 0 {
		// if the nonce is 0, set the nonce to 1 as we would when deploying a contract at
		// the address.
		if state.GetNonce(account) == 0 {
			state.SetNonce(account, 1)
		}
		state.SetCode(account, upgrade.Code)
	}
	for key, value := range upgrade.Storage {
		state.SetState(account, key, value)
	}
	return nil
}
