// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TODO: edit comments
var (
	_ StatefulPrecompileConfig = &FeeConfigManagerConfig{}

	// Singleton StatefulPrecompiledContract for minting native assets by permissioned callers.
	FeeConfigManagerPrecompile StatefulPrecompiledContract = createNativeMinterPrecompile(ContractNativeMinterAddress)

	setFeeConfigSignature = CalculateFunctionSelector("setFeeConfig(uint256, uint64, uint256, uint256, uint256, uint256, uint256, uint256)")
	getFeeConfigSignature = CalculateFunctionSelector("getFeeConfig()")

	storageKey = common.Hash{}
)

// TODO: find a common place with this and params.FeeConfig
type FeeConfig struct {
	GasLimit        *big.Int `json:"gasLimit,omitempty"`
	TargetBlockRate uint64   `json:"targetBlockRate,omitempty"`

	MinBaseFee               *big.Int `json:"minBaseFee,omitempty"`
	TargetGas                *big.Int `json:"targetGas,omitempty"`
	BaseFeeChangeDenominator *big.Int `json:"baseFeeChangeDenominator,omitempty"`

	MinBlockGasCost  *big.Int `json:"minBlockGasCost,omitempty"`
	MaxBlockGasCost  *big.Int `json:"maxBlockGasCost,omitempty"`
	BlockGasCostStep *big.Int `json:"blockGasCostStep,omitempty"`
}

// FeeConfigManagerConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the contract deployer specific precompile address.
type FeeConfigManagerConfig struct {
	AllowListConfig
}

// Address returns the address of the fee config manager contract.
func (c *FeeConfigManagerConfig) Address() common.Address {
	return FeeConfigManagerAddress
}

// Configure configures [state] with the desired admins based on [c].
func (c *FeeConfigManagerConfig) Configure(state StateDB) {
	c.AllowListConfig.Configure(state, FeeConfigManagerAddress)
}

// Contract returns the singleton stateful precompiled contract to be used for the native minter.
func (c *FeeConfigManagerConfig) Contract() StatefulPrecompiledContract {
	return FeeConfigManagerPrecompile
}

// GetContractNativeMinterStatus returns the role of [address] for the minter list.
func GetFeeConfigManagerStatus(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, FeeConfigManagerAddress, address)
}

// SetContractNativeMinterStatus sets the permissions of [address] to [role] for the
// minter list. assumes [role] has already been verified as valid.
func SetFeeConfigManagerStatus(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, FeeConfigManagerAddress, address, role)
}

func getFeeConfig(state StateDB) (FeeConfig, error) {
	// Generate the state key for [address]
	feeHash := state.GetState(FeeConfigManagerAddress, storageKey)
	v := FeeConfig{}
	err := json.Unmarshal(feeHash.Bytes(), &v)
	return v, err
}

func setFeeConfig(stateDB StateDB, precompileAddr, feeConfig FeeConfig) error {
	feeBytes, err := json.Marshal(feeConfig)
	if err != nil {
		return err
	}
	// Assign [role] to the address
	stateDB.SetState(FeeConfigManagerAddress, storageKey, common.BytesToHash(feeBytes))
	return nil
}
