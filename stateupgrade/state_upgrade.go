// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// Configure applies the state upgrade to the state.
func Configure(stateUpgrade *params.StateUpgrade, blockContext BlockContext, state StateDB) error {
	log.Info("Configuring state upgrade", "block", blockContext.Number(), "timestamp", blockContext.Timestamp())
	for account, upgrade := range stateUpgrade.Accounts {
		if err := upgradeAccount(account, upgrade, state); err != nil {
			return err
		}
	}
	return nil
}

// upgradeAccount applies the state upgrade to the given account.
func upgradeAccount(account common.Address, upgrade params.StateUpgradeAccount, state StateDB) error {
	// TODO: is this necessary?
	if !state.Exist(account) {
		// Create the account if it does not exist.
		state.CreateAccount(account)
	}

	if upgrade.BalanceChange != nil {
		state.AddBalance(account, upgrade.BalanceChange)
	}
	if upgrade.Code != nil {
		if state.GetNonce(account) == 0 {
			// If the nonce is 0, we will set it to a non-zero value
			// so the account is not considered empty.
			state.SetNonce(account, 1)
		}
		state.SetCode(account, upgrade.Code)
	}
	for key, value := range upgrade.Storage {
		state.SetState(account, key, value)
	}
	return nil
}
