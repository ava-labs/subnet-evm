package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	TRADER_TO_REFERRER_SLOT int64 = 3
	RESTRICTED_INVITES_SLOT int64 = 6
)

func restrictedInvites(stateDB contract.StateDB, referralContract common.Address) bool {
	return stateDB.GetState(referralContract, common.BigToHash(big.NewInt(RESTRICTED_INVITES_SLOT))).Big().Uint64() == 1
}

func traderToReferrer(stateDB contract.StateDB, referralContract, trader common.Address) common.Address {
	pos := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(TRADER_TO_REFERRER_SLOT).Bytes(), 32)...))
	return common.BytesToAddress(stateDB.GetState(referralContract, common.BytesToHash(pos)).Bytes())
}

func HasReferrer(stateDB contract.StateDB, trader common.Address) bool {
	referralContract := getReferralAddress(stateDB)
	return !restrictedInvites(stateDB, referralContract) || traderToReferrer(stateDB, referralContract, trader) != common.Address{}
}
