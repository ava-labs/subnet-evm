// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/libevm"
)

func (r RulesExtra) CanCreateContract(ac *libevm.AddressContext, state libevm.StateReader) error {
	// IsProhibited
	if ac.Self == constants.BlackholeAddr || modules.ReservedAddress(ac.Self) {
		return vmerrs.ErrAddrProhibited
	}

	// If the allow list is enabled, check that [ac.Origin] has permission to deploy a contract.
	if r.IsPrecompileEnabled(deployerallowlist.ContractAddress) {
		allowListRole := deployerallowlist.GetContractDeployerAllowListStatus(state, ac.Origin)
		if !allowListRole.IsEnabled() {
			ac.Gas = 0
			return fmt.Errorf("tx.origin %s is not authorized to deploy a contract", ac.Origin)
		}
	}

	return nil
}
