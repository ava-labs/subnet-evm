package bibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	MARGIN_ACCOUNT_GENESIS_ADDRESS       = "0x0300000000000000000000000000000000000001"
	ORACLE_SLOT                    int64 = 4
	SUPPORTED_COLLATERAL_SLOT      int64 = 8
	MARGIN_MAPPING_SLOT            int64 = 10
	RESERVED_MARGIN_SLOT           int64 = 11
)

func GetNormalizedMargin(stateDB contract.StateDB, trader common.Address) *big.Int {
	assets := GetCollaterals(stateDB)
	margins := getMargins(stateDB, trader)
	return hu.GetNormalizedMargin(assets, margins)
}

func getMargins(stateDB contract.StateDB, trader common.Address) []*big.Int {
	numAssets := getCollateralCount(stateDB)
	margins := make([]*big.Int, numAssets)
	for i := uint8(0); i < numAssets; i++ {
		margins[i] = getMargin(stateDB, big.NewInt(int64(i)), trader)
	}
	return margins
}

func getMargin(stateDB contract.StateDB, idx *big.Int, trader common.Address) *big.Int {
	marginStorageSlot := crypto.Keccak256(append(common.LeftPadBytes(idx.Bytes(), 32), common.LeftPadBytes(big.NewInt(MARGIN_MAPPING_SLOT).Bytes(), 32)...))
	marginStorageSlot = crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), marginStorageSlot...))
	return fromTwosComplement(stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BytesToHash(marginStorageSlot)).Bytes())
}

func getReservedMargin(stateDB contract.StateDB, trader common.Address) *big.Int {
	baseMappingHash := crypto.Keccak256(append(common.LeftPadBytes(trader.Bytes(), 32), common.LeftPadBytes(big.NewInt(RESERVED_MARGIN_SLOT).Bytes(), 32)...))
	return stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BytesToHash(baseMappingHash)).Big()
}

func GetAvailableMargin(stateDB contract.StateDB, trader common.Address) *big.Int {
	output := getNotionalPositionAndMargin(stateDB, &GetNotionalPositionAndMarginInput{Trader: trader, IncludeFundingPayments: true, Mode: uint8(1)}) // Min_Allowable_Margin
	return hu.GetAvailableMargin_(output.NotionalPosition, output.Margin, getReservedMargin(stateDB, trader), GetMinAllowableMargin(stateDB))
}

func getOracleAddress(stateDB contract.StateDB) common.Address {
	return common.BytesToAddress(stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BigToHash(big.NewInt(ORACLE_SLOT))).Bytes())
}

func GetCollaterals(stateDB contract.StateDB) []hu.Collateral {
	numAssets := getCollateralCount(stateDB)
	assets := make([]hu.Collateral, numAssets)
	for i := uint8(0); i < numAssets; i++ {
		assets[i] = getCollateralAt(stateDB, i)
	}
	return assets
}

func getCollateralCount(stateDB contract.StateDB) uint8 {
	rawVal := stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BigToHash(big.NewInt(SUPPORTED_COLLATERAL_SLOT)))
	return uint8(new(big.Int).SetBytes(rawVal.Bytes()).Uint64())
}

func getCollateralAt(stateDB contract.StateDB, idx uint8) hu.Collateral {
	// struct Collateral { IERC20 token; uint weight; uint8 decimals; }
	baseSlot := hu.Add(collateralStorageSlot(), big.NewInt(int64(idx)*3)) // collateral struct size = 3 * 32 bytes
	tokenAddress := common.BytesToAddress(stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BigToHash(baseSlot)).Bytes())
	return hu.Collateral{
		Weight:   stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BigToHash(hu.Add(baseSlot, big.NewInt(1)))).Big(),
		Decimals: uint8(stateDB.GetState(common.HexToAddress(MARGIN_ACCOUNT_GENESIS_ADDRESS), common.BigToHash(hu.Add(baseSlot, big.NewInt(2)))).Big().Uint64()),
		Price:    getUnderlyingPrice_(stateDB, tokenAddress),
	}
}

func collateralStorageSlot() *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(common.BigToHash(big.NewInt(SUPPORTED_COLLATERAL_SLOT)).Bytes()))
}
