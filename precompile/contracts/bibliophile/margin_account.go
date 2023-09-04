package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/log"
)

const (
	MARGIN_ACCOUNT_GENESIS_ADDRESS        = "0x0300000000000000000000000000000000000001"
	VAR_MARGIN_MAPPING_STORAGE_SLOT int64 = 10
	VAR_RESERVED_MARGIN_SLOT        int64 = 11
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

func getReservedMargin(stateDB contract.StateDB, trader common.Address) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(VAR_RESERVED_MARGIN_SLOT).Bytes(), 32)...))
	return stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BytesToHash(baseMappingHash)).Big()
}

// Monday, 4 September 2023 10:05:00
var V5ActivationDate *big.Int = new(big.Int).SetInt64(1693821900)

func GetAvailableMargin(stateDB contract.StateDB, trader common.Address, blockTimestamp *big.Int) *big.Int {
	includeFundingPayment := true
	mode := uint8(1) // Min_Allowable_Margin
	var output GetNotionalPositionAndMarginOutput
	if blockTimestamp != nil && blockTimestamp.Cmp(V5ActivationDate) == 1 {
		output = GetNotionalPositionAndMargin(stateDB, &GetNotionalPositionAndMarginInput{Trader: trader, IncludeFundingPayments: includeFundingPayment, Mode: mode}, blockTimestamp)
	} else {
		output = GetNotionalPositionAndMargin(stateDB, &GetNotionalPositionAndMarginInput{Trader: trader, IncludeFundingPayments: includeFundingPayment, Mode: mode}, nil)
	}
	notionalPostion := output.NotionalPosition
	margin := output.Margin
	utitlizedMargin := divide1e6(big.NewInt(0).Mul(notionalPostion, GetMinAllowableMargin(stateDB)))
	reservedMargin := getReservedMargin(stateDB, trader)
	// log.Info("GetAvailableMargin", "trader", trader, "notionalPostion", notionalPostion, "margin", margin, "utitlizedMargin", utitlizedMargin, "reservedMargin", reservedMargin)
	return big.NewInt(0).Sub(big.NewInt(0).Sub(margin, utitlizedMargin), reservedMargin)
}
