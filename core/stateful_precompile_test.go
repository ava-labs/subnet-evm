// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/txallowlist"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

// TODO: move this to precompile package once cross-import is resolved

var (
	_ precompile.BlockContext              = &mockBlockContext{}
	_ precompile.PrecompileAccessibleState = &mockAccessibleState{}

	testFeeConfig = commontype.FeeConfig{
		GasLimit:        big.NewInt(8_000_000),
		TargetBlockRate: 2, // in seconds

		MinBaseFee:               big.NewInt(25_000_000_000),
		TargetGas:                big.NewInt(15_000_000),
		BaseFeeChangeDenominator: big.NewInt(36),

		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		BlockGasCostStep: big.NewInt(200_000),
	}

	testBlockNumber = big.NewInt(7)
)

type mockBlockContext struct {
	blockNumber *big.Int
	timestamp   uint64
}

func (mb *mockBlockContext) Number() *big.Int    { return mb.blockNumber }
func (mb *mockBlockContext) Timestamp() *big.Int { return new(big.Int).SetUint64(mb.timestamp) }

type mockAccessibleState struct {
	state        *state.StateDB
	blockContext *mockBlockContext
	snowContext  *snow.Context
}

func (m *mockAccessibleState) GetStateDB() precompile.StateDB { return m.state }

func (m *mockAccessibleState) GetBlockContext() precompile.BlockContext { return m.blockContext }

func (m *mockAccessibleState) GetSnowContext() *snow.Context { return m.snowContext }

func (m *mockAccessibleState) CallFromPrecompile(caller common.Address, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	return nil, 0, nil
}

// This test is added within the core package so that it can import all of the required code
// without creating any import cycles
func TestContractDeployerAllowListRun(t *testing.T) {
	type test struct {
		caller      common.Address
		input       func() []byte
		suppliedGas uint64
		readOnly    bool

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]test{
		"set admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := deployerallowlist.GetContractDeployerAllowListStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListAdmin, res)
			},
		},
		"set deployer": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := deployerallowlist.GetContractDeployerAllowListStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListEnabled, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := deployerallowlist.GetContractDeployerAllowListStatus(state, adminAddr)
				require.Equal(t, precompile.AllowListNoRole, res)
			},
		},
		"set no role from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set deployer from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set admin from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list admin role": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list out of gas": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost - 1,
			readOnly:    true,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			deployerallowlist.SetContractDeployerAllowListStatus(state, adminAddr, precompile.AllowListAdmin)
			deployerallowlist.SetContractDeployerAllowListStatus(state, noRoleAddr, precompile.AllowListNoRole)
			require.Equal(t, precompile.AllowListAdmin, deployerallowlist.GetContractDeployerAllowListStatus(state, adminAddr))
			require.Equal(t, precompile.AllowListNoRole, deployerallowlist.GetContractDeployerAllowListStatus(state, noRoleAddr))

			blockContext := &mockBlockContext{blockNumber: common.Big0}
			ret, remainingGas, err := deployerallowlist.ContractDeployerAllowListPrecompile.Run(&mockAccessibleState{state: state, blockContext: blockContext, snowContext: snow.DefaultContextTest()}, test.caller, deployerallowlist.Address, test.input(), test.suppliedGas, test.readOnly)
			if len(test.expectedErr) != 0 {
				require.ErrorContains(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.expectedRes, ret)

			if test.assertState != nil {
				test.assertState(t, state)
			}
		})
	}
}

func TestTxAllowListRun(t *testing.T) {
	type test struct {
		caller         common.Address
		precompileAddr common.Address
		input          func() []byte
		suppliedGas    uint64
		readOnly       bool

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]test{
		"set admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := txallowlist.GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListAdmin, res)
			},
		},
		"set allowed": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := txallowlist.GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListEnabled, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := txallowlist.GetTxAllowListStatus(state, adminAddr)
				require.Equal(t, precompile.AllowListNoRole, res)
			},
		},
		"set no role from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set allowed from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set admin from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			caller:         adminAddr,
			precompileAddr: txallowlist.Address,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list admin role": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list out of gas": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost - 1,
			readOnly:    true,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			txallowlist.SetTxAllowListStatus(state, adminAddr, precompile.AllowListAdmin)
			require.Equal(t, precompile.AllowListAdmin, txallowlist.GetTxAllowListStatus(state, adminAddr))

			blockContext := &mockBlockContext{blockNumber: common.Big0}
			ret, remainingGas, err := txallowlist.TxAllowListPrecompile.Run(&mockAccessibleState{state: state, blockContext: blockContext, snowContext: snow.DefaultContextTest()}, test.caller, txallowlist.Address, test.input(), test.suppliedGas, test.readOnly)
			if len(test.expectedErr) != 0 {
				require.ErrorContains(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.expectedRes, ret)

			if test.assertState != nil {
				test.assertState(t, state)
			}
		})
	}
}

func TestContractNativeMinterRun(t *testing.T) {
	type test struct {
		caller      common.Address
		input       func() []byte
		suppliedGas uint64
		readOnly    bool
		config      *nativeminter.ContractNativeMinterConfig

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")
	testAddr := common.HexToAddress("0x123456789")

	for name, test := range map[string]test{
		"mint funds from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(noRoleAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedErr: nativeminter.ErrCannotMint.Error(),
		},
		"mint funds from enabled address": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"enabled role by config": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(testAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListEnabled).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, precompile.AllowListEnabled, nativeminter.GetContractNativeMinterStatus(state, testAddr))
			},
			config: &nativeminter.ContractNativeMinterConfig{
				AllowListConfig: precompile.AllowListConfig{EnabledAddresses: []common.Address{testAddr}},
			},
		},
		"initial mint funds": {
			caller: enabledAddr,
			config: &nativeminter.ContractNativeMinterConfig{
				InitialMint: map[common.Address]*math.HexOrDecimal256{
					enabledAddr: math.NewHexOrDecimal256(2),
				},
			},
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big2, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"mint funds from admin address": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"mint max big funds": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, math.MaxBig256)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, math.MaxBig256, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"readOnly mint with noRole fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with allow role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with admin role fails": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas mint from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read from noRole address": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {},
		},
		"read from noRole address readOnly enabled": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {},
		},
		"read from noRole address with insufficient gas": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"set allow role from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := nativeminter.GetContractNativeMinterStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListEnabled, res)
			},
		},
		"set allow role from non-admin fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			nativeminter.SetContractNativeMinterStatus(state, adminAddr, precompile.AllowListAdmin)
			nativeminter.SetContractNativeMinterStatus(state, enabledAddr, precompile.AllowListEnabled)
			nativeminter.SetContractNativeMinterStatus(state, noRoleAddr, precompile.AllowListNoRole)
			require.Equal(t, precompile.AllowListAdmin, nativeminter.GetContractNativeMinterStatus(state, adminAddr))
			require.Equal(t, precompile.AllowListEnabled, nativeminter.GetContractNativeMinterStatus(state, enabledAddr))
			require.Equal(t, precompile.AllowListNoRole, nativeminter.GetContractNativeMinterStatus(state, noRoleAddr))

			blockContext := &mockBlockContext{blockNumber: common.Big0}
			if test.config != nil {
				test.config.Configure(params.TestChainConfig, state, blockContext)
			}
			ret, remainingGas, err := nativeminter.ContractNativeMinterPrecompile.Run(&mockAccessibleState{state: state, blockContext: blockContext, snowContext: snow.DefaultContextTest()}, test.caller, nativeminter.Address, test.input(), test.suppliedGas, test.readOnly)
			if len(test.expectedErr) != 0 {
				require.ErrorContains(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.expectedRes, ret)

			if test.assertState != nil {
				test.assertState(t, state)
			}
		})
	}
}

func TestFeeConfigManagerRun(t *testing.T) {
	type test struct {
		caller       common.Address
		preCondition func(t *testing.T, state *state.StateDB)
		input        func() []byte
		suppliedGas  uint64
		readOnly     bool
		config       *feemanager.FeeConfigManagerConfig

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]test{
		"set config from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := feemanager.PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost,
			readOnly:    false,
			expectedErr: feemanager.ErrCannotChangeFee.Error(),
		},
		"set config from enabled address": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := feemanager.PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := feemanager.GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set invalid config from enabled address": {
			caller: enabledAddr,
			input: func() []byte {
				feeConfig := testFeeConfig
				feeConfig.MinBlockGasCost = new(big.Int).Mul(feeConfig.MaxBlockGasCost, common.Big2)
				input, err := feemanager.PackSetFeeConfig(feeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost,
			readOnly:    false,
			expectedRes: nil,
			config: &feemanager.FeeConfigManagerConfig{
				InitialFeeConfig: &testFeeConfig,
			},
			expectedErr: "cannot be greater than maxBlockGasCost",
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := feemanager.GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set config from admin address": {
			caller: adminAddr,
			input: func() []byte {
				input, err := feemanager.PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := feemanager.GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
				lastChangedAt := feemanager.GetFeeConfigLastChangedAt(state)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get fee config from non-enabled address": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				err := feemanager.StoreFeeConfig(state, testFeeConfig, &mockBlockContext{blockNumber: big.NewInt(6)})
				require.NoError(t, err)
			},
			input: func() []byte {
				return feemanager.PackGetFeeConfigInput()
			},
			suppliedGas: feemanager.GetFeeConfigGasCost,
			readOnly:    true,
			expectedRes: func() []byte {
				res, err := feemanager.PackFeeConfig(testFeeConfig)
				require.NoError(t, err)
				return res
			}(),
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := feemanager.GetStoredFeeConfig(state)
				lastChangedAt := feemanager.GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, big.NewInt(6), lastChangedAt)
			},
		},
		"get initial fee config": {
			caller: noRoleAddr,
			input: func() []byte {
				return feemanager.PackGetFeeConfigInput()
			},
			suppliedGas: feemanager.GetFeeConfigGasCost,
			config: &feemanager.FeeConfigManagerConfig{
				InitialFeeConfig: &testFeeConfig,
			},
			readOnly: true,
			expectedRes: func() []byte {
				res, err := feemanager.PackFeeConfig(testFeeConfig)
				require.NoError(t, err)
				return res
			}(),
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := feemanager.GetStoredFeeConfig(state)
				lastChangedAt := feemanager.GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get last changed at from non-enabled address": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				err := feemanager.StoreFeeConfig(state, testFeeConfig, &mockBlockContext{blockNumber: testBlockNumber})
				require.NoError(t, err)
			},
			input: func() []byte {
				return feemanager.PackGetLastChangedAtInput()
			},
			suppliedGas: feemanager.GetLastChangedAtGasCost,
			readOnly:    true,
			expectedRes: common.BigToHash(testBlockNumber).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := feemanager.GetStoredFeeConfig(state)
				lastChangedAt := feemanager.GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
		},
		"readOnly setFeeConfig with noRole fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := feemanager.PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with allow role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := feemanager.PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with admin role fails": {
			caller: adminAddr,
			input: func() []byte {
				input, err := feemanager.PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas setFeeConfig from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := feemanager.PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: feemanager.SetFeeConfigGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"set allow role from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := feemanager.GetFeeConfigManagerStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListEnabled, res)
			},
		},
		"set allow role from non-admin fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			feemanager.SetFeeConfigManagerStatus(state, adminAddr, precompile.AllowListAdmin)
			feemanager.SetFeeConfigManagerStatus(state, enabledAddr, precompile.AllowListEnabled)
			feemanager.SetFeeConfigManagerStatus(state, noRoleAddr, precompile.AllowListNoRole)

			if test.preCondition != nil {
				test.preCondition(t, state)
			}

			blockContext := &mockBlockContext{blockNumber: testBlockNumber}
			if test.config != nil {
				test.config.Configure(params.TestChainConfig, state, blockContext)
			}
			ret, remainingGas, err := feemanager.FeeConfigManagerPrecompile.Run(&mockAccessibleState{state: state, blockContext: blockContext, snowContext: snow.DefaultContextTest()}, test.caller, feemanager.Address, test.input(), test.suppliedGas, test.readOnly)
			if len(test.expectedErr) != 0 {
				require.ErrorContains(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.expectedRes, ret)

			if test.assertState != nil {
				test.assertState(t, state)
			}
		})
	}
}

func TestRewardManagerRun(t *testing.T) {
	type test struct {
		caller       common.Address
		preCondition func(t *testing.T, state *state.StateDB)
		input        func() []byte
		suppliedGas  uint64
		readOnly     bool
		config       *rewardmanager.RewardManagerConfig

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")
	testAddr := common.HexToAddress("0x0123")

	for name, test := range map[string]test{
		"set allow fee recipients from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := rewardmanager.PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.AllowFeeRecipientsGasCost,
			readOnly:    false,
			expectedErr: rewardmanager.ErrCannotAllowFeeRecipients.Error(),
		},
		"set reward address from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := rewardmanager.PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.SetRewardAddressGasCost,
			readOnly:    false,
			expectedErr: rewardmanager.ErrCannotSetRewardAddress.Error(),
		},
		"disable rewards from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := rewardmanager.PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.DisableRewardsGasCost,
			readOnly:    false,
			expectedErr: rewardmanager.ErrCannotDisableRewards.Error(),
		},
		"set allow fee recipients from enabled succeeds": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.AllowFeeRecipientsGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				_, isFeeRecipients := rewardmanager.GetStoredRewardAddress(state)
				require.True(t, isFeeRecipients)
			},
		},
		"set reward address from enabled succeeds": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.SetRewardAddressGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				address, isFeeRecipients := rewardmanager.GetStoredRewardAddress(state)
				require.Equal(t, testAddr, address)
				require.False(t, isFeeRecipients)
			},
		},
		"disable rewards from enabled succeeds": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.DisableRewardsGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				address, isFeeRecipients := rewardmanager.GetStoredRewardAddress(state)
				require.False(t, isFeeRecipients)
				require.Equal(t, constants.BlackholeAddr, address)
			},
		},
		"get current reward address from no role succeeds": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				rewardmanager.StoreRewardAddress(state, testAddr)
			},
			input: func() []byte {
				input, err := rewardmanager.PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.CurrentRewardAddressGasCost,
			readOnly:    false,
			expectedRes: func() []byte {
				res, err := rewardmanager.PackCurrentRewardAddressOutput(testAddr)
				require.NoError(t, err)
				return res
			}(),
		},
		"get are fee recipients allowed from no role succeeds": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				rewardmanager.EnableAllowFeeRecipients(state)
			},
			input: func() []byte {
				input, err := rewardmanager.PackAreFeeRecipientsAllowed()
				require.NoError(t, err)
				return input
			},
			suppliedGas: rewardmanager.AreFeeRecipientsAllowedGasCost,
			readOnly:    false,
			expectedRes: func() []byte {
				res, err := rewardmanager.PackAreFeeRecipientsAllowedOutput(true)
				require.NoError(t, err)
				return res
			}(),
		},
		"get initial config with address": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := rewardmanager.PackCurrentRewardAddress()
				require.NoError(t, err)
				return input
			},
			suppliedGas: rewardmanager.CurrentRewardAddressGasCost,
			config: &rewardmanager.RewardManagerConfig{
				InitialRewardConfig: &rewardmanager.InitialRewardConfig{
					RewardAddress: testAddr,
				},
			},
			readOnly: false,
			expectedRes: func() []byte {
				res, err := rewardmanager.PackCurrentRewardAddressOutput(testAddr)
				require.NoError(t, err)
				return res
			}(),
		},
		"get initial config with allow fee recipients enabled": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := rewardmanager.PackAreFeeRecipientsAllowed()
				require.NoError(t, err)
				return input
			},
			suppliedGas: rewardmanager.AreFeeRecipientsAllowedGasCost,
			config: &rewardmanager.RewardManagerConfig{
				InitialRewardConfig: &rewardmanager.InitialRewardConfig{
					AllowFeeRecipients: true,
				},
			},
			readOnly: false,
			expectedRes: func() []byte {
				res, err := rewardmanager.PackAreFeeRecipientsAllowedOutput(true)
				require.NoError(t, err)
				return res
			}(),
		},
		"readOnly allow fee recipients with allowed role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.AllowFeeRecipientsGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly set reward addresss with allowed role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.SetRewardAddressGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas set reward address from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.SetRewardAddressGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas allow fee recipients from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.AllowFeeRecipientsGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas read current reward address from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.CurrentRewardAddressGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas are fee recipients allowed from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := rewardmanager.PackAreFeeRecipientsAllowed()
				require.NoError(t, err)

				return input
			},
			suppliedGas: rewardmanager.AreFeeRecipientsAllowedGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"set allow role from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := rewardmanager.GetRewardManagerAllowListStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListEnabled, res)
			},
		},
		"set allow role from non-admin fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			rewardmanager.SetRewardManagerAllowListStatus(state, adminAddr, precompile.AllowListAdmin)
			rewardmanager.SetRewardManagerAllowListStatus(state, enabledAddr, precompile.AllowListEnabled)
			rewardmanager.SetRewardManagerAllowListStatus(state, noRoleAddr, precompile.AllowListNoRole)

			if test.preCondition != nil {
				test.preCondition(t, state)
			}

			blockContext := &mockBlockContext{blockNumber: testBlockNumber}
			if test.config != nil {
				test.config.Configure(params.TestChainConfig, state, blockContext)
			}
			ret, remainingGas, err := rewardmanager.RewardManagerPrecompile.Run(&mockAccessibleState{state: state, blockContext: blockContext, snowContext: snow.DefaultContextTest()}, test.caller, rewardmanager.Address, test.input(), test.suppliedGas, test.readOnly)
			if len(test.expectedErr) != 0 {
				require.ErrorContains(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.expectedRes, ret)

			if test.assertState != nil {
				test.assertState(t, state)
			}
		})
	}
}
