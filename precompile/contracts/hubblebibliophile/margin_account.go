package hubblebibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	MARGIN_ACCOUNT_GENESIS_ADDRESS        = "0x0300000000000000000000000000000000000070"
	VAR_MARGIN_MAPPING_STORAGE_SLOT int64 = 10
)

func GetNormalizedMargin(stateDB contract.StateDB, trader common.Address) *big.Int {
	// this is only written for single hUSD collateral
	// TODO: generalize for multiple collaterals
	return getMargin(stateDB, big.NewInt(0), trader)
}

func getMargin(stateDB contract.StateDB, collateralIdx *big.Int, trader common.Address) *big.Int {
	marginStorageSlot := crypto.Keccak256(append(common.LeftPadBytes(collateralIdx.Bytes(), 32), common.LeftPadBytes(big.NewInt(VAR_MARGIN_MAPPING_STORAGE_SLOT).Bytes(), 32)...))
	marginStorageSlot = crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), marginStorageSlot...))
	return fromTwosComplement(stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BytesToHash(marginStorageSlot)).Bytes())
}
